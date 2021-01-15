package websitecookieanalyzer

import (
	"encoding/csv"
	"io"
	"os"
)

func LoadWebsites(websitesInputFilepath string) ([]string, error) {
	file, err := os.Open(websitesInputFilepath)
	if err != nil {
		return nil, err
	}
	websites := make([]string, 0)
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 1
	reader.Comma = ','
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		websites = append(websites, record[0])
	}
	return websites, nil
}
