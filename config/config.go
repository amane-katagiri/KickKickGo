package config

import (
	"encoding/json"
	"flag"
	"log"
	"io/ioutil"
)

// Config is for kick-kick-go
type Config struct {
	Server        server   `json:"server"`
	StaticDir     string   `json:"static_dir"`
	TemplateFiles []string `json:"template_files"`
	WsURL         wsURL    `json:"ws_url"`
}
type server struct {
	Key	        string `json:"key"`
	Cert        string `json:"cert"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	CheckOrigin bool   `json:"check_origin"`
	WsPath      string `json:"ws_path"`
}

type flagParam struct {
	ConfigFile  string
	Key         string
	Cert        string
	Host        string
	Port        int
	CheckOrigin bool
	WsURL       wsURL
}
type wsURL struct {
	Ssl  bool   `json:"ssl"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Path string `json:"path"`
}

var f flagParam

func loadDefaultConfig() *Config {
	return &Config{
		Server: server{
			Key: "", Cert: "",
			Host: "localhost",
			Port: 8000,
			CheckOrigin: false,
			WsPath: "/ws",
		},
		StaticDir:     "static",
		TemplateFiles: []string{"templates/index.tmpl"},
		WsURL:         wsURL{Ssl: false, Host: "", Port: -1, Path: ""},
	}
}

func loadFile(filename string) ([]byte, error) {
	if filename != "" {
		return ioutil.ReadFile(filename)
	}
	b, err := ioutil.ReadFile("config/config.json")
	if err == nil {
		return b, err
	}
	b, err = ioutil.ReadFile("config.json")
	if err == nil {
		log.Println("default config file path is `config/config.json` from v0.4")
		return b, err
	}
	return nil, err
}

// LoadFlag add flag parameters (call after other LoadFlag functions)
func LoadFlag() {
	flag.StringVar(&f.ConfigFile, "config", "", "config file in JSON format")
	flag.StringVar(&f.Key, "key", "", "server key")
	flag.StringVar(&f.Cert, "cert", "", "server cert")
	flag.StringVar(&f.Host, "host", "", "hostname")
	flag.BoolVar(&f.CheckOrigin, "checkorigin", false, "check origin for websocket")
	flag.IntVar(&f.Port, "port", -1, "port number")
	flag.BoolVar(&f.WsURL.Ssl, "wsurl.ssl", false, "set ssl in `WsURL` (template variable)")
	flag.StringVar(&f.WsURL.Host, "wsurl.host", "", "set hostname in `WsURL` (template variable)")
	flag.IntVar(&f.WsURL.Port, "wsurl.port", -1, "set port in `WsURL` (template variable)")
	flag.StringVar(&f.WsURL.Path, "wsurl.path", "", "set path (starts with /) in `WsURL` (template variable)")

	flag.Parse()
}

// LoadConfig load config (call after LoadFlag)
func LoadConfig() (*Config, error) {
	config := loadDefaultConfig()
	{
		j, err := loadFile(f.ConfigFile)
		if err != nil && f.ConfigFile != "" {
			return nil, err
		}
		json.Unmarshal(j, &config)
	}
	config.update()
	return config, nil
}

func (config *Config) update() {
	if f.Key != "" {
		config.Server.Key = f.Key
	}
	if f.Cert != "" {
		config.Server.Cert = f.Cert
	}
	if f.Host != "" {
		config.Server.Host = f.Host
	}
	if f.Port != -1 {
		config.Server.Port = f.Port
	}
	if f.CheckOrigin == true {
		config.Server.CheckOrigin = f.CheckOrigin
	}

	if f.WsURL.Ssl == true || config.Server.Key != "" {
		config.WsURL.Ssl = true
	}

	if f.WsURL.Host != "" {
		config.WsURL.Host = f.WsURL.Host
	} else if config.WsURL.Host == "" {
		config.WsURL.Host = config.Server.Host
	}

	if f.WsURL.Port != -1 {
		config.WsURL.Port = f.WsURL.Port
	} else if config.WsURL.Port == -1 {
		config.WsURL.Port = config.Server.Port
	}

	if f.WsURL.Path != "" {
		config.WsURL.Path = f.WsURL.Path
	} else if config.WsURL.Path == "" {
		config.WsURL.Path = config.Server.WsPath
	}
}
