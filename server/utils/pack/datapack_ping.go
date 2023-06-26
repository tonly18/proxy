package pack

import (
	"bytes"
	"encoding/binary"
	"proxy/server/utils/pack/iface"
)

//DataPackPing 封包拆包类实例，暂时不需要成员
type DataPackPing struct{}

//DataPackPing 封包拆包实例初始化方法
func NewDataPackPing() iface.PacketPing {
	return &DataPackPing{}
}

//获取包头长度方法:存储二进制数据流总长度
func (dp *DataPackPing) GetHeadLen() uint32 {
	//包头长度: msgLen: uint32(4字节) + cmd: uint32(4字节)
	return 8
}

//Pack 打包方法(压缩数据)
func (dp *DataPackPing) Pack(msg iface.IMessagePing) ([]byte, error) {
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

	return dataBuff.Bytes(), nil
}

//Unpack 拆包方法(解压数据)
func (dp *DataPackPing) UnPack(binaryData []byte) (iface.IMessagePing, error) {
	//创建一个存放bytes字节的缓冲
	dataReader := bytes.NewReader(binaryData)

	//messagePing
	msg := &MessagePing{}

	//读: msglen
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}

	//读: cmd
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}

	return msg, nil
}
