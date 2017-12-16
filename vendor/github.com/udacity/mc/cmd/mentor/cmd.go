package mentor

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/creds"
)

// shared flags
var (
	uid     string
	outJSON bool
)

// Cmd configures the root mentor verb
var Cmd = &cobra.Command{
	Use:   "mentor",
	Short: "Collect and update mentor biographical and enrollment information",
	Long: `Mentors include people who have either applied for any kind of mentorship
and/or who are active as a mentor.`,
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
		fetchCmd,
		updateCmd,
	)
}
