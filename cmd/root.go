package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// shared flags
var (
	region string
	vault  string
	svc    *glacier.Glacier
)

// RootCmd shows usage.
var RootCmd = &cobra.Command{
	Use:   "glacier",
	Short: "Upload files to AWS Glacier",
	Long:  `If you want to use a non-default AWS profile from your credentials file, specify the profile you want to use with the --profile flag.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			return err
		}

		region, err := cmd.Flags().GetString("region")
		if err != nil {
			return err
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

		svc = glacier.New(sess)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

var genDocsCmd = &cobra.Command{
	Use:   "gen-docs",
	Short: "Generate the markdown documentation for the command tree",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Generating docs...")
		return doc.GenMarkdownTree(RootCmd, "./docs")
	},
}

func init() {
	RootCmd.AddCommand(
		inventoryCmd,
		uploadCmd,
		genDocsCmd,
	)

	RootCmd.PersistentFlags().String("profile", "default", "AWS credentials profile to use")
	RootCmd.PersistentFlags().String("region", "us-east-1", "AWS region of the vault")
}
