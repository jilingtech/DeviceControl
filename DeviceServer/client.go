package main

import (
	"log"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
	Missions   map[string]chan []byte
}

func (c *Client) read() {
	defer func() {
		manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		var rm = new(ResponseMessage)
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			manager.Unregister <- c
			c.Socket.Close()
			break
		}
		log.Printf("recv: %s", message)
		err = json.Unmarshal(message, rm)
		if err != nil {
			log.Println("Unformatted data")
			continue
		}
		r, ok := c.Missions[rm.Id]
		if !ok {
			log.Println("Unknown message id")
			continue
		}
		r <- rm.Detail
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