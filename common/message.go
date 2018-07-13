package common

import(
	"github.com/satori/go.uuid"
	"github.com/gin-gonic/gin/json"
)

type Message struct {
	Id string `json:"message_id,omitempty"`
	Type MessageType
	Detail []byte
}

type GetInfo struct {
	BoxID string `json:"boxid,omitempty"`
}

type DetailRegister struct {
	BoxId string `json:"boxid,omitempty"`
	SI []byte `json:"sysinfo,omitempty"`
}

type DetailError struct {
	Code ErrorCode `json:"errcode,omitempty"`
	ErrorDetail string `json:"errdetail,omitempty"`
}

func NewMessageByDetail(t MessageType, datail []byte) (*Message, error) {
	var m = new(Message)
	ui, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	uis := ui.String()
	m.Id = uis
	m.Type = t
	m.Detail = datail
	return m, nil
}

func NewMessageByObj(t MessageType, obj interface{}) (*Message, error) {
	var m = new(Message)
	ui, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	uis := ui.String()
	m.Id = uis
	m.Type = t
	datail, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	m.Detail = datail
	return m, nil
}