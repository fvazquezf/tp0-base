package common

import (
	"encoding/csv"
	"io"
	"os"
)

type CsvReader struct {
	file   *os.File
	reader *csv.Reader
	end    int64
}

func NewCsvReader(filename string) (*CsvReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	end, _ := file.Seek(0, io.SeekEnd)
	file.Seek(0, io.SeekStart)
	return &CsvReader{
		file:   file,
		reader: reader,
		end:    end,
	}, nil
}

func (p *CsvReader) ReadLine() ([]string, error) {
	record, err := p.reader.Read()
	if err != nil {
		if err.Error() == "EOF" {
			return nil, err
		} else {
			return nil, err
		}
	}

	return record, nil
}

func (p *CsvReader) IsAtEnd() bool {
	pos, _ := p.file.Seek(0, io.SeekCurrent)
	// log.Infof("action: isAtEnd is true: current: %v, end: %v ", pos, p.end)
	return pos == p.end
}

func (p *CsvReader) Close() error {
	return p.file.Close()
}
