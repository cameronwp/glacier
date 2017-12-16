package guru

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/display"
	"github.com/udacity/mc/guru"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch a guru using a UID",
	RunE: func(cmd *cobra.Command, args []string) error {
		isStaging, err := cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		g, err := guru.FetchGuru(isStaging, uid)
		if err != nil {
			return err
		}

		if outJSON {
			return display.AsJSON(g)
		}

		return display.AsTable([]guru.RespGuru{g})
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
