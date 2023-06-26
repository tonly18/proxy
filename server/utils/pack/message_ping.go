package pack

import (
	"proxy/server/utils/pack/iface"
)

//MessagePing 消息
type MessagePing struct {
	MsgLen uint32 //消息总长度
	Cmd    uint32 //消息ID长度
}

//NewMessagePing 创建一个Message消息包
func NewMessagePing(cmd uint32) iface.IMessagePing {
	msg := &MessagePing{
		MsgLen: 4 + 4, //msg总长度
		Cmd:    cmd,   //ping下发cmd
	}

	return msg
}

//GetMsgLen 获取二进制流总长度
func (msg *MessagePing) GetMsgLen() uint32 {
	return msg.MsgLen
}
func (msg *MessagePing) SetMsgLen(length uint32) {
	msg.MsgLen = length
}

//GetCmd 获取cmd数据
func (msg *MessagePing) GetCmd() uint32 {
	return msg.Cmd
}
func (msg *MessagePing) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}
