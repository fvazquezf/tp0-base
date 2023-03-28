package common

import (
	"encoding/csv"
	"os"
)

type CsvReader struct {
	file   *os.File
	reader *csv.Reader
}

func NewCsvReader(filename string) (*CsvReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	return &CsvReader{
		file:   file,
		reader: reader,
	}, nil
}

func (p *CsvReader) ReadLine() ([]string, error) {
	record, err := p.reader.Read()
	if err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return record, nil
}

func (p *CsvReader) Close() error {
	return p.file.Close()
}