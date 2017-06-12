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

var pushCount = make(chan int)
var clientCount = make(chan int)
var sockID = make(chan int)
var push = make(chan int)
var join = make(chan int)
var leave = make(chan int)

var socks = map[int]*ws.Conn{}
var socksLock = new(sync.RWMutex)

// ServeCounts increment/decrement websocket client count, increment pushed count, increment socket ID
func ServeCounts(getter func() int, setter func(int)) {
	total := getter()
	client := 0
	seq := 0
	for {
		select {
		case pushCount <- total:
		case clientCount <- client:
		case <-push:
			total++
			go setter(total)
		case <-join:
			client++
		case <-leave:
			client--
		case sockID <- seq:
			seq++
		}
	}
}

func send(conn *ws.Conn, message string) error {
	data, _ := json.Marshal(Message{Message: message, Count: <-pushCount, Clients: <-clientCount})
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
	join <- 0
	socksLock.Lock()
	socks[id] = c
}

func delSocket(c *ws.Conn, id int) {
	defer pole(id)
	defer func() { leave <- 0 }()
	defer socksLock.Unlock()
	c.Close()
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

	id := <-sockID
	defer delSocket(conn, id)

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
		push <- 0
		err = ping(id)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
