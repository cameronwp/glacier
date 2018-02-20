package cmd

import (
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
	totalJobs, err := addJobs(files, jq, chunker)
	if err != nil {
		return err
	}

	// where uploading individual parts happens
	upload := func(c *jobqueue.Chunk) error {
		// note that this does not necessarily initiate a new upload
		_, err := InitiateMultiPartUpload(svc, c.Path, fs.DefaultPartSize)
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

		_, err = fb.FetchAndHash(f, c.Path, c.StartB, c.EndB)
		if err != nil {
			return err
		}
		// fmt.Printf("uploading %s, %d-%d\n", c.Path, c.StartB, c.EndB)
		// fmt.Printf("buf len %d | hash %x\n", len(fileChunk.Buf), fileChunk.SHA256)
		return nil
	}

	drain := drain.NewDrain(upload)
	go drain.Drain(jq)

	// listen for upload status updates
	for i := 0; i < totalJobs*2; i++ {
		status := <-drain.Schan
		if status.State == jobqueue.Completed {
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

			// TODO: actually complete the upload
			fmt.Printf("%s done | treehash %x\n", status.Chunk.Path, treehash)
		}
	}

	// TODO: setup https://github.com/gizak/termui

	// TODO:
	// * upload files
	// * calculate rate of upload
	// * display file status
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

// InitiateMultiPartUpload will either return the UploadID because it was
// already initiated, or actually initiate it.
func InitiateMultiPartUpload(svc glacieriface.GlacierAPI, path string, partsize int64) (string, error) {
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
func GetDescription(target string, path string) (string, error) {
	baseName, err := filepath.Rel(target, path)
	if err != nil {
		return "", err
	}

	if baseName == "." {
		// the target is a single file
		baseName = filepath.Base(path)
	}

	return baseName, nil
}

// func uploadFileMultipart(svc *glacier.Glacier, fp string) error {
// 	var partSize = int64(1 << 20) // 1MB
// 	baseName := filepath.Base(fp)

// 	// TODO: move this to awsiface
// 	initResult, err := svc.InitiateMultiPartUpload(&glacier.InitiateMultiPartUploadInput{
// 		AccountId:          aws.String("-"),
// 		ArchiveDescription: aws.String(baseName),
// 		PartSize:           aws.String(fmt.Sprintf("%d", partSize)),
// 		VaultName:          aws.String(vault),
// 	})
// 	if err != nil {
// 		return formatAWSError(err)
// 	}

// 	f, err := os.Open(fp)
// 	if err != nil {
// 		return err
// 	}

// 	// TODO: zip the file first

// 	hashes := [][]byte{}

// 	var totalSize int64
// 	if stats, err := f.Stat(); err == nil {
// 		totalSize = stats.Size()
// 	} else {
// 		return err
// 	}

// 	bar := pb.ProgressBarTemplate(fmt.Sprintf(`%s: {{bar . | green}} {{counters . | blue }}`, baseName)).Start64(totalSize)

// 	startB := int64(0)
// 	var wg sync.WaitGroup
// 	for {
// 		wg.Add(1)

// 		// either the part size, or the amount of file remaining, whichever is smaller
// 		contentLength := int(math.Min(float64(partSize), float64(totalSize-startB)))
// 		buf := make([]byte, contentLength)
// 		n, _ := io.ReadFull(f, buf)
// 		if n == 0 {
// 			wg.Done()
// 			break
// 		}

// 		endB := startB + int64(n)

// 		hash := sha256.Sum256(buf[:n])
// 		hashes = append(hashes, hash[:])

// 		go func(b []byte, s int64, e int64, h string) {
// 			_, err := svc.UploadMultipartPart(&glacier.UploadMultipartPartInput{
// 				AccountId: aws.String("-"),
// 				Body:      bytes.NewReader(buf),
// 				Checksum:  aws.String(h),
// 				Range:     aws.String(fmt.Sprintf("bytes %d-%d/*", s, e-1)),
// 				UploadId:  aws.String(*initResult.UploadId),
// 				VaultName: aws.String(vault),
// 			})
// 			if err != nil {
// 				// TODO: queue the part to be reuploaded
// 				panic(formatAWSError(err))
// 			}
// 			bar.Add(contentLength)
// 			wg.Done()
// 		}(buf, startB, endB, fmt.Sprintf("%x", hash))

// 		startB = endB
// 	}

// 	wg.Wait()

// 	input := &glacier.CompleteMultipartUploadInput{
// 		AccountId:   aws.String("-"),
// 		ArchiveSize: aws.String(fmt.Sprintf("%d", totalSize)),
// 		Checksum:    aws.String(fmt.Sprintf("%x", glacier.ComputeTreeHash(hashes))),
// 		UploadId:    aws.String(*initResult.UploadId),
// 		VaultName:   aws.String(vault),
// 	}
// 	result, err := svc.CompleteMultipartUpload(input)
// 	if err != nil {
// 		return formatAWSError(err)
// 	}

// 	bar.Finish()

// 	// TODO: sync the archive with an S3 bucket
// 	fmt.Println(result)
// 	fmt.Println(*initResult.UploadId)

// 	return nil
// }

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
