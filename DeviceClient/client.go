package main

import (
	"github.com/gorilla/websocket"
				)

type Client struct {
	BoxId string
	Socket *websocket.Conn
	Send   chan []byte
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
