package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/udacity/mc/cmd/auth"
	"github.com/udacity/mc/cmd/cm"
	"github.com/udacity/mc/cmd/guru"
	"github.com/udacity/mc/cmd/livehelp"
	"github.com/udacity/mc/cmd/mentor"
	"github.com/udacity/mc/cmd/nd"
	"github.com/udacity/mc/cmd/reviews"
)

// RootCmd shows usage.
var RootCmd = &cobra.Command{
	Use:   "mc",
	Short: "Mentor CLI connects student services to our mentorship-related services",
	Long: `mc commands generally follow this pattern:

	mc [service] [verb] [--flag argument]

See below for more help. All commands run in production unless run with the
'--staging' flag.`,
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

// Ops has already been trained to use `mc login`. this provides a shortcut
var loginCmd = auth.LoginCmd

func init() {
	RootCmd.AddCommand(
		auth.Cmd,
		cm.Cmd,
		guru.Cmd,
		livehelp.Cmd,
		mentor.Cmd,
		nd.Cmd,
		reviews.Cmd,
		genDocsCmd,
		loginCmd,
	)

	RootCmd.PersistentFlags().Bool("staging", false, "Run against staging environments")
}
