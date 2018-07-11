package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

type Worker interface {
	Work()
}

type CameraWorker struct {

}

type SysInfoWorker struct {

}

func (s *SysInfoWorker) Work(c *Client) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v0/diag/sys", *gateway))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.Send <- body
	return nil
}