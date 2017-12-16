package impersonate

import (
	"github.com/spf13/cobra"
)

// Cmd is the root impersonation command.
var Cmd = &cobra.Command{
	Use:   "impersonate",
	Short: "Perform commands as a different Udacity staff member",
	Long: `Great for testing and troubleshooting. (You should probably ask the other
person before you impersonate them - it's the polite thing to do.)

CAVEAT: you cannot currently impersonate someone when fetching/updating ND
enrollments or Udacity email addresses. All commands will still work, but those
specific actions will be performed using your credentials, not the impersonatee.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	Cmd.AddCommand(
		startCmd,
		stopCmd,
	)
}
