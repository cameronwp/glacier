package cmd

import (
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to Glacier",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
}
