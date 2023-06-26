package pack

import (
	"bytes"
	"encoding/binary"
	"proxy/server/utils/pack/iface"
)

//DataPackKickOut 封包拆包类实例，暂时不需要成员
type DataPackKickOut struct{}

//NewDataPackKickOut 封包拆包实例初始化方法
func NewDataPackKickOut() iface.PacketKickOut {
	return &DataPackKickOut{}
}

//获取包头长度方法:存储二进制数据流总长度
func (dp *DataPackKickOut) GetHeadLen() uint32 {
	//包头长度: msgLen: uint32(4字节) + cmd: uint32(4字节) + errorcode:uint32(4字节)
	return 12
}

//Pack 打包方法(压缩数据)
func (dp *DataPackKickOut) Pack(msg iface.IMessageKickOut) ([]byte, error) {
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
	//写ErrorCode
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetErrorCode()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//UnPack 拆包方法(解压数据)
func (dp *DataPackKickOut) UnPack(binaryData []byte) (iface.IMessageKickOut, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//messageDown
	msg := &MessageDown{}

	//读dataLen: 二进制流包部总长度
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}
	//读cmd
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}
	//读ErrorCode
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ErrorCode); err != nil {
		return nil, err
	}

	return msg, nil
}
