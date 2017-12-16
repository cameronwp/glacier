package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/creds"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display who you are performing commands as",
	Long: `Shows both the email address and whether or not this user is being impersonated.
This command is useful when you are impersonating other Udacity staff members.`,
	Run: func(cmd *cobra.Command, args []string) {
		email, impersonating, err := creds.LoadI()
		if err != nil {
			fmt.Println("You are not logged in")
			return
		}

		impersonatingS := ""
		if impersonating {
			impersonatingS = fmt.Sprintf(" (impersonatee)")
		}

		fmt.Printf("You are %s%s\n", email, impersonatingS)
	},
}
