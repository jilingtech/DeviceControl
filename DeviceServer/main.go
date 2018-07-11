package main

import (
	"fmt"
	"flag"
	"github.com/gin-gonic/gin"
	)

var (
	port = flag.Int("port", 1234, "")
)

var manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}

func main() {
	flag.Parse()
	r := gin.Default()
	r.GET("/ws", Sockets)
	r.POST("/api/cmds/send", SendCommandBase)
	r.GET("/api/node/info", GetNodeInfo)

	go manager.start()

	r.Run(fmt.Sprintf(":%d", *port))
}