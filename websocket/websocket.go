package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sync"

	ws "github.com/gorilla/websocket"
)

// Message is format to send information
// Server send Message and client send "ping"
type Message struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
	Clients int    `json:"clients"`
}

var origin = ""

// SetOrigin set websocket origin "/scheme://host(:port)/" or "" (will not check)
func SetOrigin(org string) {
	origin = org
}

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		org := r.Header.Get("Origin")
		if len(org) == 0 {
			return true
		}
		u, err := url.Parse(org)
		if err != nil {
			return false
		}
		return origin == "" || org == origin || u.Host == r.Host
	},
}

var addCount = make(chan int)
var getCount = make(chan int)
var addClients = make(chan int)
var getClients = make(chan int)
var seq = make(chan int)

var socks = map[int]*ws.Conn{}
var socksLock = new(sync.RWMutex)

// ServeCount increment pushed count
func ServeCount(i int, setter func(int)) {
	for {
		select {
		case getCount <- i:
		case c := <-addCount:
			i += c
			go setter(i)
		}
	}
}

// ServeClients increment/decrement websocket client count
func ServeClients() {
	i := 0
	for {
		select {
		case getClients <- i:
		case c := <-addClients:
			i += c
		}
	}
}

// ServeID return new websocket client id
func ServeID() {
	i := 0
	for {
		seq <- i
		i++
	}
}

func send(conn *ws.Conn, message string) error {
	data, _ := json.Marshal(Message{Message: message, Count: <-getCount, Clients: <-getClients})
	err := conn.WriteMessage(ws.TextMessage, data)
	return err
}

func wait(conn *ws.Conn) error {
	_, _, err := conn.ReadMessage()
	return err
}

func sendAll(id int, pub string, prv string) error {
	defer socksLock.RUnlock()
	socksLock.RLock()
	for i, c := range socks {
		if i != id {
			send(c, pub)
		} else {
			err := send(c, prv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func addSocket(c *ws.Conn, id int) {
	defer socksLock.Unlock()
	socksLock.Lock()
	socks[id] = c
}

func delSocket(id int) {
	defer socksLock.Unlock()
	socksLock.Lock()
	delete(socks, id)
}

func pole(id int) error {
	return sendAll(id, "pole", "pole")
}

func ping(id int) error {
	return sendAll(id, "ping", "pong")
}

// WsHandler serve websocket
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := <-seq
	defer pole(id)
	defer func() { addClients <- -1 }()
	defer delSocket(id)
	defer conn.Close()

	addClients <- 1
	addSocket(conn, id)

	err = pole(id)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		err := wait(conn)
		if err != nil {
			log.Println(err)
			break
		}
		addCount <- 1
		err = ping(id)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
