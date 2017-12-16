package impersonate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/creds"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop impersonating and run commands as yourself",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := creds.StopImpersonating()
		if err != nil {
			return err
		}

		email, _, err := creds.Load()
		if err != nil {
			return err
		}

		fmt.Printf("Success! Now running commands as %s\n", email)
		return nil
	},
}
