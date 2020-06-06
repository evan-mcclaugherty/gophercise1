package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type FileAdapter struct {
	fileLocation string
}

func NewFileAdapter() *FileAdapter {
	return &FileAdapter{}
}

func (C *FileAdapter) GetRecords() ([][]string, error) {
	if C.fileLocation == "" {
		return nil, errors.New("no file location is set")
	}
	file, err := C.openFile(C.fileLocation)
	if err != nil {
		return nil, err
	}
	readFile, err := C.readFile(file)
	if err != nil {
		return nil, err
	}
	return readFile, nil
}

func (C *FileAdapter) SetFileLocation(fileLocation string) {
	C.fileLocation = fileLocation
}

func (C *FileAdapter) openFile(fileLocation string) (*os.File, error) {
	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	return file, nil
}

func (C *FileAdapter) readFile(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return records, nil
}
