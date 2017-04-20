package storage

import (
	"github.com/amane-katagiri/kick-kick-go/storage/redis"
)

// Storage is an interface of saving count backend
type Storage interface {
	GetCount() int
	SetCount(int)
}

// LoadFlag add flag parameters (call before config.LoadFlag)
func LoadFlag() {
	redis.LoadFlag()
}

// LoadConfig load config (call after LoadFlag)
func LoadConfig() error {
	err := redis.LoadConfig()
	if err != nil {
		return err
	}
	return nil
}
