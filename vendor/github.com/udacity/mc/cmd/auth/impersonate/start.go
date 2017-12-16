package impersonate

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/udacity/mc/creds"
	"github.com/udacity/mc/students"
)

var email string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Impersonate another Udacity staff member",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Pulling account info for '%s'...\n", email)
		users, err := students.Search(email)
		if err != nil {
			return err
		}

		var user students.User
		switch len(users) {
		case 0:
			fmt.Printf("Sorry, no users found who match '%s'\n", email)
			return nil
		case 1:
			user = users[0]
		default:
			fmt.Printf("Found multiple users who matched '%s'. Please check the accuracy of the email address.\n", email)
			return nil
		}

		fmt.Printf("Attempting to impersonate %s...\n", user.Email)
		err = creds.Impersonate(user.UID, user.Email)
		if err != nil {
			return err
		}

		fmt.Printf("Success! You are now impersonating %s\n", email)
		return nil
	},
}

func init() {
	startCmd.Flags().StringVarP(&email, "email", "u", "", "Udacity staff member email")
	err := startCmd.MarkFlagRequired("email")
	if err != nil {
		log.Fatal(err)
	}
}
