package storage

import (
	"github.com/amane-katagiri/kick-kick-go/storage/redis"
)

type Storage interface {
	GetCount() int
	SetCount(int)
}

func LoadFlag() {
	redis.LoadFlag()
}

func LoadConfig() error {
	err := redis.LoadConfig()
	if err != nil {
		return err
	}
	return nil
}
