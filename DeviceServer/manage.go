package main

import "log"

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
			if _, ok := manager.Clients[client.Id]; ok {
				close(client.Send)
				delete(manager.Clients, client.Id)
				log.Printf("A socket client is disconnected. %v", client)
			}
		}
	}
}