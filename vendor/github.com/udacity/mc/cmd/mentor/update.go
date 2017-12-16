package mentor

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/display"
	"github.com/udacity/mc/mentor"
)

// flags
var (
	country     string
	language    string
	paypalEmail string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a mentor",
	Long: `Use the command line flags to change a mentor's information in mentor-api.
Only the flags you specify will be changed. Eg. if you do not include
'--language', the mentor's language won't change.

The output from this commands (a mentor's info) is a table by default. You
should widen your terminal to see the output displayed correctly (or even
legibly).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		isStaging, err := cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		fields := make(map[string]string)
		fields["uid"] = uid
		if country != "" {
			fields["country"] = country
		}

		if language != "" {
			fields["language"] = language
		}

		if paypalEmail != "" {
			fields["paypal_email"] = paypalEmail
		}

		m, err := mentor.UpdateMentor(isStaging, fields)
		if err != nil {
			return err
		}

		if m.UID == "" {
			fmt.Printf("no mentor found with UID '%s'\n", uid)
			return nil
		}

		if outJSON {
			return display.AsJSON(m)
		}

		return display.AsTable([]mentor.Mentor{m})
	},
}

func init() {
	updateCmd.Flags().StringVarP(&country, "country", "c", "", "Country code (ISO 3166-1 alpha-2 format)")

	updateCmd.Flags().StringVarP(&language, "language", "l", "", "Language (BCP 47 / INTL Team standards)")

	updateCmd.Flags().StringVarP(&paypalEmail, "paypal_email", "p", "", "Paypal email address")

	updateCmd.Flags().StringVarP(&uid, "uid", "u", "", "Udacity UID")
	err := updateCmd.MarkFlagRequired("uid")
	if err != nil {
		log.Fatal(err)
	}

	updateCmd.Flags().BoolVarP(&outJSON, "json", "j", false, "Output JSON instead of a table")
}
