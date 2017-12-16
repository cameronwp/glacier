package nd

import (
	"log"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/udacity/mc/gae"

	"github.com/spf13/cobra"
)

var unenrollCmd = &cobra.Command{
	Use:   "unenroll",
	Short: "Unnroll a mentor from nd050",
	RunE: func(cmd *cobra.Command, args []string) error {
		results := make(map[string]string)
		results[uid] = "maybe"

		for u := range results {
			unenrolled, err := gae.UnenrollFromMentorshipND(u)
			if err != nil {
				return err
			}

			if !unenrolled {
				enrolled, err := checkMentorNDEnrollment(uid)
				if err != nil {
					return err
				}
				unenrolled = !enrolled
			}

			if unenrolled {
				results[u] = "no"
			} else {
				results[u] = "yes"
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"udacity key", "nodekey", "enrolled"})
		for k, v := range results {
			table.Append([]string{k, "nd050", v})
		}

		table.Render()
		return nil
	},
}

func init() {
	unenrollCmd.Flags().StringVarP(&uid, "uid", "u", "", "Udacity UID")
	err := unenrollCmd.MarkFlagRequired("uid")
	if err != nil {
		log.Fatal(err)
	}
}
