package csv

type Port interface {
	GetRecords() ([][]string, error)
}

type FilePort interface {
	Port
	SetFileLocation(fileLocation string)
}
