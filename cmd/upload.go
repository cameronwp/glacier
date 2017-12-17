package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/spf13/cobra"
)

var (
	region string
	target string
	vault  string
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file or directory to Glacier",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			return err
		}

		files := make(map[string]struct{})
		getFiles(target, files)
		if len(files) == 0 {
			return fmt.Errorf("invalid target: no file(s) found")
		}

		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewSharedCredentials("", profile),
		})
		if err != nil {
			return err
		}

		_, err = sess.Config.Credentials.Get()
		if err != nil {
			return fmt.Errorf("AWS credentials error | %s", err)
		}

		svc := glacier.New(sess)

		for fp := range files {
			err = upload(svc, fp)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "AWS region of the vault")

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

func upload(svc *glacier.Glacier, fp string) error {
	partSize := int64(1 << 22)

	initResult, err := svc.InitiateMultipartUpload(&glacier.InitiateMultipartUploadInput{
		AccountId: aws.String("-"),
		PartSize:  aws.String(fmt.Sprintf("%d", partSize)), // 4MB part size
		VaultName: aws.String(vault),
	})
	if err != nil {
		return formatAWSError(err)
	}

	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	var totalSize int64
	if stats, err := f.Stat(); err != nil {
		totalSize = stats.Size()
	} else {
		return err
	}

	for i := int64(0); i < totalSize; i = i + partSize {
		start := i
		end := i + partSize - 1
		if end > totalSize {
			end = totalSize
		}

		uploadInput := &glacier.UploadMultipartPartInput{
			AccountId: aws.String("-"),
			Body:      f,
			Range:     aws.String(fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize)),
			UploadId:  aws.String(*initResult.UploadId),
			VaultName: aws.String(vault),
		}

		uploadResult, err := svc.UploadMultipartPart(uploadInput)
		if err != nil {
			return formatAWSError(err)
		}

		fmt.Println(uploadResult)
	}

	completeResult, err := svc.CompleteMultipartUpload(&glacier.CompleteMultipartUploadInput{
		AccountId:   aws.String("-"),
		ArchiveSize: aws.String(fmt.Sprintf("%d", totalSize)),
		UploadId:    aws.String(*initResult.UploadId),
		VaultName:   aws.String(vault),
	})

	fmt.Println(completeResult)

	if err != nil {
		return formatAWSError(err)
	}

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
