package cameraid

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/bary321/DeviceControl/camera/format"
	"io/ioutil"
	"net/http"
)

type CamInfo struct {
	BoxId   string
	RtspUrl string
}

func GetInfo() *CamInfo {
	return &CamInfo{BoxId: "12345678901234567890", RtspUrl: "rtsp://192.168.2.90"}
}

func AddToPlatForm(rtsp, user, password, id, url string) (*format.RecMessage, error) {
	var m = new(format.RecMessage)
	client := &http.Client{}

	// data, err := format.SentXmlMarshal("rtsp://192.168.0.202:554/0", "admin", "123456", "5ef03573-2852-4457-b")
	data, err := format.SentXmlMarshal(rtsp, user, password, id)
	if err != nil {
		fmt.Println(err)
		return m, err
	}

	r := bytes.NewBuffer(data)

	// req, err := http.NewRequest("POST", "http://222.188.110.108:9999/device/addformzy", r)
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/xml")
	fmt.Println(req)
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return m, err
	}

	fmt.Println(xml.Unmarshal([]byte(body), &m))
	fmt.Println(m)
	fmt.Println(string(body))
	return m, nil
}
