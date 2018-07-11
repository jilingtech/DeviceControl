package main

import "time"

func ReportSysInfo(c *Client) {
	var sw = new(SysInfoWorker)
	for {
		time.Sleep(time.Duration(*delay) * time.Second)
		err := sw.Work(c)
		if err != nil {
			log.Error(err)
		}
	}
}
