package guru

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/display"
	"github.com/udacity/mc/guru"
)

var ndkeys []string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a guru using a UID and list of 1+ NDs",
	RunE: func(cmd *cobra.Command, args []string) error {
		isStaging, err := cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		g, err := guru.CreateGuru(isStaging, uid, ndkeys)
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
	createCmd.Flags().StringVarP(&uid, "uid", "u", "", "Udacity UID")
	err := createCmd.MarkFlagRequired("uid")
	if err != nil {
		log.Fatal(err)
	}

	createCmd.Flags().StringSliceVarP(&ndkeys, "ndkeys", "n", nil, "List of ND Keys")
	err = createCmd.MarkFlagRequired("ndkeys")
	if err != nil {
		log.Fatal(err)
	}

	createCmd.Flags().BoolVarP(&outJSON, "json", "j", false, "Output JSON instead of a table")
}
