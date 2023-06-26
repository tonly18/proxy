package pack

import "proxy/server/utils/pack/iface"

//MessageUp 消息
type MessageUp struct {
	MsgLen uint32 //消息总长度
	Cmd    uint32 //消息ID长度

	DataLen uint32 //data长度 = mode(1) + uin头(2) + uin数据长度(UinLen) + pb数据长度(pbLen)

	Mode uint8 //Mode 固定值为0

	UinLen uint16 //uin长度
	Uin    []byte //uin内容:string

	Pb []byte //protobuf内容

	KeyLen uint32 //Key的长度
	Key    []byte //Key内容:string

	DebugLen uint32   //debug长度
	Debug    []uint64 //debug内容

	TroopLen uint32   //troop长度
	Troop    []uint32 //军队内容

	Data    []byte //消息全部的内容
	RawData []byte //消息全部的内容(二进制流)
}

//NewMessageUp 创建一个Message消息包
func NewMessageUp(cmd uint32, pb []byte, uin []byte, key []byte, debug []uint64, troop []uint32) iface.IMessageUp {
	msg := &MessageUp{
		MsgLen: 0,
		Cmd:    cmd,

		Mode:    uint8(0), //固定值:0
		DataLen: 0,

		UinLen: uint16(len(uin)),
		Uin:    uin,

		Pb: pb,

		KeyLen: uint32(len(key)),
		Key:    key,

		DebugLen: 8 * uint32(len(debug)),
		Debug:    debug,

		TroopLen: 4 * uint32(len(troop)),
		Troop:    troop,

		Data:    nil,
		RawData: nil,
	}

	//dataLen = mode(1) + uin头(2) + UinLen + pbLen
	msg.DataLen = 1 + 2 + uint32(msg.UinLen) + uint32(len(msg.Pb))

	//msg总长度: msgLen头(4) + cmd(4) + DataLen头(4) + DataLen
	msg.MsgLen = 4 + 4 + 4 + msg.DataLen + (4 + msg.KeyLen) + (4 + msg.DebugLen) + (4 + msg.TroopLen)

	return msg
}

//GetDataToalLen 获取二进制流总长度
func (msg *MessageUp) GetMsgLen() uint32 {
	return msg.MsgLen
}

//SetDataToalLen 设置消息数据长度
func (msg *MessageUp) SetMsgLen(len uint32) {
	msg.MsgLen = len
}

//GetCmd 获取cmd数据
func (msg *MessageUp) GetCmd() uint32 {
	return msg.Cmd
}

//SetCmd 设置cmd数据
func (msg *MessageUp) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

//GetMode 获取cmd数据
func (msg *MessageUp) GetMode() uint8 {
	return msg.Mode
}

//SetMode 设置cmd数据
func (msg *MessageUp) SetMode(mode uint8) {
	msg.Mode = mode
}

//GetDataLen 获取数据长度
func (msg *MessageUp) GetDataLen() uint32 {
	return msg.DataLen
}

//SetDataLen 获取pb数据长度
func (msg *MessageUp) SetDataLen(length uint32) {
	msg.DataLen = length
}

//GetUinLen 获取uin数据长度
func (msg *MessageUp) GetUinLen() uint16 {
	return msg.UinLen
}

//GetUin 获取uin数据
func (msg *MessageUp) GetUin() []byte {
	return msg.Uin
}

//SetUinLen 获取uin数据长度
func (msg *MessageUp) SetUinLen(length uint16) {
	msg.UinLen = length
}

//SetUin 获取uin数据
func (msg *MessageUp) SetUin(uin []byte) {
	msg.Uin = uin
}

//GetPb 获取pb数据
func (msg *MessageUp) GetPb() []byte {
	return msg.Pb
}

//SetPb 获取pb数据长度
func (msg *MessageUp) SetPb(pb []byte) {
	msg.Pb = pb
}

//GetKeyLen 获取key数据长度
func (msg *MessageUp) GetKeyLen() uint32 {
	return msg.KeyLen
}

//GetKey 获取key数据
func (msg *MessageUp) GetKey() []byte {
	return msg.Key
}

//SetKeyLen 获取key数据长度
func (msg *MessageUp) SetKeyLen(length uint32) {
	msg.KeyLen = length
}

//SetKey 获取key数据
func (msg *MessageUp) SetKey(key []byte) {
	msg.Key = key
}

//GetDebugLen 获取debug数据长度
func (msg *MessageUp) GetDebugLen() uint32 {
	return msg.DebugLen
}

//GetDebug 获取debug数据
func (msg *MessageUp) GetDebug() []uint64 {
	return msg.Debug
}

//SetDebugLen 获取debug数据长度
func (msg *MessageUp) SetDebugLen(length uint32) {
	msg.DebugLen = length
}

//SetDebug 获取debug数据
func (msg *MessageUp) SetDebug(debug []uint64) {
	msg.Debug = debug
}

//GetTroopLen 获取troop数据长度
func (msg *MessageUp) GetTroopLen() uint32 {
	return msg.TroopLen
}

//GetTroop 获取troop数据
func (msg *MessageUp) GetTroop() []uint32 {
	return msg.Troop
}

//SetTroopLen 获取troop数据长度
func (msg *MessageUp) SetTroopLen(length uint32) {
	msg.TroopLen = length
}

//SetTroop 获取troop数据
func (msg *MessageUp) SetTroop(troop []uint32) {
	msg.Troop = troop
}

//GetData 获取消息内容
func (msg *MessageUp) GetData() []byte {
	return msg.Data
}

//SetData 设置消息内容
func (msg *MessageUp) SetData(head, body []byte) {
	msg.Data = append(head, body...)
}
