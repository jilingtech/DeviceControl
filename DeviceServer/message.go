package main

type RegisterMessage struct {
	Id string `json:"id"`
}

type ResponseMessage struct {
	Id string
	Detail []byte
}

type CommandMessage struct {
	Id string
	Detail []byte
}

type GetInfo struct {
	BoxID string `json:"boxid,omitempty"`
}