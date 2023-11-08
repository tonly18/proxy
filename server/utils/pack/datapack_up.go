package pack

import (
	"bytes"
	"encoding/binary"
	"proxy/server/utils/pack/iface"
)

// DataPackUp 封包拆包类实例，暂时不需要成员
type DataPackUp struct{}

// NewDataPackUp 封包拆包实例初始化方法
func NewDataPackUp() iface.PacketUp {
	return &DataPackUp{}
}

// 获取包头长度方法:存储二进制数据流总长度
func (dp *DataPackUp) GetHeadLen() uint32 {
	//包头长度: datalen uint32(4字节) + cmd uint32(4字节)
	return 8
}

// Pack 封包方法(压缩数据)
func (dp *DataPackUp) Pack(msg iface.IMessageUp) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写DataLen值
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//写cmd值
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetCmd()); err != nil {
		return nil, err
	}
	//写data长度
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法(解压数据)
func (dp *DataPackUp) UnPack(binaryData []byte) (iface.IMessageUp, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &MessageUp{}

	//读msgLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}

	//读cmd
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}

	//读data
	msg.Data = make([]byte, msg.MsgLen-dp.GetHeadLen())
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}

	return msg, nil
}
