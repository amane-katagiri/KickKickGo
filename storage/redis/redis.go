package redis

import (
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "time"

    redigo "github.com/garyburd/redigo/redis"
)

type Config struct {
    Address string
    Key string
    MaxIdle int
    IdleTimeOut int
}
var config *Config

type Flag struct {
    ConfigFile string
    Address string
    Key string
    MaxIdle int
    IdleTimeOut int
}
var f Flag

func loadDefaultConfig() *Config {
    return &Config{
        Address: "localhost:6379",
        Key: "push_count",
        MaxIdle: 3,
        IdleTimeOut: 240,
    }
}

func loadFile(filename string) ([]byte, error) {
    if filename != "" {
        return ioutil.ReadFile(filename)
    } else {
        return ioutil.ReadFile("redis.json")
    }
}

func LoadFlag() {
    flag.StringVar(&f.ConfigFile, "redis.config", "", "config file in JSON format")
    flag.StringVar(&f.Address, "redis.address", "", "redis server addr 'host:port'")
    flag.StringVar(&f.Key, "redis.key", "", "your redis dict key")
    flag.IntVar(&f.MaxIdle, "redis.maxidle", -1, "redis.Pool.MaxIdle")
    flag.IntVar(&f.IdleTimeOut, "redis.idletimeout", -1, "redis.Pool.IdleTimeout (in sec)")
}

func LoadConfig() error {
    config = loadDefaultConfig()
    {
        j, err := loadFile(f.ConfigFile)
        if err != nil && f.ConfigFile != "" {
            return err
        } else {
            json.Unmarshal(j, &config);
        }
    }
    config.update()
    return nil
}

func (self *Config) update() {
    if f.Address != "" {
        self.Address = f.Address
    }
    if f.Key != "" {
        self.Key = f.Key
    }
    if f.MaxIdle != -1 {
        self.MaxIdle = f.MaxIdle
    }
    if f.IdleTimeOut != -1 {
        self.IdleTimeOut = f.IdleTimeOut
    }
}

type RedisStorage struct {
    key string
    pool redigo.Pool
}
func (s RedisStorage) GetCount() int {
    conn := s.pool.Get()
    defer conn.Close()
    i, err := redigo.Int(conn.Do("GET", s.key))
    if err != nil {
        log.Println(err)
        return 0
    }
    return i
}
func (s RedisStorage) SetCount(i int) {
    conn := s.pool.Get()
    defer conn.Close()
    _, err := conn.Do("SET", s.key, i)
    if err != nil {
        log.Println(err)
    }
}

func NewRedisStorage() *RedisStorage {
    return &RedisStorage{
        key: config.Key,
        pool: redigo.Pool{
            MaxIdle: config.MaxIdle,
            IdleTimeout: time.Duration(config.IdleTimeOut) * time.Second,
            Dial: func () (redigo.Conn, error) { return redigo.Dial("tcp", config.Address) },
        },
    }
}
