package iface

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessageKickOut interface {
	GetMsgLen() uint32 //获取消息包总长度
	SetMsgLen(uint32)  //设置消息包总长度

	GetCmd() uint32 //获取cmd数据
	SetCmd(uint32)  //设置cmd数据

	GetErrorCode() uint32 //获取Mode数据
	SetErrorCode(uint32)  //设置Mode数据
}
