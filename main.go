package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/amane-katagiri/kick-kick-go/config"
	"github.com/amane-katagiri/kick-kick-go/storage"
	"github.com/amane-katagiri/kick-kick-go/storage/null"
	"github.com/amane-katagiri/kick-kick-go/storage/redis"
	"github.com/amane-katagiri/kick-kick-go/websocket"
)

var wsURL = "wss?://host:port/path/to/ws"
var origin = "https?://host/path:port"
var secure = ""
var tmpl *template.Template

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index", map[string]interface{}{"WsUrl": wsURL})
}

func main() {
	storage.LoadFlag()
	config.LoadFlag()
	err := storage.LoadConfig()
	if err != nil {
		panic(err)
	}
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	defaultPort := 80
	if config.WsURL.Ssl {
		secure = "s"
		defaultPort = 443
	}
	if config.WsURL.Port != defaultPort {
		wsURL = fmt.Sprintf("ws%s://%s:%d%s", secure, config.WsURL.Host, config.WsURL.Port, config.WsURL.Path)
		origin = fmt.Sprintf("http%s://%s:%d", secure, config.WsURL.Host, config.WsURL.Port)
	} else {
		wsURL = fmt.Sprintf("ws%s://%s%s", secure, config.WsURL.Host, config.WsURL.Path)
		origin = fmt.Sprintf("http%s://%s", secure, config.WsURL.Host,)
	}
	if config.Server.CheckOrigin {
		websocket.SetOrigin(origin)
	}
	tmpl, err = template.New("").ParseFiles(config.TemplateFiles...)
	if err != nil {
		log.Println(err)
		tmpl, _ = template.New("err").Parse("template file is not found")
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc(config.Server.WsPath, websocket.WsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticDir))))

	var s storage.Storage
	s, err = redis.NewStorage()
	if err != nil {
		log.Println(err)
		s, err = null.NewStorage()
	}
	go websocket.ServeCount(s.GetCount(), s.SetCount)
	go websocket.ServeClients()
	go websocket.ServeID()

	if config.Server.Key != "" {
		log.Printf("Serving at https://%s:%d", config.Server.Host, config.Server.Port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Server.Port), config.Server.Key, config.Server.Cert, nil))
	} else {
		log.Printf("Serving at http://%s:%d", config.Server.Host, config.Server.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil))
	}
}
