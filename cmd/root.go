package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// shared flags
var (
	region string
	vault  string
)

// RootCmd shows usage.
var RootCmd = &cobra.Command{
	Use:   "glacier",
	Short: "Upload files to AWS Glacier",
	Long: `If you want to use a non-default AWS profile from your credentials file, specify
the profile you want to use with the --profile flag.`,
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
		uploadCmd,
		genDocsCmd,
	)

	RootCmd.PersistentFlags().String("profile", "default", "AWS credentials profile to use")
}
