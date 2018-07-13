package common

//type用来标识返回信息的种类
// 0：注册信息
// 1：错误信息
// 2: 状态上报信息
// 3：命令下发使用的信息
// 4：命令的相应信息
type MessageType int

const (
	RegisterType MessageType = iota
	RegisterOkType
	ErrorType
	StatusType
	CommandType
	CommandResponseType
)

type ErrorCode int

const (
	DuplicateId ErrorCode = iota
)
