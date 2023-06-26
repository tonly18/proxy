package iface

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessageLogin interface {
	GetMsgLen() uint32 //获取消息包总长度
	SetMsgLen(uint32)  //设置消息包总长度

	GetCmd() uint32 //获取cmd数据
	SetCmd(uint32)  //设置cmd数据

	GetUinLen() uint32 //获取UinLen数据
	GetUin() []byte    //获取Uin数据
	SetUinLen(uint32)  //设置UinLen数据
	SetUin([]byte)     //设置Uin数据

	GetServerLen() uint32 //获取server数据
	GetServer() []byte    //获取server数据
	SetServerLen(uint32)  //设置serverLen数据
	SetServer([]byte)     //设置server数据

	GetSessionKeyLen() uint32 //获取sessionKeyLen数据
	GetSessionKey() []byte    //获取sessionKey数据
	SetSessionKeyLen(uint32)  //设置sessionKeyLen数据
	SetSessionKey([]byte)     //设置sessionKey数据

	GetVersionLen() uint32 //获取versionLen数据
	GetVersion() []byte    //获取version数据
	SetVersionLen(uint32)  //设置versionLen数据
	SetVersion([]byte)     //设置sversion数据

	GetData() []byte        //获取消息内容
	SetData([]byte, []byte) //设置消息内容
}
