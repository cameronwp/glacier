package cm

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/cm"
	"github.com/udacity/mc/display"
)

// flags
var (
	ndkey string
	uid   string
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch classroom mentors",
	Long: `Fetches classroom mentors from Classroom-Mentor API filtered by ndkey
or UID. Includes info for the mentor(s) from Classroom-Content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		isStaging, err := cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		if ndkey != "" {
			filename, err := cm.FetchByND(isStaging, ndkey)
			if err != nil {
				return err
			}

			fmt.Printf("\nA new file is available at %s!\n", filename)
		} else if uid != "" {
			cm, err := cm.FetchByUID(isStaging, uid)
			if err != nil {
				return err
			}

			return display.AsJSON(cm)
		} else {
			return fmt.Errorf("error: must specify either --uid or --ndkey")
		}

		return nil
	},
}

func init() {
	fetchCmd.Flags().StringVarP(&ndkey, "ndkey", "k", "", "Nanodegree key (eg. nd013)")
	fetchCmd.Flags().StringVarP(&uid, "uid", "u", "", "Udacity UID")
}
