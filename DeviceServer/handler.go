package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"errors"
	"strconv"
	"github.com/satori/go.uuid"
	"io"
	"github.com/bary321/DeviceControl/Structs"
	"time"
	)

func Sockets(c *gin.Context) {
	var rm = new(Structs.RegisterMessage)
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("write error ", err)
		return
	}

	fmt.Println(string(message))
	err = json.Unmarshal(message, rm)
	if err != nil {
		fmt.Println("Incorrect registration information")
		fmt.Println(string(message))
		return
	}
	if client, ok := manager.Clients[rm.Id]; ok {
		if client.Online == false {
			client.Socket = conn
			client.UpdateTime = time.Now()
			client.Online = true
			client.Send = make(chan []byte)
			fmt.Println("a client reconnect", client)
			go client.read()
			go client.write()
		} else {
			fmt.Println("id is already in use", rm.Id)
			var cm = new(Structs.RequestMessage)
			cm.Id = "__replicate_id__"
			cm.Type = 3
			conn.WriteJSON(cm)
			conn.Close()
			c.AbortWithError(http.StatusBadRequest, errors.New("replicate id"))
			return
		}
	} else {
		client := &Client{Id: rm.Id, Socket: conn, Send: make(chan []byte), Missions: map[string]chan []byte{},
						Online:true, RegisterTime:time.Now(), UpdateTime:time.Now(), SysInfo:new(Structs.SysInfo)}
		manager.Register <- client
		go client.read()
		go client.write()
	}
}

func SendCommand(c *gin.Context) {
	var qos bool
	var timeout int
	var data []byte
	device_id, ok:= c.GetQuery("device_id")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("device_id not found"))
	}
	url_qos, ok := c.GetQuery("qos")
	if !ok {
		url_qos = "1"
	}
	if url_qos == "0" {
		qos = false
	} else if url_qos == "1" {
		qos = true
	} else {
		c.AbortWithError(http.StatusBadRequest, errors.New("qos value not allow"))
	}
	url_timeout, ok := c.GetQuery("timeout")
	if !ok {
		url_timeout = "0"
	}
	timeout, err := strconv.Atoi(url_timeout)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("timeout value not allow"))
	}
	client, ok := manager.Clients[device_id]
	if !ok {
		c.String(http.StatusNotAcceptable, "Device not found")
		return
	}
	_, err = c.Request.Body.Read(data)
	if !qos {
		client.Send <- data
		c.String(http.StatusOK, "ok")
	} else {
		fmt.Println(timeout)
	}
}

func SendCommandBase(c *gin.Context) {
	buf := make([]byte, 1024)
	device_id, ok:= c.GetQuery("device_id")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("device_id not found"))
	}
	var cm = new(Structs.RequestMessage)
	client, ok := manager.Clients[device_id]
	if !ok {
		c.String(http.StatusNotAcceptable, "Device not found")
		return
	}
	if ! client.Online {
		c.String(http.StatusNotAcceptable, "Device registered but not online")
		return
	}
	ui, _ := uuid.NewV4()
	uis := ui.String()
	cm.Id = uis

 	i, err := c.Request.Body.Read(buf)
	if err != io.EOF && err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "read body err")
		return
	}
 	cm.Detail = buf[:i]
	fmt.Println(cm.Detail)
	data, err := json.Marshal(cm)
	if err != nil {
		c.String(http.StatusInternalServerError, "covert err")
		return
	}
	client.Send <- data
	client.Missions[uis] = make(chan []byte)
	result := <- client.Missions[uis]
	fmt.Println(string(result))
	delete(client.Missions, uis)
	c.Data(http.StatusOK, "application/json", result)
}

func GetNodeInfo(c *gin.Context) {
	device_id, ok:= c.GetQuery("device_id")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("device_id not found"))
	}
	client, ok := manager.Clients[device_id]
	if !ok {
		c.String(http.StatusNotAcceptable, "Device not found")
		return
	}
	c.JSON(http.StatusOK, client)
}