package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

// CreateCSVFile takes a filename template, suffix, header, and data rows
// and creates a CSV file as a result, returning either the filename and nil or an empty string and relevant error.
func CreateCSVFile(filenameTemplate, suffix string, header []string, rows [][]string) (string, error) {
	filename := fmt.Sprintf(filenameTemplate, suffix)
	csvFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		err := csvFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.Write(header)
	if err != nil {
		return "", err
	}
	err = csvWriter.WriteAll(rows)
	if err != nil {
		return "", err
	}
	err = csvWriter.Error()
	if err != nil {
		return "", err
	}
	return filename, nil
}

// NewReader wraps encoding/csv.NewReader.
func NewReader(r io.Reader) *csv.Reader {
	return csv.NewReader(r)
}
