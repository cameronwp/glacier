package main

import (
	"os"

	"github.com/cameronwp/glacier/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
