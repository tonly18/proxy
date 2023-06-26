package pack

import (
	"bytes"
	"encoding/binary"
	"proxy/server/utils/pack/iface"
)

//DataPackLogin 封包拆包类实例，暂时不需要成员
type DataPackLogin struct{}

//DataPackLogin 封包拆包实例初始化方法
func NewDataPackLogin() iface.PacketLogin {
	return &DataPackLogin{}
}

//获取包头长度方法:存储二进制数据流总长度
func (dp *DataPackLogin) GetHeadLen() uint32 {
	//包头长度: msgLen: uint32(4字节) + cmd: uint32(4字节)
	return 8
}

//Pack 打包方法(压缩数据)
func (dp *DataPackLogin) Pack(msg iface.IMessageLogin) ([]byte, error) {
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

	//写Uin
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetUinLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetUin()); err != nil {
		return nil, err
	}

	//写server
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetServerLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetServer()); err != nil {
		return nil, err
	}

	//写sessionKey
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetSessionKeyLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetSessionKey()); err != nil {
		return nil, err
	}

	//写version
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetVersionLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetVersion()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//Unpack 拆包方法(解压数据)
func (dp *DataPackLogin) UnPack(binaryData []byte) (iface.IMessageLogin, error) {
	//创建一个存放bytes字节的缓冲
	dataReader := bytes.NewReader(binaryData)

	//messageLogin
	msg := &MessageLogin{}

	//读: msglen
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}
	//读: cmd
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.Cmd); err != nil {
		return nil, err
	}

	//读userID
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.UinLen); err != nil {
		return nil, err
	}
	msg.Uin = make([]byte, msg.GetUinLen())
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.Uin); err != nil {
		return nil, err
	}

	//读server
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.ServerLen); err != nil {
		return nil, err
	}
	msg.Server = make([]byte, msg.GetServerLen())
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.Server); err != nil {
		return nil, err
	}

	//读sessionKey
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.SessionKeyLen); err != nil {
		return nil, err
	}
	msg.SessionKey = make([]byte, msg.GetSessionKeyLen())
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.SessionKey); err != nil {
		return nil, err
	}

	//读Version
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.VersionLen); err != nil {
		return nil, err
	}
	msg.Version = make([]byte, msg.GetVersionLen())
	if err := binary.Read(dataReader, binary.LittleEndian, &msg.Version); err != nil {
		return nil, err
	}

	return msg, nil
}
