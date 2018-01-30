package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all archives in a vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		// profile, err := cmd.Flags().GetString("profile")
		// if err != nil {
		// 	return err
		// }

		// sess, err := session.NewSession(&aws.Config{
		// 	Region:      aws.String(region),
		// 	Credentials: credentials.NewSharedCredentials("", profile),
		// })
		// if err != nil {
		// 	return err
		// }

		// _, err = sess.Config.Credentials.Get()
		// if err != nil {
		// 	return fmt.Errorf("AWS credentials error | %s", err)
		// }

		// svc := glacier.New(sess)

		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "AWS region of the vault")

	listCmd.Flags().StringVarP(&vault, "vault", "v", "", "Vault name")
	err := listCmd.MarkFlagRequired("vault")
	if err != nil {
		log.Fatal(err)
	}
}
