package pack

//MessageKickOut 消息
type MessageKickOut struct {
	MsgLen    uint32 //消息总长度
	Cmd       uint32 //消息ID长度
	ErrorCode uint32 //错误码
}

//NewMessageKickOut 创建一个Message消息包
func NewMessageKickOut(cmd, errcode uint32) *MessageKickOut {
	return &MessageKickOut{
		MsgLen:    12, //包总长度
		Cmd:       cmd,
		ErrorCode: errcode,
	}
}

//GetMsgLen 二进制流总长度
func (msg *MessageKickOut) GetMsgLen() uint32 {
	return msg.MsgLen
}
func (msg *MessageKickOut) SetMsgLen(len uint32) {
	msg.MsgLen = len
}

//GetCmd 获取cmd
func (msg *MessageKickOut) GetCmd() uint32 {
	return msg.Cmd
}
func (msg *MessageKickOut) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

//GetErrorCode ErrorCode数据
func (msg *MessageKickOut) GetErrorCode() uint32 {
	return msg.ErrorCode
}
func (msg *MessageKickOut) SetErrorCode(errcode uint32) {
	msg.ErrorCode = errcode
}
