package auth

import (
	"fmt"

	"github.com/udacity/mc/creds"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// LoginCmd logs in. Exported to be used with `mc login` as well as `mc auth
// login`.
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Udacity (required for basically every command)",
	Long: `Use your Udacity staff email and password. This login command may also
authenticate you for other mentorship CLIs. Note: 'mc login' is a shortcut for
'mc auth login' (they are the same command).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Udacity Email: ")
		var email string
		_, err := fmt.Scanf("%s", &email)
		if err != nil {
			return err
		}

		fmt.Print("Udacity Password: ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			return err
		}

		if len(password) == 0 {
			return fmt.Errorf("no password entered")
		}

		fmt.Println("\nLogging in...")
		err = creds.Login(email, string(password))
		if err != nil {
			return fmt.Errorf("cannot sign in (check your password) | %s", err)
		}

		fmt.Println("\nLogged in!")
		return nil
	},
}
