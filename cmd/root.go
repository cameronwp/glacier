package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// RootCmd shows usage.
var RootCmd = &cobra.Command{
	Use:   "glacier",
	Short: "Upload and manage files in AWS Glacier",
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
}
