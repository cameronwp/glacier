package nd

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/gocarina/gocsv"
	"github.com/olekukonko/tablewriter"

	"github.com/udacity/mc/gae"

	"github.com/spf13/cobra"
)

// flag
var (
	infile string
)

var enrollCmd = &cobra.Command{
	Use:   "enroll",
	Short: "Enroll mentors in nd050",
	Long: `You can set the mentors you'd like to enroll with --uid and --infile. Feel
free to use both at the same time too. The --infile CSV must include a
'udacity_key' column. Currently enrolls students in v1.0.0 of nd050.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if infile == "" && uid == "" {
			return fmt.Errorf("error: you must set --infile and/or --uid")
		}

		results := make(map[string]string)
		defaultResult := "no"

		if infile != "" {
			uids, err := parseCSV(infile)
			if err != nil {
				return err
			}
			for _, u := range uids {
				results[u] = defaultResult
			}
		}

		if uid != "" {
			results[uid] = defaultResult
		}

		count := len(results)
		bar := pb.StartNew(count)
		for u := range results {
			enrolled, err := gae.EnrollInMentorshipND(u)
			if err != nil {
				return err
			}

			if !enrolled {
				enrolled, err = checkMentorNDEnrollment(u)
				if err != nil {
					return err
				}
			}

			if enrolled {
				results[u] = "yes"
			}

			bar.Increment()
		}

		bar.Finish()

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
	enrollCmd.Flags().StringVarP(&uid, "uid", "u", "", "Udacity UID")
	enrollCmd.Flags().StringVarP(&infile, "infile", "i", "", "Input CSV with 'udacity_key' column")
}

func parseCSV(filepath string) ([]string, error) {
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

	rows := []struct {
		UdacityKey string `csv:"udacity_key"`
	}{}

	err = gocsv.UnmarshalFile(file, &rows)
	if err != nil {
		return nil, err
	}

	uids := []string{}
	for _, r := range rows {
		uids = append(uids, r.UdacityKey)
	}
	return uids, nil
}

func checkMentorNDEnrollment(uid string) (bool, error) {
	enrollments, err := gae.FetchEnrollments(uid)
	if err != nil {
		return false, err
	}

	for _, e := range enrollments {
		if e.NodeKey == "nd050" && e.State == "enrolled" {
			return true, nil
		}
	}

	return false, nil
}
