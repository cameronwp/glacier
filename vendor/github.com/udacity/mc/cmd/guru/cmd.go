package guru

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/creds"
)

// shared flags
var (
	outJSON bool
	uid     string
)

// Cmd configures the root guru verb
var Cmd = &cobra.Command{
	Use:   "guru",
	Short: "Perform guru-specific tasks",
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
		createCmd,
		fetchCmd,
	)
}
