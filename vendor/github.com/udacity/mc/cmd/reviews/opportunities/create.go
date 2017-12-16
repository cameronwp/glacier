package opportunities

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/udacity/mc/lang"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/udacity/mc/reviews"
)

// flags
var (
	infile        string
	days          int
	language      string
	noSubRequired bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create opportunities from a CSV. Note that flags are default values",
	Long: `This command attempts to generate opportunities based on the
information in the CSV file. If some piece of information is missing, it uses
the flags provided to fill in the gaps.

The CSV must include a "udacity_key" column.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		isStaging, err = cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		opportunities, err := parseCSV(infile)
		if err != nil {
			return err
		}

		fmt.Println("Finished reading CSV. Posting opportunities.")

		erred := postOpportunities(opportunities)
		if erred {
			fmt.Println("\nDone, but there were errors. See ~/.mc/errs.log")
		} else {
			fmt.Println("\nDone!")
		}

		return nil
	},
}

func init() {
	createCmd.Flags().IntVarP(&days, "days", "d", 7, "Number of days from now a candidate has to accept the opportunity. Defaults to 7 days")

	createCmd.Flags().StringVarP(&infile, "infile", "i", "", "Relative path to a CSV with candidate information")
	err := createCmd.MarkFlagRequired("infile")
	if err != nil {
		log.Fatal(err)
	}

	createCmd.Flags().StringVarP(&language, "language", "l", "", "Language for the opportunities")

	createCmd.Flags().BoolVarP(&noSubRequired, "no_submission_required", "n", false, "If set, means that candidates do not need a submission to create a cert. If not set, candidates need a passed submission against the specified project")

	createCmd.Flags().StringVarP(&projectID, "project_id", "p", "", "Project ID for which you want to create opportunities")
}

func parseCSV(filepath string) ([]reviews.Opportunity, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	opportunities := []reviews.Opportunity{}

	err = gocsv.UnmarshalFile(file, &opportunities)
	if err != nil {
		return nil, err
	}

	for i := range opportunities {
		setDefaultIfNecessary(&opportunities[i])
		opportunities[i].Language = lang.Normalize(opportunities[i].Language)
	}

	return opportunities, nil
}

func setDefaultIfNecessary(o *reviews.Opportunity) {
	pID, _ := strconv.Atoi(projectID)

	defaultO := reviews.Opportunity{
		Language:           language,
		ExpiresAt:          time.Now().AddDate(0, 0, days),
		ProjectID:          pID,
		SubmissionRequired: !noSubRequired,
	}

	emptyO := reviews.Opportunity{}

	if o.Language == emptyO.Language {
		o.Language = defaultO.Language
	}

	if o.ExpiresAt == emptyO.ExpiresAt {
		if o.Days != 0 {
			o.ExpiresAt = time.Now().AddDate(0, 0, o.Days)
		} else {
			o.ExpiresAt = defaultO.ExpiresAt
		}
	}

	if projectID != "" && o.ProjectID == emptyO.ProjectID {
		o.ProjectID = defaultO.ProjectID
	}

	if o.SubmissionRequired == emptyO.SubmissionRequired {
		o.SubmissionRequired = defaultO.SubmissionRequired
	}
}

func postOpportunities(opportunities []reviews.Opportunity) bool {
	erred := false
	count := len(opportunities)
	bar := pb.StartNew(count)

	for _, o := range opportunities {
		if o.ProjectID == 0 {
			err := fmt.Sprintf("could not find project ID for %s", o.UdacityKey)
			fmt.Println(err)
			log.Errorf(err)
			return true
		}

		err := reviews.PostOpportunity(isStaging, o)
		if err != nil {
			erred = true
			log.Error(err)
		}

		bar.Increment()
	}

	bar.Finish()

	return erred
}
