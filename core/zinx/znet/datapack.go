package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"proxy/core/zinx/zconf"
	"proxy/core/zinx/ziface"
)

//MessagePack 封包拆包类实例，暂时不需要成员
type MessagePack struct{}

//MessagePack 封包拆包实例初始化方法
func NewMessagePack() ziface.Packet {
	return &MessagePack{}
}

//获取包头长度方法:存储二进制数据流总长度
func (dp *MessagePack) GetHeadLen() uint32 {
	//包头长度: msgLen: uint32(4字节) + cmd: uint32(4字节)
	return 8
}

//Pack 打包方法(压缩数据)
func (dp *MessagePack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写MsgLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//写cmd
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetCmd()); err != nil {
		return nil, err
	}
	//写data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//UnPack 拆包方法(解压数据)
func (dp *MessagePack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//message
	msg := &Message{}

	//读dataLen: 二进制流包部总长度
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}

	//读消息ID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if zconf.GlobalObject.MaxPacketSize > 0 && msg.MsgLen > zconf.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	return msg, nil
}
