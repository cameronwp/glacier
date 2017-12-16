package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/creds"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of mc and any other mentorship CLIs",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := creds.Logout()
		if err != nil {
			return fmt.Errorf("something went wrong logging out | %s", err)
		}

		fmt.Println("\nLogged out!")
		return nil
	},
}
