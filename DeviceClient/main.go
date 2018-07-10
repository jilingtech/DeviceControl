package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"encoding/json"
	)

type RegisterMessage struct {
	Id string `json:"id"`
}

type ResponseMessage struct {
	Id string
	Detail []byte
}

type CommandMessage struct {
	Id string

	Detail []byte
}

type Client struct {
	BoxId string
	Socket *websocket.Conn
	Send   chan []byte
}

var (
	host = flag.String("host", "192.168.2.90", "http service host")
	port = flag.Int("port", 1234, "http service port")
	id = flag.String("id", "test", "id")
)

func main() {
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", *host, *port), Path: "/ws"}
	var dialer *websocket.Dialer
	var c = new(Client)
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.BoxId = *id
	c.Socket = conn
	c.Send = make(chan []byte)

	go c.write()

	rm, err := json.Marshal(RegisterMessage{Id:c.BoxId})
	c.Send <- rm
	// go timeWriter(conn)

	for {
		var rm = new(ResponseMessage)
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
		var cm = new(CommandMessage)
		err = json.Unmarshal(message, cm)
		if err != nil {
			fmt.Println("covert err", err)
			continue
		}
		rm.Id = string(cm.Id)
		rm.Detail = cm.Detail
		rmj, err := json.Marshal(rm)
		c.Send <- rmj
	}
}

func (c *Client) write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}


func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}