package pack

//MessageDown 消息
type MessageDown struct {
	MsgLen    uint32 //消息总长度
	Cmd       uint32 //消息ID长度
	ErrorCode uint32 //错误码

	PbLen      uint32 //pb长度
	Pb         []byte //protobuf内容
	BuffKeyLen uint32 //buff key长度
	BuffKey    []byte //buff key
}

//NewMsgPackageDown 创建一个Message消息包
func NewMessageDown(cmd, errcode uint32, pb, buffKey []byte) *MessageDown {
	msg := &MessageDown{
		MsgLen:     0,
		Cmd:        cmd,
		ErrorCode:  errcode,
		PbLen:      uint32(len(pb)),
		Pb:         pb,
		BuffKeyLen: uint32(len(buffKey)),
		BuffKey:    buffKey,
	}

	//msg总长度
	msg.MsgLen = 4 + 4 + 4 + 4 + msg.PbLen + 4 + msg.BuffKeyLen

	return msg
}

//GetMsgLen 二进制流总长度
func (msg *MessageDown) GetMsgLen() uint32 {
	return msg.MsgLen
}
func (msg *MessageDown) SetMsgLen(len uint32) {
	msg.MsgLen = len
}

//GetCmd 获取cmd
func (msg *MessageDown) GetCmd() uint32 {
	return msg.Cmd
}
func (msg *MessageDown) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

//GetErrorCode ErrorCode数据
func (msg *MessageDown) GetErrorCode() uint32 {
	return msg.ErrorCode
}
func (msg *MessageDown) SetErrorCode(errcode uint32) {
	msg.ErrorCode = errcode
}

//GetPbLen Pb数据
func (msg *MessageDown) GetPbLen() uint32 {
	return msg.PbLen
}
func (msg *MessageDown) GetPb() []byte {
	return msg.Pb
}

//GetBuffKeyLen Buff Key数据
func (msg *MessageDown) GetBuffKeyLen() uint32 {
	return msg.BuffKeyLen
}
func (msg *MessageDown) GetBuffKey() []byte {
	return msg.BuffKey
}
