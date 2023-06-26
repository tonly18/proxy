package iface

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessageUp interface {
	GetMsgLen() uint32 //获取消息包总长度
	SetMsgLen(uint32)  //设置消息包总长度

	GetCmd() uint32 //获取cmd数据
	SetCmd(uint32)  //设置cmd数据

	GetMode() uint8 //获取Mode数据
	SetMode(uint8)  //设置Mode数据

	GetDataLen() uint32 //获取数据长度
	SetDataLen(uint32)  //设置pb数据长度

	GetUinLen() uint16 //获取uin数据长度
	GetUin() []byte    //获取uin数据
	SetUinLen(uint16)  //设置uin数据长度
	SetUin([]byte)     //设置uin数据

	GetPb() []byte //获取pb数据
	SetPb([]byte)  //设置pb数据

	GetKeyLen() uint32 //获取key数据长度
	GetKey() []byte    //获取key数据
	SetKeyLen(uint32)  //设置key数据长度
	SetKey([]byte)     //设置key数据

	GetDebugLen() uint32 //获取debug数据长度
	GetDebug() []uint64  //获取debug数据
	SetDebugLen(uint32)  //设置debug数据长度
	SetDebug([]uint64)   //设置debug数据

	GetTroopLen() uint32 //获取troop数据长度
	GetTroop() []uint32  //获取troop数据
	SetTroopLen(uint32)  //设置troop数据长度
	SetTroop([]uint32)   //设置troop数据

	GetData() []byte        //获取消息内容
	SetData([]byte, []byte) //设置消息内容
}
