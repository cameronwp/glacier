package main

import (
	"os"

	"github.com/udacity/mc/cmd"
)

func main() {
	err = cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
