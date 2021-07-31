package utils

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
	"os"
)

// WriteOutput will write output to console with table | json | text | yaml format
func WriteOutput(headers []string, rows [][]string, outputFormat string) error {
	result := make([]interface{}, 0)
	for _, row := range rows {
		headerResult := make(map[string]string)
		for rowIndex, rowIndexValue := range row {
			headerResult[headers[rowIndex]] = rowIndexValue
		}
		result = append(result, headerResult)
	}

	switch {
	case outputFormat == "table":
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)
		table.AppendBulk(rows)
		table.Render()
	case outputFormat == "yaml", outputFormat == "text":
		out, err := yaml.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	case outputFormat == "json":
		out, err := json.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	default:
		out, err := json.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	}

	return nil
}
