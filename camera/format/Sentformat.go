package format

import (
	"encoding/xml"
	"fmt"
)

type SentMessage struct {
	XMLName xml.Name `xml:"Message"`
	Boxinfo BoxInfo `xml:"BoxInfo"`
	DevicesCount DeviceCount `xml:"DeviceCount"`
}

type BoxInfo struct {
	XMLName xml.Name `xml:"BoxInfo"`
	Id string `xml:"ID,attr"`
}

type DeviceCount struct {
	XMLName xml.Name `xml:"DeviceCount"`
	Num string `xml:"num,attr"`
	Type string `xml:"type,attr"`
	Devices []Device `xml:"Device"`
}

type Device struct {
	XMLName xml.Name `xml:"Device"`
	Rtsp string `xml:"Rtsp,attr"`
	User string `xml:"User,attr"`
	PassWord string `xml:"PassWord,attr"`
}

func SentXmlMarshal(rtsp, user, password, id string) ([]byte, error) {
	var m = new(SentMessage)
	var bi = new(BoxInfo)
	var dc = new(DeviceCount)
	var dv = new(Device)
	/*
	dv.User = "admin"
	dv.Rtsp = "rtsp://192.168.2.90"
	dv.PassWord = "admin"
	bi.Id = "11111111111111111111"
	*/
	dv.User = user
	dv.PassWord = password
	bi.Id = id
	dv.Rtsp = rtsp
	dc.Devices = append(dc.Devices, *dv)
	dc.Type = "0"
	fmt.Println(len(dc.Devices))
	dc.Num = fmt.Sprintf("%d", len(dc.Devices))
	m.Boxinfo = *bi
	m.DevicesCount = *dc
	data, err := xml.MarshalIndent(m, "", "  ");
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}