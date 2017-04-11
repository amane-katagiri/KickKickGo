package config

import (
    "encoding/json"
    "flag"
    "io/ioutil"
)

type Config struct {
    Server Server `json:"server"`
    StaticDir string `json:"static_dir"`
    TemplateFiles []string `json:"template_files"`
}
type Server struct {
    Key string `json:"key"`
    Cert string `json:"cert"`
    Host string `json:"host"`
    Port int `json:"port"`
    WsPath string `json:"ws_path"`
}

type Flag struct {
    ConfigFile string
    Key string
    Cert string
    Host string
    Port int
}
var f Flag

func loadDefaultConfig() *Config {
    return &Config{
        Server: Server{
            Key: "", Cert: "",
            Host: "localhost",
            Port: 8000,
            WsPath: "/ws",
        },
        StaticDir: "static",
        TemplateFiles: []string{"templates/index.tmpl"},
    }
}

func loadFile(filename string) ([]byte, error) {
    if filename != "" {
        return ioutil.ReadFile(filename)
    } else {
        return ioutil.ReadFile("config.json")
    }
}

func LoadFlag() {
    flag.StringVar(&f.ConfigFile, "config", "", "config file in JSON format")
    flag.StringVar(&f.Key, "key", "", "server key")
    flag.StringVar(&f.Cert, "cert", "", "server cert")
    flag.StringVar(&f.Host, "host", "", "hostname")
    flag.IntVar(&f.Port, "port", -1, "port number")

    flag.Parse()
}

func LoadConfig() (*Config, error) {
    config := loadDefaultConfig()
    {
        j, err := loadFile(f.ConfigFile)
        if err != nil && f.ConfigFile != "" {
            return nil, err
        } else {
            json.Unmarshal(j, &config);
        }
    }
    config.update()
    return config, nil
}

func (self *Config) update() {
    if f.Key != "" {
        self.Server.Key = f.Key
    }
    if f.Cert != "" {
        self.Server.Cert = f.Cert
    }
    if f.Host != "" {
        self.Server.Host = f.Host
    }
    if f.Port != -1 {
        self.Server.Port = f.Port
    }
}
