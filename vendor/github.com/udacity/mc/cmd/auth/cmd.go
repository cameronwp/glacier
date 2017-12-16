package auth

import (
	"github.com/spf13/cobra"
	"github.com/udacity/mc/cmd/auth/impersonate"
)

// Cmd configures the root auth verb
var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Perform authentication related commands (logging in, etc)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	Cmd.AddCommand(
		impersonate.Cmd,
		LoginCmd,
		logoutCmd,
		whoamiCmd,
	)
}
