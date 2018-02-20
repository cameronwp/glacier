package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/cameronwp/glacier/drain"
	"github.com/cameronwp/glacier/filebuffer"
	"github.com/cameronwp/glacier/fs"
	"github.com/cameronwp/glacier/jobqueue"
	"github.com/spf13/cobra"
)

var target string

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file or directory to Glacier",
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := fs.GetFilepaths(target)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return fmt.Errorf("invalid target: no files found")
		}

		jq := jobqueue.NewJobQueue(0, true)
		osChunker := &fs.OSChunker{}
		osBuffer := &filebuffer.OSBuffer{}

		return execute(svc, jq, osChunker, osBuffer, files)
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&target, "target", "t", "", "Path to file or directory to upload")
	err := uploadCmd.MarkFlagRequired("target")
	if err != nil {
		log.Fatal(err)
	}

	uploadCmd.Flags().StringVarP(&vault, "vault", "v", "", "Vault name")
	err = uploadCmd.MarkFlagRequired("vault")
	if err != nil {
		log.Fatal(err)
	}
}

func execute(svc glacieriface.GlacierAPI, jq jobqueue.FIFOQueuer, chunker fs.Chunker, fb filebuffer.BufferFetcherHasher, files []string) error {
	_, err := addJobs(files, jq, chunker)
	if err != nil {
		return err
	}

	// where uploading individual parts happens
	upload := func(c *jobqueue.Chunk) error {
		uploadID, err := GetUploadID(svc, c.Path, fs.DefaultPartSize)
		if err != nil {
			return err
		}

		f, err := os.Open(c.Path)
		if err != nil {
			return err
		}
		defer func() {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}()

		fileChunk, err := fb.FetchAndHash(f, c.Path, c.StartB, c.EndB)
		if err != nil {
			return err
		}

		_, err = svc.UploadMultipartPart(&glacier.UploadMultipartPartInput{
			AccountId: aws.String("-"),
			Body:      bytes.NewReader(fileChunk.Buf),
			Checksum:  aws.String(fmt.Sprintf("%x", fileChunk.SHA256)),
			Range:     aws.String(fmt.Sprintf("bytes %d-%d/*", c.StartB, c.EndB-1)),
			UploadId:  aws.String(uploadID),
			VaultName: aws.String(vault),
		})
		if err != nil {
			return formatAWSError(err)
		}
		return nil
	}

	drain := drain.NewDrain(upload)
	go drain.Drain(jq)

	collectResults(svc, chunker, fb, drain)

	return nil
}

func addJobs(files []string, jq jobqueue.FIFOQueuer, chunker fs.Chunker) (int, error) {
	numJobs := 0
	for _, path := range files {
		chunks, err := fs.ChunkFile(chunker, path, fs.DefaultPartSize)
		if err != nil {
			return numJobs, err
		}

		for _, c := range chunks {
			_, err := jq.Add(c)
			if err != nil {
				return numJobs, err
			}
			numJobs++
		}
	}
	return numJobs, nil
}

var uploadIDs = make(map[string]string)

// GetUploadID will either return the `UploadID` because it was already
// initiated, or actually initiate it.
func GetUploadID(svc glacieriface.GlacierAPI, path string, partsize int64) (string, error) {
	if uploadID, ok := uploadIDs[path]; ok {
		return uploadID, nil
	}

	desc, err := GetDescription(target, path)
	if err != nil {
		return "", err
	}

	initResult, err := svc.InitiateMultipartUpload(&glacier.InitiateMultipartUploadInput{
		AccountId:          aws.String("-"),
		ArchiveDescription: aws.String(desc),
		PartSize:           aws.String(fmt.Sprintf("%d", partsize)),
		VaultName:          aws.String(vault),
	})
	if err != nil {
		return "", formatAWSError(err)
	}

	uploadIDs[path] = *initResult.UploadId

	return *initResult.UploadId, nil
}

// GetDescription gets the relative path to the file from the target directory.
// `target` can be a relative path, while `path` must be an absolute path.
func GetDescription(target string, path string) (string, error) {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}

	baseName, err := filepath.Rel(absTarget, path)
	if err != nil {
		return "", err
	}

	if baseName == "." {
		// the target is a single file
		baseName = filepath.Base(path)
	}

	return baseName, nil
}

func collectResults(svc glacieriface.GlacierAPI, chunker fs.Chunker, fb filebuffer.BufferFetcherHasher, drain *drain.Drain) {
	// termui here
	go onComplete(svc, chunker, fb, drain)
	go onError(drain)
}

func onComplete(svc glacieriface.GlacierAPI, chunker fs.Chunker, fb filebuffer.BufferFetcherHasher, drain *drain.Drain) {
	for {
		status := <-drain.Schan

		if status.State != jobqueue.Completed {
			continue
		}

		// check if the whole file is done, and get the TreeHash if so
		treehash, err := fb.TreeHash(chunker, status.Chunk.Path)
		if err != nil {
			if err == filebuffer.ErrMissingFileChunks {
				// file isn't done, no biggie
				continue
			} else {
				fmt.Printf("unexpected treehash error | %s\n", err)
				continue
			}
		}

		// file is done, complete the upload
		totalSize, err := chunker.GetFilesize(status.Chunk.Path)
		if err != nil {
			fmt.Printf("unexpected error before completing upload | %s\n", err)
			continue
		}

		uploadID, err := GetUploadID(svc, status.Chunk.Path, fs.DefaultPartSize)
		if err != nil {
			fmt.Printf("unexpected error before completing upload | %s\n", err)
			continue
		}

		result, err := svc.CompleteMultipartUpload(&glacier.CompleteMultipartUploadInput{
			AccountId:   aws.String("-"),
			ArchiveSize: aws.String(fmt.Sprintf("%d", totalSize)),
			Checksum:    aws.String(fmt.Sprintf("%x", treehash)),
			UploadId:    aws.String(uploadID),
			VaultName:   aws.String(vault),
		})
		if err != nil {
			fmt.Println(formatAWSError(err))
		}

		// TODO: this needs to be logged somewhere
		fmt.Println(result)
	}
}

func onError(drain *drain.Drain) {
	for {
		status := <-drain.Echan
		fmt.Println(status)
	}
}

func formatAWSError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case glacier.ErrCodeResourceNotFoundException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeResourceNotFoundException, aerr.Error())
		case glacier.ErrCodeInvalidParameterValueException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeInvalidParameterValueException, aerr.Error())
		case glacier.ErrCodeMissingParameterValueException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeMissingParameterValueException, aerr.Error())
		case glacier.ErrCodeRequestTimeoutException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeRequestTimeoutException, aerr.Error())
		case glacier.ErrCodeServiceUnavailableException:
			return fmt.Errorf("%s | %s", glacier.ErrCodeServiceUnavailableException, aerr.Error())
		default:
			return fmt.Errorf("%s", aerr.Error())
		}
	}
	return fmt.Errorf(err.Error())
}
