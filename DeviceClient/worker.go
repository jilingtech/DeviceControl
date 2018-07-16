package main

import (
	"encoding/json"
	"github.com/bary321/DeviceControl/common"
	api "github.com/bary321/go-zfs-api"
)

type Worker interface {
	Work()
}

type CameraWorker struct {
}

type SysInfoWorker struct {
}

func (s *SysInfoWorker) Work(c *Client) error {
	body, err := api.GetSysInfo(*gateway)
	if err != nil {
		return err
	}
	rm, err := common.NewMessageByDetail(common.StatusType, body)
	data, err := json.Marshal(rm)
	if err != nil {
		return err
	}
	c.Send <- data
	return nil
}
