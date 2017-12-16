package livehelp

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/udacity/mc/livehelp"
)

// flags
var (
	startdate string
	enddate   string
)

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		isStaging, err := cmd.Flags().GetBool("staging")
		if err != nil {
			return err
		}

		filename, err := livehelp.CalculatePayment(isStaging, startdate, enddate)
		if err != nil {
			return err
		}
		fmt.Printf("\nA new file is available at %s!\n", filename)

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
