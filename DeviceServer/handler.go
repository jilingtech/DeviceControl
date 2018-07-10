package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
		"github.com/gorilla/websocket"
	"errors"
	"strconv"
)

func Sockets(res http.ResponseWriter, req *http.Request) {
	var rm = new(RegisterMessage)
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
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

	client := &Client{Id: rm.Id, Socket: conn, Send: make(chan []byte), Missions: map[string]chan []byte{}}

	manager.Register <- client

	go client.read()
	go client.write()
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

	}

	/*
	var reqInfo GetInfo
	err := c.BindJSON(&reqInfo)
	if err != nil {
		c.String(http.StatusBadRequest, "error")
		return
	} else {
		var cm = new(CommandMessage)
		client, ok := manager.Clients[reqInfo.BoxID]
		if !ok {
			c.String(http.StatusNotAcceptable, "Device not found")
			return
		}

		ui, _ := uuid.NewV4()
		fmt.Println(ui)
		uis := ui.String()
		cm.Id = uis
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
		c.String(http.StatusOK, "ok")
	}
	*/
}