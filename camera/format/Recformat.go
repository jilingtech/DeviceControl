package format

import (
	"encoding/xml"
)

type RecMessage struct {
	XMLName      xml.Name       `xml:"Message"`
	Code         string         `xml:"Code"`
	Descript     string         `xml:"Descript"`
	DevicesCount RecDeviceCount `xml:"DeviceCount"`
}

type RecDeviceCount struct {
	XMLName xml.Name    `xml:"DeviceCount"`
	Num     string      `xml:"num,attr"`
	Type    string      `xml:"type,attr"`
	Devices []RecDevice `xml:"Device"`
}

type RecDevice struct {
	XMLName xml.Name `xml:"Device"`
	Id      string   `xml:"ID,attr"`
}

func RecXmlMarshal(data []byte) (*RecMessage, error) {
	var m = new(RecMessage)
	xml.Unmarshal([]byte(data), &m)
	return m, nil
}
