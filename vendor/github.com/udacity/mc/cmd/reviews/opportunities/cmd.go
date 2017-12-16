package opportunities

import (
	"github.com/spf13/cobra"
	"github.com/udacity/mc/logging"
)

// flags for this package
var (
	isStaging bool
	projectID string
)

var log = logging.FileLogger()

// Cmd collects opportunities-related tasks.
var Cmd = &cobra.Command{
	Use:   "opportunities",
	Short: "Perform opportunities-related tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func init() {
	Cmd.AddCommand(
		candidatesCmd,
		createCmd,
	)
}
