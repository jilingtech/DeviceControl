package Structs

type RegisterMessage struct {
	Id string `json:"id"`
}

//type用来标识返回信息的种类
// 1：状态上报的信息
// 2：执行命令后的响应
// 3: 异常关闭的消息
type ResponseMessage struct {
	Id string
	Type int
	Detail []byte
}

type RequestMessage struct {
	Id string
	Type int
	Detail []byte
}

type GetInfo struct {
	BoxID string `json:"boxid,omitempty"`
}