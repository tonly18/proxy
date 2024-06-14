// IServer 服务mod接口
// IRouter 路由mod接口
// IConnection 连接mod层接口
// IMessage 消息mod接口
// IDataPack 消息拆解接口
// IMsgHandler 消息处理及协程池接口
package ziface

import (
	"context"
	"net"
	"time"
)

// 定义连接接口
type IConnection interface {
	Start()                   //启动连接，让当前连接开始工作
	Stop()                    //停止连接，结束当前连接状态M
	Context() context.Context //返回ctx，用于用户自定义的go程获取连接退出状态

	GetTCPServer() IServer          //获取当前链接对应的server
	GetTCPConnection() *net.TCPConn //从当前连接获取原始的socket TCPConn
	GetConnID() uint64              //获取当前连接ID
	GetConnMgr() IConnManager       //获取connection管理器
	GetMsgHandler() IMsgHandle      //获取消息处理器
	GetRemoteAddr() net.Addr        //获取远程客户端地址信息
	GetLocalAddr() net.Addr         //获取服务端地址信息
	GetRemoteIP() string            //获取远程地址:ip
	GetRemotePort() string          //获取远程地址:port

	SendMsg(uint32, []byte) error     //直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	SendBuffMsg(uint32, []byte) error //直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendByteMsg([]byte) error         //直接将二进制流发送给远程的TCP客户端(有缓冲)

	SetProperty(string, any) //设置链接属性
	GetProperty(string) any  //获取链接属性
	RemoveProperty(string)   //移除链接属性

	//context
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(any) any

	GetCreateTime() time.Time //conn创建时间

	IsAlive() bool                   //判断当前连接是否存活
	SetHeartBeat(IHeartbeatChecker)  //设置心跳检测器
	GetHeartBeat() IHeartbeatChecker //获取心跳检测器
}
