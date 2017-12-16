package mentor

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/display"
	"github.com/udacity/mc/mentor"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a mentor by UID",
	Long: `The output from this command (a mentor's info) is
a table by default. You should widen your terminal to see the output displayed
correctly (or even legibly).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		isStaging, err := cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		m, err := mentor.FetchMentor(isStaging, uid)
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
	fetchCmd.Flags().StringVarP(&uid, "uid", "u", "", "Udacity UID")
	err := fetchCmd.MarkFlagRequired("uid")
	if err != nil {
		log.Fatal(err)
	}

	fetchCmd.Flags().BoolVarP(&outJSON, "json", "j", false, "Output JSON instead of a table")
}
