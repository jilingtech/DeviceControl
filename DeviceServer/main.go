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
	r.GET("/ws", func(c *gin.Context) {
		Sockets(c.Writer, c.Request)
	})
	r.POST("/api/cmds/send", SendCommand)
	go manager.start()

	r.Run(fmt.Sprintf(":%d", *port))
}