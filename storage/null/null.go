package null

import "log"

// Storage do nothing (count will be only in memory)
type Storage struct {
}

// GetCount return count as 0
func (s Storage) GetCount() int {
	return 0
}

// SetCount do nothing
func (s Storage) SetCount(i int) {
	log.Printf("count: %d", i)
	return
}

// NewStorage return new default Storage object
func NewStorage() (*Storage, error) {
	return &Storage{}, nil
}
