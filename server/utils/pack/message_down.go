package pack

import "proxy/server/utils/pack/iface"

// MessageDown 消息
type MessageDown struct {
	MsgLen uint32 //消息总长度
	Cmd    uint32 //消息ID长度
	Code   uint32 //错误码
	Data   []byte //protobuf内容
}

// NewMsgPackageDown 创建一个Message消息包
func NewMessageDown(cmd, code uint32, data []byte) iface.IMessageDown {
	msg := &MessageDown{
		MsgLen: 0,
		Cmd:    cmd,
		Code:   code,
		Data:   data,
	}

	//msg总长度
	msg.MsgLen = 4 + 4 + 4 + uint32(len(msg.Data))

	return msg
}

// GetMsgLen 包总长度
func (msg *MessageDown) GetMsgLen() uint32 {
	return msg.MsgLen
}
func (msg *MessageDown) SetMsgLen(len uint32) {
	msg.MsgLen = len
}

// GetCmd 获取cmd
func (msg *MessageDown) GetCmd() uint32 {
	return msg.Cmd
}
func (msg *MessageDown) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

// GetCode Code数据
func (msg *MessageDown) GetCode() uint32 {
	return msg.Code
}
func (msg *MessageDown) SetCode(code uint32) {
	msg.Code = code
}

// GetData Data数据
func (msg *MessageDown) GetData() []byte {
	return msg.Data
}
func (msg *MessageDown) SetData(data []byte) {
	msg.Data = data
}
