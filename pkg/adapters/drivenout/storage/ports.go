package storage

type Port interface {
	SaveRecords(records [][]string)
	Records() [][]string
}
