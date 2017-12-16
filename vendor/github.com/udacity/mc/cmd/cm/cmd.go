package cm

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/creds"
)

// Cmd configures the root classmentor verb
var Cmd = &cobra.Command{
	Use:   "cm",
	Short: "Perform classroom mentor-specific tasks",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !creds.LoggedIn() {
			return fmt.Errorf("not logged in: run `mc login` to log in")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	Cmd.AddCommand(
		appliedCmd,
		fetchCmd,
		paymentsCmd,
	)
}
