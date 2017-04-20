package config

import (
	"encoding/json"
	"flag"
	"log"
	"io/ioutil"
)

type Config struct {
	Server        Server   `json:"server"`
	StaticDir     string   `json:"static_dir"`
	TemplateFiles []string `json:"template_files"`
	WsUrl         WsUrl    `json:"ws_url"`
}
type Server struct {
	Key	        string `json:"key"`
	Cert        string `json:"cert"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	CheckOrigin bool   `json:"check_origin"`
	WsPath      string `json:"ws_path"`
}

type Flag struct {
	ConfigFile  string
	Key         string
	Cert        string
	Host        string
	Port        int
	CheckOrigin bool
	WsUrl       WsUrl
}
type WsUrl struct {
	Ssl  bool   `json:"ssl"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Path string `json:"path"`
}

var f Flag

func loadDefaultConfig() *Config {
	return &Config{
		Server: Server{
			Key: "", Cert: "",
			Host: "localhost",
			Port: 8000,
			CheckOrigin: false,
			WsPath: "/ws",
		},
		StaticDir:     "static",
		TemplateFiles: []string{"templates/index.tmpl"},
		WsUrl:         WsUrl{Ssl: false, Host: "", Port: -1, Path: ""},
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

func LoadFlag() {
	flag.StringVar(&f.ConfigFile, "config", "", "config file in JSON format")
	flag.StringVar(&f.Key, "key", "", "server key")
	flag.StringVar(&f.Cert, "cert", "", "server cert")
	flag.StringVar(&f.Host, "host", "", "hostname")
	flag.BoolVar(&f.CheckOrigin, "checkorigin", false, "check origin for websocket")
	flag.IntVar(&f.Port, "port", -1, "port number")
	flag.BoolVar(&f.WsUrl.Ssl, "wsurl.ssl", false, "set ssl in `WsUrl` (template variable)")
	flag.StringVar(&f.WsUrl.Host, "wsurl.host", "", "set hostname in `WsUrl` (template variable)")
	flag.IntVar(&f.WsUrl.Port, "wsurl.port", -1, "set port in `WsUrl` (template variable)")
	flag.StringVar(&f.WsUrl.Path, "wsurl.path", "", "set path (starts with /) in `WsUrl` (template variable)")

	flag.Parse()
}

func LoadConfig() (*Config, error) {
	config := loadDefaultConfig()
	{
		j, err := loadFile(f.ConfigFile)
		if err != nil && f.ConfigFile != "" {
			return nil, err
		} else {
			json.Unmarshal(j, &config)
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
	if f.CheckOrigin == true {
		self.Server.CheckOrigin = f.CheckOrigin
	}

	if f.WsUrl.Ssl == true || self.Server.Key != "" {
		self.WsUrl.Ssl = true
	}

	if f.WsUrl.Host != "" {
		self.WsUrl.Host = f.WsUrl.Host
	} else if self.WsUrl.Host == "" {
		self.WsUrl.Host = self.Server.Host
	}

	if f.WsUrl.Port != -1 {
		self.WsUrl.Port = f.WsUrl.Port
	} else if self.WsUrl.Port == -1 {
		self.WsUrl.Port = self.Server.Port
	}

	if f.WsUrl.Path != "" {
		self.WsUrl.Path = f.WsUrl.Path
	} else if self.WsUrl.Path == "" {
		self.WsUrl.Path = self.Server.WsPath
	}
}
