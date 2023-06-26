package pack

import (
	"proxy/server/utils/pack/iface"
)

//MessageLogin 消息
type MessageLogin struct {
	MsgLen uint32 //消息总长度
	Cmd    uint32 //消息ID长度

	UinLen uint32 //uin长度
	Uin    []byte //uin内容:string

	ServerLen uint32 //server长度
	Server    []byte //server内容:string

	SessionKeyLen uint32 //SessionKey长度
	SessionKey    []byte //SessionKey内容:string

	VersionLen uint32 //Version长度
	Version    []byte //Version内容:string

	Data []byte //消息全部的内容
}

//NewMsgLoginPackage 创建一个Message消息包
func NewMessageLogin(cmd uint32, uin []byte, server []byte, sessionKey []byte, version []byte) iface.IMessageLogin {
	msg := &MessageLogin{
		MsgLen: 0,
		Cmd:    cmd,

		UinLen: uint32(len(uin)),
		Uin:    uin,

		ServerLen: uint32(len(server)),
		Server:    server,

		SessionKeyLen: uint32(len(sessionKey)),
		SessionKey:    sessionKey,

		VersionLen: uint32(len(version)),
		Version:    version,

		Data: nil,
	}

	//msg总长度
	msg.MsgLen = 4 + 4 + (4 + msg.UinLen) + (4 + msg.ServerLen) + (4 + msg.SessionKeyLen) + (4 + msg.VersionLen)

	return msg
}

//GetMsgLen 获取二进制流总长度
func (msg *MessageLogin) GetMsgLen() uint32 {
	return msg.MsgLen
}
func (msg *MessageLogin) SetMsgLen(length uint32) {
	msg.MsgLen = length
}

//GetCmd 获取cmd数据
func (msg *MessageLogin) GetCmd() uint32 {
	return msg.Cmd
}
func (msg *MessageLogin) SetCmd(cmd uint32) {
	msg.Cmd = cmd
}

//GetUin 获取Uin数据
func (msg *MessageLogin) GetUinLen() uint32 {
	return msg.UinLen
}
func (msg *MessageLogin) GetUin() []byte {
	return msg.Uin
}
func (msg *MessageLogin) SetUinLen(length uint32) {
	msg.UinLen = length
}
func (msg *MessageLogin) SetUin(userID []byte) {
	msg.Uin = userID
}

//GetServer 获取server数据
func (msg *MessageLogin) GetServerLen() uint32 {
	return msg.ServerLen
}
func (msg *MessageLogin) GetServer() []byte {
	return msg.Server
}
func (msg *MessageLogin) SetServerLen(length uint32) {
	msg.ServerLen = length
}
func (msg *MessageLogin) SetServer(server []byte) {
	msg.Server = server
}

//GetSessionKey 获取sessionKey数据
func (msg *MessageLogin) GetSessionKeyLen() uint32 {
	return msg.SessionKeyLen
}
func (msg *MessageLogin) GetSessionKey() []byte {
	return msg.SessionKey
}
func (msg *MessageLogin) SetSessionKeyLen(length uint32) {
	msg.SessionKeyLen = length
}
func (msg *MessageLogin) SetSessionKey(sessionKey []byte) {
	msg.SessionKey = sessionKey
}

//GetVersion 获取version数据
func (msg *MessageLogin) GetVersionLen() uint32 {
	return msg.VersionLen
}
func (msg *MessageLogin) GetVersion() []byte {
	return msg.Version
}
func (msg *MessageLogin) SetVersionLen(length uint32) {
	msg.VersionLen = length
}
func (msg *MessageLogin) SetVersion(version []byte) {
	msg.Version = version
}

//GetData 获取消息内容
func (msg *MessageLogin) GetData() []byte {
	return msg.Data
}

//SetData 设计消息内容
func (msg *MessageLogin) SetData(head, data []byte) {
	msg.Data = append(head, data...)
}
