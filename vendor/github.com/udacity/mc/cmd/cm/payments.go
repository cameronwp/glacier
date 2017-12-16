package cm

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/cm"
)

// flags
var (
	startdate string
	enddate   string
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Fetch payments information over a specified date/time interval",
	RunE: func(cmd *cobra.Command, args []string) error {
		staging, err := cmd.Flags().GetBool("staging")
		messagesFilename, err := cm.MessagePayments(startdate, enddate, staging)
		if err != nil {
			return err
		}
		fmt.Printf("\nA new file is available at %s!\n", messagesFilename)

		ratingsFilename, err := cm.RatingPayments(startdate, enddate, staging)
		if err != nil {
			return err
		}
		fmt.Printf("\nA new file is available at %s!\n", ratingsFilename)

		return nil
	},
}

func init() {
	paymentsCmd.Flags().StringVarP(&startdate, "startdate", "s", "", "Start Date (eg. 2017-08-01")
	err := paymentsCmd.MarkFlagRequired("startdate")
	if err != nil {
		log.Fatal(err)
	}

	paymentsCmd.Flags().StringVarP(&enddate, "enddate", "e", "", "End Date (eg. 2017-08-01)")
	err = paymentsCmd.MarkFlagRequired("enddate")
	if err != nil {
		log.Fatal(err)
	}
}
