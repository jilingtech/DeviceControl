package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Client struct {
	BoxId      string
	Socket     *websocket.Conn
	Send       chan []byte
	Registered bool
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

func (c *Client) WriteJson(obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	c.Send <- data
	return nil
}
