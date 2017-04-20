package redis

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

// Config is for redis
type Config struct {
	Address     string
	Key         string
	MaxIdle     int
	IdleTimeOut int
}

var config *Config

// Flag has flag parameters for redis
type flagParam struct {
	ConfigFile  string
	Address     string
	Key         string
	MaxIdle     int
	IdleTimeOut int
}

var f flagParam

func loadDefaultConfig() *Config {
	return &Config{
		Address:     "",
		Key:         "push_count",
		MaxIdle:     3,
		IdleTimeOut: 240,
	}
}

func loadFile(filename string) ([]byte, error) {
	if filename != "" {
		return ioutil.ReadFile(filename)
	}
	b, err := ioutil.ReadFile("config/redis.json")
	if err == nil {
		return b, err
	}
	b, err = ioutil.ReadFile("redis.json")
	if err == nil {
		log.Println("default redis config file path is `config/redis.json` from v0.4")
		return b, err
	}
	return nil, err
}

// LoadFlag add flag parameters (call from storage.LoadFlag)
func LoadFlag() {
	flag.StringVar(&f.ConfigFile, "redis.config", "", "config file in JSON format")
	flag.StringVar(&f.Address, "redis.address", "", "redis server addr 'host:port'")
	flag.StringVar(&f.Key, "redis.key", "", "your redis dict key")
	flag.IntVar(&f.MaxIdle, "redis.maxidle", -1, "redis.Pool.MaxIdle")
	flag.IntVar(&f.IdleTimeOut, "redis.idletimeout", -1, "redis.Pool.IdleTimeout (in sec)")
}

// LoadConfig load config (call from LoadConfig)
func LoadConfig() error {
	config = loadDefaultConfig()
	j, err := loadFile(f.ConfigFile)
	if err != nil && f.ConfigFile != "" {
		return err
	}
	json.Unmarshal(j, &config)
	config.update()
	return nil
}

func (config *Config) update() {
	if f.Address != "" {
		config.Address = f.Address
	}
	if f.Key != "" {
		config.Key = f.Key
	}
	if f.MaxIdle != -1 {
		config.MaxIdle = f.MaxIdle
	}
	if f.IdleTimeOut != -1 {
		config.IdleTimeOut = f.IdleTimeOut
	}
}

// Storage save count to redis
type Storage struct {
	key  string
	pool redigo.Pool
}

// GetCount request count from redis
func (s Storage) GetCount() int {
	conn := s.pool.Get()
	defer conn.Close()
	i, err := redigo.Int(conn.Do("GET", s.key))
	if err != nil {
		log.Println(err)
		return 0
	}
	return i
}

// SetCount store count to redis
func (s Storage) SetCount(i int) {
	conn := s.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", s.key, i)
	if err != nil {
		log.Println(err)
	}
}

// NewStorage return new Storage object from config
func NewStorage() (*Storage, error) {
	if config.Address != "" {
		return &Storage{
			key: config.Key,
			pool: redigo.Pool{
				MaxIdle:     config.MaxIdle,
				IdleTimeout: time.Duration(config.IdleTimeOut) * time.Second,
				Dial:        func() (redigo.Conn, error) { return redigo.Dial("tcp", config.Address) },
			},
		}, nil
	}
	return nil, errors.New("`Address` for RedisStorage is not specified")
}
