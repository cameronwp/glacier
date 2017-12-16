package nd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/creds"
)

// shared flag
var uid string

// Cmd configures the root mentor verb
var Cmd = &cobra.Command{
	Use:   "nd",
	Short: "Collect and update (some) ND enrollment information for any user",
	Long: `Check the enrollment info for any student, mentor or not. Additionally,
you can enroll (and unenroll) mentors into nd050.

PS: the '--staging' flag has no effect here - all commands run against
production.

PPS: impersonation does not work here at the moment. You will always run commands as yourself.`,
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
		enrollCmd,
		enrollmentsCmd,
		unenrollCmd,
	)
}
