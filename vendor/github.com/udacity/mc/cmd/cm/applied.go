package cm

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/cm"
	"github.com/udacity/mc/csv"
	"github.com/udacity/mc/mentor"
)

var appliedCmd = &cobra.Command{
	Use:   "applied",
	Short: "Fetch prospective classroom mentors who have applied for a specific Nanodegree",
	Long: `Fetches all classroom mentors from Classroom-Mentor API, fetches all
mentors from Mentor API, and, for those prospective mentors who applied for
Classroom Mentorship for the supplied ndkey and who are not already classroom
mentors for the supplied ndkey, fetches Classroom-Content info for those users
and outputs the records to a CSV file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		isStaging, err := cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		classMentors, err := cm.FetchClassMentors(isStaging)

		currentClassMentors := make(map[string]bool)
		for _, classMentor := range classMentors {
			currentClassMentors[classMentor.UID] = true
		}

		// Fetch mentors.
		mentors, err := mentor.FetchMentors(isStaging)
		if err != nil {
			return err
		}

		// Create rows.
		rowsChan := make(chan []string)
		var wg sync.WaitGroup
		for _, mentor := range mentors {
			wg.Add(1)
			go cm.CreateClassMentorAppliedRow(ndkey, currentClassMentors, mentor, &wg, rowsChan)
		}
		ticker := time.NewTicker(time.Second * 1)
		go func() {
			time.Sleep(time.Second * 1)
			fmt.Print("Processing")
			for range ticker.C {
				fmt.Print(".")
			}
		}()
		go func() {
			wg.Wait()
			close(rowsChan)
			ticker.Stop()
			fmt.Println()
		}()

		var rows [][]string
		for row := range rowsChan {
			rows = append(rows, row)
		}

		// Create CSV and check for errors.
		header := []string{"uid", "first_name", "last_name", "email", "nanodegrees", "paypal_email", "country", "languages", "bio", "educational_background", "intro_msg", "github_url", "linkedin_url", "avatar_url", "application", "created_at", "updated_at"}
		filename, err := csv.CreateCSVFile(cm.ClassMentorsByNanodegreeFilename, ndkey, header, rows)

		if err != nil {
			return err
		}

		fmt.Printf("\nA new file is available at %s!\n", filename)
		return nil
	},
}

func init() {
	appliedCmd.Flags().StringVarP(&ndkey, "ndkey", "k", "", "Nanodegree Key (eg. nd013)")
	err := appliedCmd.MarkFlagRequired("ndkey")
	if err != nil {
		log.Fatal(err)
	}
}
