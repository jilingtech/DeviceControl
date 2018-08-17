package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bary321/DeviceControl/common"
	api "github.com/bary321/go-zfs-api"
	"github.com/gorilla/websocket"
	logs "gx/ipfs/QmQvJiADDe7JR4m968MwXobTCCzUqQkP87aRHe29MEBGHV/go-logging"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"net/url"
	"strings"
)

var (
	host = flag.String("host", "localhost", "http service host")
	port = flag.Int("port", 1234, "http service port")
	// id    = flag.String("id", "test", "id")
	delay    = flag.Int("delay", 3, "delay")
	gateway  = flag.String("gateway", "192.168.2.92:5001", "")
	loglevel = flag.String("loglevel", "DEBUG", "")
	log      = logging.Logger("main")
)

func main() {
	var lvl logs.Level
	var dialer *websocket.Dialer
	var c = new(Client)
	var dr = new(common.DetailRegister)

	flag.Parse()

	switch strings.ToUpper(*loglevel) {
	case "CRITICAL":
		lvl = logs.CRITICAL
	case "ERROR":
		lvl = logs.ERROR
	case "WARNING":
		lvl = logs.WARNING
	case "NOTICE":
		lvl = logs.NOTICE
	case "INFO":
		lvl = logs.INFO
	case "DEBUG":
		lvl = logs.DEBUG
	default:
		lvl = logs.WARNING
	}

	logging.SetAllLoggers(lvl)

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", *host, *port), Path: "/ws"}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	ido, err := api.GetID(*gateway)
	if err != nil {
		log.Error(err)
		return
	}
	c.BoxId = ido.ID

	log.Info("This box id is", c.BoxId)

	c.Socket = conn
	c.Send = make(chan []byte)
	c.Registered = false

	go c.write()

	sysinfo, err := api.GetSysInfo(*gateway)
	dr.BoxId = c.BoxId
	dr.SI = sysinfo
	rm, err := common.NewMessageByObj(common.RegisterType, dr)
	err = c.WriteJson(rm)
	if err != nil {
		log.Fatal(err)
		return
	}
	// go timeWriter(conn)

	go ReportSysInfo(c)

	for {
		var rec, res = new(common.Message), new(common.Message)
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			return
		}
		fmt.Printf("received: %s\n", message)
		err = json.Unmarshal(message, rec)
		if err != nil {
			log.Error("covert err", err)
			continue
		}

		if rec.Type == common.RegisterOkType {
			log.Info("Register OK")
			c.Registered = true
		} else if rec.Type == common.ErrorType {
			var de = new(common.DetailError)
			json.Unmarshal(rec.Detail, de)
			log.Error(de.Code, de.ErrorDetail)
			c.Socket.Close()
			return
		} else if rec.Type == common.CommandType {
			res.Id = string(rec.Id)
			res.Type = common.CommandResponseType
			res.Detail = rec.Detail
			err = c.WriteJson(res)
			if err != nil {
				log.Error(err)
				return
			}
		} else {
			log.Error("Error message type", rec.Type)
		}
	}
}
