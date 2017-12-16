package cm

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func checkFileExists(filename, startDate, endDate string) error {
	dateFormat := "2006-01-02"
	endDateTime, err := time.Parse(dateFormat, endDate)
	if err != nil {
		return err
	}
	// Check for file existence and age.
	// If file does not exist or the file is too old (i.e., it was generated before
	// or on the endDate), attempt to download a new version.
	fileInfo, err := os.Stat(filename)
	if err != nil || fileInfo.ModTime().Before(endDateTime.Add(time.Hour*24)) {
		if guruSQL[filename] != "" {
			fmt.Printf("\nDownload '%s' using your favorite DB tool and the following query:\n%s", filename, fmt.Sprintf(guruSQL[filename], startDate, endDate))
		} else {
			URL := fmt.Sprintf(mentorshipSegmentDashboardURL, startDate, endDate)
			// Depending on user settings, this command is able to open
			// the default browser and trigger a download.
			cmd := exec.Command("open", URL)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				return err
			}
			if stderr.String() != "" {
				return errors.New(stderr.String())
			}
		}
		return fmt.Errorf("move '%s' to the current directory", filename)
	}
	return nil
}

func readCSV(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("error: closing file '%s': %s", filename, err)
		}
	}()
	return csv.NewReader(f).ReadAll()
}
