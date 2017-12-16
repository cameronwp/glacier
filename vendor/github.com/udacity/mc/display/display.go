package display

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

// AsJSON marshals an interface and prints the result.
func AsJSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

// AsTable nicely displays data in a table for human consumption. PS: `in`
// should be []struct, otherwise this will panic. PPS: Note that the entire
// struct needs to be exportable, otherwise this will panic.
func AsTable(in interface{}) error {
	header, rows, err := generateTableInfo(in)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, row := range rows {
		table.Append(row)
	}

	table.Render() // Send output
	return nil
}

// Get the header and rows for the table.
func generateTableInfo(in interface{}) ([]string, [][]string, error) {
	header := []string{}
	rows := [][]string{}

	v := reflect.ValueOf(in)
	values := v.Slice(0, v.Len())

	// create rows
	for i := 0; i < values.Len(); i++ {
		r := values.Index(i)
		row := []string{}

		// create header with the first row
		if i == 0 {
			for j := 0; j < r.NumField(); j++ {
				field := r.Type().Field(j).Name
				header = append(header, field)
			}
		}

		// fill in the values for the row
		for j := 0; j < r.NumField(); j++ {
			val := r.Field(j).Interface()
			value := fmt.Sprintf("%v", val)
			row = append(row, value)
		}
		rows = append(rows, row)
	}

	return header, rows, nil
}
