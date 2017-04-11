package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    ws "github.com/gorilla/websocket"
)

type Message struct {
    Message string `json:"message"`
    Count int `json:"count"`
    Clients int `json:"clients"`
}

var upgrader = ws.Upgrader{}

var add_count = make(chan int)
var get_count = make(chan int)
var add_clients = make(chan int)
var get_clients = make(chan int)
var seq = make(chan int)

var socks = map[int]*ws.Conn{}
var socksLock = new(sync.RWMutex)

// TODO: limit call setter frequency
func ServeCount(i int, setter func (int)) {
    for {
        select {
            case get_count <- i:
            case c := <-add_count:
                i += c
                go setter(i)
        }
    }
}

func ServeClients() {
    i := 0
    for {
        select {
            case get_clients <- i:
            case c := <-add_clients:
                i += c
        }
    }
}

func ServeId() {
    i := 0
    for {
        seq <- i
        i++
    }
}

func send(conn *ws.Conn, message string) error {
    data, _ := json.Marshal(Message{Message: message, Count: <-get_count, Clients: <-get_clients})
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

func WsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    id := <-seq
    defer pole(id)
    defer func() {add_clients <- -1}()
    defer delSocket(id)
    defer conn.Close()

    add_clients <- 1
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
        add_count <- 1
        err = ping(id)
        if err != nil {
            log.Println(err)
            break
        }
    }
}
