package pack

import "proxy/server/utils/pack/iface"

// MessageUp 消息
type MessageUp struct {
	MsgLen uint32 //消息总长度
	Cmd    uint32 //消息ID长度
	Data   []byte //消息全部的内容
}

// NewMessageUp 创建一个Message消息包
func NewMessageUp(cmd uint32, data []byte) iface.IMessageUp {
	msg := &MessageUp{
		MsgLen: 0,
		Cmd:    cmd,
		Data:   data,
	}

	//msg总长度: msgLen头(4) + cmd(4) + len(data)
	msg.MsgLen = 4 + 4 + uint32(len(data))

	return msg
}

// GetDataTotalLen 获取二进制流总长度
func (msg *MessageUp) GetMsgLen() uint32 {
	return msg.MsgLen
}

// SetDataTotalLen 设置消息数据长度
func (msg *MessageUp) SetMsgLen(len uint32) {
	msg.MsgLen = len
}

// GetCmd 获取cmd数据
func (msg *MessageUp) GetCmd() uint32 {
	return msg.Cmd
}

// SetCmd 设置cmd数据
func (msg *MessageUp) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

// GetData 获取消息内容
func (msg *MessageUp) GetData() []byte {
	return msg.Data
}

// SetData 设置消息内容
func (msg *MessageUp) SetData(head, body []byte) {
	msg.Data = append(head, body...)
}
