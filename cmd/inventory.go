package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/spf13/cobra"
)

var sns string

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Trigger a vault inventory",
	Long:  `The inventory will be published to the given SNS.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initiateJob()
	},
}

func init() {
	inventoryCmd.Flags().StringVarP(&sns, "sns", "s", "", "SNS topic to publish to")
	err := inventoryCmd.MarkFlagRequired("sns")
	if err != nil {
		log.Fatal(err)
	}

	inventoryCmd.Flags().StringVarP(&vault, "vault", "v", "", "Vault name")
	err = inventoryCmd.MarkFlagRequired("vault")
	if err != nil {
		log.Fatal(err)
	}
}

func initiateJob() error {
	input := &glacier.InitiateJobInput{
		AccountId: aws.String("-"),
		JobParameters: &glacier.JobParameters{
			Description: aws.String("My inventory job"),
			Format:      aws.String("CSV"),
			SNSTopic:    aws.String(sns),
			Type:        aws.String("inventory-retrieval"),
		},
		VaultName: aws.String(vault),
	}

	result, err := svc.InitiateJob(input)
	if err != nil {
		return formatAWSError(err)
	}

	fmt.Println(result)
	return nil
}
