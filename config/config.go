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
	WsURL        wsURL    `json:"ws_url"` // WsURL is ignored from v0.5
}
type server struct {
	Key	        string `json:"key"`
	Cert        string `json:"cert"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	CheckOrigin bool   `json:"check_origin"`
	WsURL       string `json:"ws_url"`
	WsPath     string `json:"ws_path"` // WsPath is ignored from v0.5
}

type flagParam struct {
	ConfigFile  string
	Key         string
	Cert        string
	Host        string
	Port        int
	CheckOrigin bool
	WsURL       string
	_WsURL      wsURL // -wsurl is replaced from v0.5
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
			WsURL: "ws://localhost:8000/ws",
			WsPath: "",
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
	flag.IntVar(&f.Port, "port", -1, "port number")
	flag.StringVar(&f.WsURL, "wsurl", "", "websocket url")
	flag.BoolVar(&f.CheckOrigin, "checkorigin", false, "check origin for websocket")

	flag.BoolVar(&f._WsURL.Ssl, "wsurl.ssl", false, "(is ignored)")
	flag.StringVar(&f._WsURL.Host, "wsurl.host", "", "(is ignored)")
	flag.IntVar(&f._WsURL.Port, "wsurl.port", -1, "(is ignored)")
	flag.StringVar(&f._WsURL.Path, "wsurl.path", "", "(is ignored)")

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
		if config.Server.WsPath != "" || config.WsURL.Ssl == true || config.WsURL.Host != "" || config.WsURL.Port != -1 || config.WsURL.Path != "" {
			log.Println("`ws_url` and `server.ws_path` are ignored from v0.5. Use `server.ws_url` instead.")
		}
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
	if f.WsURL != "" {
		config.Server.WsURL = f.WsURL
	}

	if f._WsURL.Ssl == true || f._WsURL.Host != "" || f._WsURL.Port != -1 || f._WsURL.Path != "" {
		log.Println("`-wsurl.*` is ignored from v0.5. Use `-wsurl` instead.")
	}
}
