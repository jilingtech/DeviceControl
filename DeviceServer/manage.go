package main

import (
	"log"
	"time"
)

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func (manager *ClientManager) start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client.Id] = client
			log.Printf("A new socket client connected. %v", client)
		case client := <-manager.Unregister:
			if c, ok := manager.Clients[client.Id]; ok {
				// close(c.Send)
				//delete(manager.Clients, client.Id)
				c.Online = false
				c.Socket.Close()
				c.DropTime = time.Now()
				log.Printf("A socket client is disconnected. %v", client)
			}
		}
	}
}
