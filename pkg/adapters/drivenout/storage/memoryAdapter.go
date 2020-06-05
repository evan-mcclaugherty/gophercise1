package storage

import "sync"

type MemoryAdapter struct {
	csvMap [][]string
	sync.Mutex
}

func NewMemoryAdapter() *MemoryAdapter {
	return &MemoryAdapter{
		csvMap: [][]string{},
	}
}

func (m *MemoryAdapter) SaveRecords(records [][]string) {
	m.Lock()
	defer m.Unlock()
	m.csvMap = records
}

func (m *MemoryAdapter) Records() [][]string {
	return m.csvMap
}
