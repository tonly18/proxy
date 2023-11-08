package pack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"proxy/server/utils/pack/iface"
)

// DataPackDown 封包拆包类实例，暂时不需要成员
type DataPackDown struct{}

// DataPackDown 封包拆包实例初始化方法
func NewDataPackDown() iface.PacketDown {
	return &DataPackDown{}
}

// 获取包头长度方法
func (dp *DataPackDown) GetHeadLen() uint32 {
	//包头长度: msgLen: uint32(4字节) + cmd: uint32(4字节) + code: uint32(4字节)
	return 12
}

// Pack 打包方法(压缩数据)
func (dp *DataPackDown) Pack(msg iface.IMessageDown) ([]byte, error) {
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
	//写Code
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetCode()); err != nil {
		return nil, err
	}
	//写data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// UnPack 拆包方法(解压数据)
func (dp *DataPackDown) UnPack(binaryData []byte) (iface.IMessageDown, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//messageDown
	msg := &MessageDown{}

	//读dataLen: 二进制流包部总长度(包头+包体)
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}
	//读cmd
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}
	//读Code
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Code); err != nil {
		return nil, err
	}

	//读data
	if dataBuff.Len() > 0 {
		msg.Data = make([]byte, msg.MsgLen-dp.GetHeadLen())
		fmt.Println("msg.Data:::", len(msg.Data))
		if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
			return nil, err
		}
	}

	return msg, nil
}
