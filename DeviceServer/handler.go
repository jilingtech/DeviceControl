package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bary321/DeviceControl/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Sockets(c *gin.Context) {
	var rm = new(common.Message)
	var rd = new(common.DetailRegister)

	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("read error ", err)
		return
	}

	fmt.Println("Rec: ", string(message))
	err = json.Unmarshal(message, rm)
	if err != nil {
		fmt.Println("Incorrect registration information")
		fmt.Println(string(message))
		return
	}

	err = json.Unmarshal(rm.Detail, rd)
	if err != nil {
		fmt.Println("Incorrect registration information")
		fmt.Println(string(rm.Detail))
		return
	}
	if client, ok := manager.Clients[rd.BoxId]; ok {
		if client.Online == false {
			client.Socket = conn
			client.UpdateTime = time.Now()
			client.Online = true
			client.Send = make(chan []byte)
			fmt.Println("a client reconnect", client)

			ro, _ := common.NewMessageByDetail(common.RegisterOkType, []byte{})
			conn.WriteJSON(ro)

			go client.read()
			go client.write()
		} else {
			fmt.Printf("%s is already registed\n", rd.BoxId)
			var de = new(common.DetailError)
			de.Code = common.DuplicateId
			de.ErrorDetail = fmt.Sprintf("%s is already registed", rd.BoxId)
			cm, err := common.NewMessageByObj(common.ErrorType, de)
			if err != nil {
				fmt.Println("Send message error", err)
				return
			}
			conn.WriteJSON(cm)
			conn.Close()
		}
	} else {
		client := &Client{Id: rd.BoxId, Socket: conn, Send: make(chan []byte), Missions: map[string]chan []byte{},
			Online: true, RegisterTime: time.Now(), UpdateTime: time.Now(), SysInfo: new(common.SysInfo)}
		manager.Register <- client

		ro, _ := common.NewMessageByDetail(common.RegisterOkType, []byte{})
		conn.WriteJSON(ro)

		go client.read()
		go client.write()
	}
}

func SendCommand(c *gin.Context) {
	var qos bool
	var timeout int
	var data []byte
	device_id, ok := c.GetQuery("device_id")
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
	device_id, ok := c.GetQuery("device_id")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("device_id not found"))
	}

	i, err := c.Request.Body.Read(buf)
	if err != io.EOF && err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "read body err")
		return
	}

	client, ok := manager.Clients[device_id]
	if !ok {
		c.String(http.StatusNotAcceptable, "Device not found")
		return
	}
	if !client.Online {
		c.String(http.StatusNotAcceptable, "Device registered but not online")
		return
	}

	cm, err := common.NewMessageByDetail(common.CommandType, buf[:i])
	data, err := json.Marshal(cm)
	if err != nil {
		c.String(http.StatusInternalServerError, "covert err")
		return
	}
	client.Missions[cm.Id] = make(chan []byte)
	client.Send <- data

	result := <-client.Missions[cm.Id]
	fmt.Println(string(result))
	delete(client.Missions, cm.Id)
	c.Data(http.StatusOK, "application/json", result)
}

func GetNodeInfo(c *gin.Context) {
	device_id, ok := c.GetQuery("device_id")
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
