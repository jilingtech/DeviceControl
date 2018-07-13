package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/bary321/DeviceControl/common"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id           string
	Socket       *websocket.Conn `json:"-"`
	Online       bool
	SysInfo      *common.SysInfo
	RegisterTime time.Time
	UpdateTime   time.Time
	DropTime     time.Time
	Send         chan []byte            `json:"-"`
	Missions     map[string]chan []byte `json:"-"`
}

func (c *Client) read() {
	defer func() {
		manager.Unregister <- c
		// c.Socket.Close()
	}()

	for {
		var rm = new(common.Message)
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
		if rm.Type == common.CommandResponseType {
			r, ok := c.Missions[rm.Id]
			if !ok {
				log.Println("Unknown message id")
				continue
			}
			r <- rm.Detail
		} else if rm.Type == common.StatusType {
			err = json.Unmarshal(rm.Detail, c.SysInfo)
			if err != nil {
				log.Println("can't unmarshal data", string(rm.Detail))
				continue
			}
			c.UpdateTime = time.Now()
		}
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
