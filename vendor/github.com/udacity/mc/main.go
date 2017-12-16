package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/udacity/mc/cmd"
	"github.com/udacity/mc/config"
)

func main() {
	err := createMcDir()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func createMcDir() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	mcpath := filepath.Join(usr.HomeDir, config.Dirname)
	_, err = os.Stat(mcpath)
	if os.IsNotExist(err) {
		err = os.Mkdir(mcpath, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
