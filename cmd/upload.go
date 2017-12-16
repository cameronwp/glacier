package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/spf13/cobra"
)

var (
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

		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String("us-west-2"),
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
		body, err := getBody()
		if err != nil {
			return err
		}

		result, err := svc.UploadArchive(&glacier.UploadArchiveInput{
			AccountId: aws.String("-"),
			VaultName: &vault,
			Body:      body,
		})
		if err != nil {
			return err
		}

		log.Println("Uploaded to archive", *result.ArchiveId)
		return nil
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&target, "target", "t", "", "Path to file or directory to upload")
	uploadCmd.Flags().StringVarP(&vault, "vault", "v", "", "Vault name")
}

func getBody() (io.ReadSeeker, error) {
	maybeFileOrDirectory, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("target cannot be found | %s", err)
	}

	if maybeFileOrDirectory.IsDir() {
		// do dir stuff
		return nil, nil
	}

	// just return the file
	return os.Open(target)
}
