package ziface

/*
将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	GetMsgLen() uint32 //获取消息包总长度
	SetMsgLen(uint32)  //设置消息包总长度

	GetCmd() uint32 //获取cmd数据
	SetCmd(uint32)  //设置cmd数据

	SetData([]byte)  //设置data数据
	GetData() []byte //获取data数据
}
