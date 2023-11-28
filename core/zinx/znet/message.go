package znet

// Message 消息
type Message struct {
	MsgLen uint32 //消息:包总长度
	Cmd    uint32 //消息:cmd
	Data   []byte //消息:包体
}

// NewMessage 创建一个Message消息包
func NewMessage(cmd uint32, data []byte) *Message {
	return &Message{
		MsgLen: 4 + 4 + uint32(len(data)), //总长:uint32(4字节) + 消息ID:uint32(4字节) + len(data)
		Cmd:    cmd,
		Data:   data,
	}
}

// GetMsgLen 获取二进制流总长度
func (msg *Message) GetMsgLen() uint32 {
	return msg.MsgLen
}

// SetMsgLen 设置二进制流总长度
func (msg *Message) SetMsgLen(len uint32) {
	msg.MsgLen = len
}

// GetCmd 获取cmd数据
func (msg *Message) GetCmd() uint32 {
	return msg.Cmd
}

// SetCmd 设置cmd数据
func (msg *Message) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

// SetData 设置data数据
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

// GetData 获取data数据
func (msg *Message) GetData() []byte {
	return msg.Data
}
