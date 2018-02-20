package cmd

import (
	"github.com/cameronwp/glacier/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all archives in a vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Render()

		return nil
	},
}

func init() {
	// listCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "AWS region of the vault")

	// listCmd.Flags().StringVarP(&vault, "vault", "v", "", "Vault name")
	// err := listCmd.MarkFlagRequired("vault")
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
