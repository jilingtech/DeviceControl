package main

import (
	"flag"
	"fmt"
	"net/url"
	"github.com/gorilla/websocket"
	"encoding/json"
	"github.com/bary321/DeviceControl/Structs"
	"github.com/ipfs/go-ipfs-api"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	)

var (
	host  = flag.String("host", "localhost", "http service host")
	port  = flag.Int("port", 1234, "http service port")
	// id    = flag.String("id", "test", "id")
	delay = flag.Int("delay", 30, "delay")
	gateway = flag.String("gateway", "192.168.2.92:5001", "")
	log = logging.Logger("main")
)

func main() {
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", *host, *port), Path: "/ws"}
	var dialer *websocket.Dialer
	var c = new(Client)
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()
	var sh = shell.NewShell(*gateway)
	ido, err := sh.ID()
	if err != nil {
		log.Error(err)
		return
	}
	c.BoxId = ido.ID
	c.Socket = conn
	c.Send = make(chan []byte)

	go c.write()

	rm, err := json.Marshal(Structs.RegisterMessage{Id:c.BoxId})
	c.Send <- rm
	// go timeWriter(conn)

	for {
		var rm = new(Structs.ResponseMessage)
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
		var cm = new(Structs.RequestMessage)
		err = json.Unmarshal(message, cm)
		if err != nil {
			log.Error("covert err", err)
			continue
		}
		if cm.Type == 3 {
			log.Error(cm.Id)
			c.Socket.Close()
			return
		}
		rm.Id = string(cm.Id)
		rm.Detail = cm.Detail
		rmj, err := json.Marshal(rm)
		c.Send <- rmj
	}
}