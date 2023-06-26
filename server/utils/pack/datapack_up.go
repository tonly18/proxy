package pack

import (
	"bytes"
	"encoding/binary"
	"proxy/server/utils/pack/iface"
)

//DataPackUp 封包拆包类实例，暂时不需要成员
type DataPackUp struct{}

//NewDataPackUp 封包拆包实例初始化方法
func NewDataPackUp() iface.PacketUp {
	return &DataPackUp{}
}

//获取包头长度方法:存储二进制数据流总长度
func (dp *DataPackUp) GetHeadLen() uint32 {
	//包头长度: datalen uint32(4字节) + cmd uint32(4字节)
	return 8
}

//Pack 封包方法(压缩数据)
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
	//写pb长度
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写mode值
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMode()); err != nil {
		return nil, err
	}

	//写nin
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetUinLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetUin()); err != nil {
		return nil, err
	}
	//写pb值
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetPb()); err != nil {
		return nil, err
	}
	//写key
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetKeyLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetKey()); err != nil {
		return nil, err
	}

	//写debug
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDebugLen()); err != nil {
		return nil, err
	}
	for _, v := range msg.GetDebug() {
		if err := binary.Write(dataBuff, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}
	//写troop
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetTroopLen()); err != nil {
		return nil, err
	}
	for _, v := range msg.GetTroop() {
		if err := binary.Write(dataBuff, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}

	return dataBuff.Bytes(), nil
}

//Unpack 拆包方法(解压数据)
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

	//读数据流长度dataLen: mode:1, uin头:2, uin数据长度:msg.UinLen, pb数据长度:pbLen
	//dataLen = mode(1) + uin头(2) + UinLen + pbLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读mode
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Mode); err != nil {
		return nil, err
	}

	//读uin长度
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.UinLen); err != nil {
		return nil, err
	}
	msg.Uin = make([]byte, msg.UinLen)
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Uin); err != nil {
		return nil, err
	}

	//读pb: pbLen = dataLen - mode(1) - uin头(2) - UinLen
	msg.Pb = make([]byte, (msg.DataLen - 1 - 2 - uint32(msg.UinLen)))
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Pb); err != nil {
		return nil, err
	}

	//读Key长度
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.KeyLen); err != nil {
		return nil, err
	}
	msg.Key = make([]byte, msg.KeyLen)
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Key); err != nil {
		return nil, err
	}

	//读debug: debugLen(4字节)｜debug
	//debugLen = len(debug) + troophead(4) + len(troop)
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DebugLen); err != nil {
		return nil, err
	}
	if msg.DebugLen > 0 {
		debugLen := int(msg.DebugLen)
		if msg.DebugLen > 48 {
			debugLen = 48 //六种资源int64长度: 48 = 8 * 6
		}
		msg.Debug = make([]uint64, 0, debugLen/8)
		for x := 0; x < debugLen/8; x++ {
			var debug uint64
			if err := binary.Read(dataBuff, binary.LittleEndian, &debug); err != nil {
				return nil, err
			}
			msg.Debug = append(msg.Debug, debug)
		}
	}

	//读troop: troopLen(4字节)｜troop
	//troopLen = len(troop)
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.TroopLen); err != nil {
		return nil, err
	}
	if msg.TroopLen > 0 {
		msg.Troop = make([]uint32, 0, msg.TroopLen/4)
		for x := 0; x < int(msg.TroopLen/4); x++ {
			var troop uint32
			if err := binary.Read(dataBuff, binary.LittleEndian, &troop); err != nil {
				return nil, err
			}
			msg.Troop = append(msg.Troop, troop)
		}
	}

	return msg, nil
}
