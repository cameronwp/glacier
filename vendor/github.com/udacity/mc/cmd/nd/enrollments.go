package nd

import (
	"fmt"
	"log"

	"github.com/udacity/mc/display"

	"github.com/udacity/mc/gae"

	"github.com/spf13/cobra"
)

var (
	ndkey   string
	outJSON bool
)

var enrollmentsCmd = &cobra.Command{
	Use:   "enrollments",
	Short: "Get a user's ND enrollments (not just mentors)",
	Long: `If you want to search for a specific enrollment, you can use the --ndkey
flag, which will only show you a ND (or course) enrollment that matches the
ndkey.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		enrollments, err := gae.FetchEnrollments(uid)
		if err != nil {
			return err
		}

		filteredFor := ""
		if ndkey != "" {
			filterForND(&enrollments, ndkey)
			filteredFor = fmt.Sprintf(" filtered for %s", ndkey)
		}

		if len(enrollments) == 0 {
			fmt.Printf("no enrollments found for %s\n", uid)
			return nil
		}

		fmt.Printf("Enrollment(s) for %s%s\n", uid, filteredFor)
		if outJSON {
			return display.AsJSON(enrollments)
		}

		return display.AsTable(enrollments)
	},
}

func init() {
	enrollmentsCmd.Flags().StringVarP(&uid, "uid", "u", "", "Udacity UID")
	err := enrollmentsCmd.MarkFlagRequired("uid")
	if err != nil {
		log.Fatal(err)
	}

	enrollmentsCmd.Flags().StringVarP(&ndkey, "ndkey", "k", "", "ND key for an enrollment you want to check")

	enrollmentsCmd.Flags().BoolVarP(&outJSON, "json", "j", false, "Output JSON instead of a table")
}

func filterForND(enrollments *[]gae.Enrollment, ndkey string) {
	found := []gae.Enrollment{}
	for _, e := range *enrollments {
		if e.NodeKey == ndkey {
			found = append(found, e)
		}
	}
	*enrollments = found
}
