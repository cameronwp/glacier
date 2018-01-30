package cmd

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/nickvanw/treehash"
	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v2"
)

var target string

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file or directory to Glacier",
	RunE: func(cmd *cobra.Command, args []string) error {
		files := make(map[string]struct{})
		getFiles(target, files)
		if len(files) == 0 {
			return fmt.Errorf("invalid target: no file(s) found")
		}

		for fp := range files {
			err := uploadFileMultipart(svc, fp)
			if err != nil {
				return err
			}
		}

		return nil
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

// Just the filepaths of the files to get uploaded
func getFiles(fp string, files map[string]struct{}) {
	maybeFileOrDirectory, err := os.Stat(fp)
	if err != nil {
		// no more files
		return
	}

	// recurse over the directory until files are found
	if maybeFileOrDirectory.IsDir() {
		err := filepath.Walk(fp, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// don't run against the root of the dir again
			if fp == path {
				return nil
			}

			getFiles(path, files)
			return nil
		})

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	files[fp] = struct{}{}
}

func uploadFileMultipart(svc *glacier.Glacier, fp string) error {
	partSize := int64(1 << 20) // 1MB
	baseName := filepath.Base(fp)

	initResult, err := svc.InitiateMultipartUpload(&glacier.InitiateMultipartUploadInput{
		AccountId:          aws.String("-"),
		ArchiveDescription: aws.String(baseName),
		PartSize:           aws.String(fmt.Sprintf("%d", partSize)),
		VaultName:          aws.String(vault),
	})
	if err != nil {
		return formatAWSError(err)
	}

	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	var th treehash.MultiTreeHash

	var totalSize int64
	if stats, err := f.Stat(); err == nil {
		totalSize = stats.Size()
	} else {
		return err
	}

	bar := pb.ProgressBarTemplate(fmt.Sprintf(`%s: {{bar . | green}} {{counters . | blue }}`, baseName)).Start64(totalSize)

	startB := int64(0)
	var wg sync.WaitGroup
	for {
		wg.Add(1)

		// either the part size, or the amount of file remaining, whichever is smaller
		contentLength := int(math.Min(float64(partSize), float64(totalSize-startB)))
		buf := make([]byte, contentLength)
		n, _ := io.ReadFull(f, buf)
		if n == 0 {
			wg.Done()
			break
		}

		endB := startB + int64(n)

		hash := fmt.Sprintf("%x", sha256.Sum256(buf[:n]))
		th.Add(hash)

		go func(b []byte, s int64, e int64, h string) {
			_, err := svc.UploadMultipartPart(&glacier.UploadMultipartPartInput{
				AccountId: aws.String("-"),
				Body:      bytes.NewReader(buf),
				Checksum:  aws.String(h),
				Range:     aws.String(fmt.Sprintf("bytes %d-%d/*", s, e-1)),
				UploadId:  aws.String(*initResult.UploadId),
				VaultName: aws.String(vault),
			})
			if err != nil {
				// TODO: queue the part to be reuploaded
				panic(formatAWSError(err))
			}
			bar.Add(contentLength)
			wg.Done()
		}(buf, startB, endB, hash)

		startB = endB
	}

	wg.Wait()

	input := &glacier.CompleteMultipartUploadInput{
		AccountId:   aws.String("-"),
		ArchiveSize: aws.String(fmt.Sprintf("%d", totalSize)),
		Checksum:    aws.String(th.Hash()),
		UploadId:    aws.String(*initResult.UploadId),
		VaultName:   aws.String(vault),
	}
	result, err := svc.CompleteMultipartUpload(input)
	if err != nil {
		return formatAWSError(err)
	}

	bar.Finish()

	// TODO: write the archive ID to a file
	fmt.Println(result)
	fmt.Println(*initResult.UploadId)

	return nil
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
