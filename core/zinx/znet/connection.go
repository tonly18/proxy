package znet

import (
	"context"
	"errors"
	"github.com/spf13/cast"
	"io"
	"net"
	"proxy/core/zinx/zconf"
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/zlog"
	"proxy/library/command"
	"strings"
	"sync"
	"time"
)

//Connection 链接
type Connection struct {
	//当前Conn属于哪个Server
	TCPServer ziface.IServer
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint64
	//消息管理MsgID和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	//用户发消息的Lock
	msgLock sync.RWMutex
	//链接属性
	property map[string]any
	//玩家ID
	userId uint64
	//保护当前property的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool
	//当前链接的远程地址
	remoteAddr net.Addr
	//当前链接的本地地址
	localAddr net.Addr
	//conn创建时间
	createTime int32
	//最后一次活动时间
	lastActivityTime time.Time
	//心跳检测器
	hc ziface.IHeartbeatChecker
}

//NewConnection 创建连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint64, msgHandler ziface.IMsgHandle) *Connection {
	//初始化Conn属性
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgBuffChan: make(chan []byte, zconf.GlobalObject.MaxMsgChanLen),
		property:    make(map[string]any),
		remoteAddr:  conn.RemoteAddr(),
		localAddr:   conn.LocalAddr(),
		createTime:  int32(time.Now().Unix()),
	}

	//将新创建的Conn添加到链接管理中
	c.TCPServer.GetConnMgr().Add(c)

	//return
	return c
}

//StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	zlog.Info("[Conn Write] Goroutine is Running!", c.GetRemoteAddr().String())

	defer func() {
		zlog.Info("[Conn Write] Goroutine is Exit!", c.GetRemoteAddr().String())
		if err := recover(); err != nil {
			zlog.Error("[Conn Write] Goroutine is Exit Error!", c.GetRemoteAddr().String())
		}
		c.Stop()
	}()

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//设置写入数据流时间(100毫秒)
				if zconf.GlobalObject.MaxConnWriteTime > 0 {
					c.Conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(zconf.GlobalObject.MaxConnWriteTime)))
				}

				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					zlog.Error("[Conn Write] Send Buff Data Error:", err, ", Conn Writer exit")
					break
				}
			} else {
				zlog.Error("[Conn Write] MsgBuffChan is Closed!]")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

//StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	zlog.Info("[Conn Read] Goroutine is Running!", c.GetRemoteAddr().String())

	defer func() {
		zlog.Info("[Conn Read] Goroutine is Exit!", c.GetRemoteAddr().String())
		if err := recover(); err != nil {
			zlog.Error("[Conn Read] Goroutine is Exit Error!", c.GetRemoteAddr().String())
		}
		c.Stop()
	}()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//设置读取数据流时间
			if zconf.GlobalObject.MaxConnReadTime > 0 {
				c.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(zconf.GlobalObject.MaxConnReadTime)))
			}

			//读取客户端的Msg head
			msgHeadBuffer := make([]byte, c.TCPServer.Packet().GetHeadLen())
			if _, err := io.ReadFull(c.Conn, msgHeadBuffer); err != nil {
				zlog.Errorf(`[Conn Read] Read Msg Head Error:%v, Address:%v`, err, c.GetRemoteAddr())
				return
			}

			//拆包:得到msgID和datalen放在msg中
			msg, err := c.TCPServer.Packet().UnPack(msgHeadBuffer)
			if err != nil {
				zlog.Errorf(`[Conn Read] Unpack Error:%v, Address:%v`, err, c.GetRemoteAddr())
				return
			}

			//根据dataLen读取data,放在msg.Data中
			if msg.GetMsgLen() < c.TCPServer.Packet().GetHeadLen() {
				continue
			}
			var msgBodyBuffer []byte
			if msg.GetMsgLen() > 0 {
				msgBodyBuffer = make([]byte, msg.GetMsgLen()-c.TCPServer.Packet().GetHeadLen())
				if _, err := io.ReadFull(c.Conn, msgBodyBuffer); err != nil {
					zlog.Error("[Conn Read] Read Msg Data Error:", err)
					return
				}
			}
			msg.SetRawData(msgHeadBuffer, msgBodyBuffer) //设置原始二进制流和data

			//正常读取到对端数据,更新心跳检测Active状态
			if c.hc != nil {
				c.updateActivity()
			}

			//Request 得到当前客户端请求的Request数据
			req := NewRequest(c, msg)
			//设置链路追踪ID
			req.SetTraceId(command.GenTraceID())

			if zconf.GlobalObject.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.MsgHandler.SendMsgToTaskQueue(req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.MsgHandler.DoMsgHandler(req)
			}
		}
	}
}

//Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	zlog.Infof(`[Conn Start] Goroutine is Running! Addr:%v`, c.GetRemoteAddr())

	defer func() {
		zlog.Infof(`[Conn Start] Goroutine is Exit! Addr:%v`, c.GetRemoteAddr())
		if err := recover(); err != nil {
			zlog.Infof(`[Conn Start] Goroutine is Exit! Addr:%v, Error:%v`, c.GetRemoteAddr(), err)
		}
	}()

	//context
	c.ctx, c.cancel = context.WithCancel(context.Background())

	//开启心跳检测器
	if c.hc != nil {
		c.updateActivity()
		c.hc.Start()
	}

	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TCPServer.CallOnConnStart(c)

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

//Stop 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	c.cancel()
}

//GetTCPServer 获取TCPServer
func (c *Connection) GetTCPServer() ziface.IServer {
	return c.TCPServer
}

//GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint64 {
	return c.ConnID
}

//RemoteAddr 获取远程客户端地址信息
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.remoteAddr
}

//GetLocalAddr 获取服务端地址信息
func (c *Connection) GetLocalAddr() net.Addr {
	return c.localAddr
}

//GetRemoteIP ip
func (c *Connection) GetRemoteIP() string {
	return strings.Split(c.remoteAddr.String(), ":")[0]
}

//GetRemotePort port
func (c *Connection) GetRemotePort() string {
	return strings.Split(c.remoteAddr.String(), ":")[1]
}

//SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	//将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		zlog.Error("[Conn SendMsg] Pack error msg ID:", msgID, ", err:", err)
		return errors.New("pack error msg")
	}

	//写回客户端: 设置写入数据流时间(100毫秒)
	if zconf.GlobalObject.MaxConnWriteTime > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(time.Duration(zconf.GlobalObject.MaxConnWriteTime) * time.Millisecond))
	}
	_, err = c.Conn.Write(msg)

	return err
}

//SendBuffMsg  发生BuffMsg
func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	idleTimeout := time.NewTimer(10 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	//将data封包，并且发送
	dp := c.TCPServer.Packet()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		zlog.Error("[Conn SendBuffMsg] Pack error msg ID = ", msgID, " Err: ", err)
		return errors.New("pack error msg")
	}

	// 发送超时
	select {
	case <-idleTimeout.C:
		zlog.Error("[conn SendBuffMsg] Send Buff Msg Timeout")
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
	//写回客户端
	//c.msgBuffChan <- msg

	return nil
}

//SendByteMsg  发生BuffMsg
func (c *Connection) SendByteMsg(data []byte) error {
	//lock
	c.msgLock.RLock()
	defer c.msgLock.RUnlock()

	//time out
	idleTimeout := time.NewTimer(10 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	//发送超时
	select {
	case <-idleTimeout.C:
		zlog.Error("[conn SendByteMsg] Send Buff Msg Timeout")
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- data:
		return nil
	}

	return nil
}

//SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value any) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) any {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	return c.property[key]
}

//RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

//Context 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}

//SetProxyId 网关
func (c *Connection) SetProxyId(proxyId uint32) {
	c.SetProperty("proxy_id", proxyId)
}
func (c *Connection) GetProxyId() uint32 {
	return cast.ToUint32(c.GetProperty("proxy_id"))
}

//SetServerId 区服ID
func (c *Connection) SetServerId(serverId uint32) {
	c.SetProperty("server_id", serverId)
}
func (c *Connection) GetServerId() uint32 {
	return cast.ToUint32(c.GetProperty("server_id"))
}

//SetUIN 帐号ID
func (c *Connection) SetUIN(uin uint64) {
	c.SetProperty("uin", uin)
}
func (c *Connection) GetUIN() uint64 {
	return cast.ToUint64(c.GetProperty("uin"))
}

//SetUserId 玩家ID
func (c *Connection) SetUserId(userId uint64) {
	c.userId = userId
}
func (c *Connection) GetUserId() uint64 {
	return c.userId
}

func (c *Connection) finalizer() {
	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	//如果用户注册了该链接的关闭回调业务,那么在此刻应该显示调用
	c.TCPServer.CallOnConnStop(c)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	//停止心跳检测器
	if c.hc != nil {
		c.hc.Stop()
	}

	//关闭socket链接
	c.Conn.Close()
	//将链接从连接管理器中删除
	c.TCPServer.GetConnMgr().Remove(c)
	//关闭该链接全部管道
	if c.msgBuffChan != nil {
		close(c.msgBuffChan)
	}
	//设置标志位
	c.isClosed = true

	zlog.Infof(`[Conn Finalizer] Conn Stop ConnID:%v, UserID:%v, Address:%v`, c.ConnID, c.GetUserId(), c.GetRemoteAddr())
}

//Deadline
func (c *Connection) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

//Done
func (c *Connection) Done() <-chan struct{} {
	return c.ctx.Done()
}

//Err
func (c *Connection) Err() error {
	return c.ctx.Err()
}

//Value
func (c *Connection) Value(key any) any {
	if k, ok := key.(string); ok {
		if k == "client_ip" {
			return c.GetRemoteIP()
		}
	}

	return c.GetProperty(cast.ToString(key))
}

//GetCreateTime 链接创建时间
func (c *Connection) GetCreateTime() int32 {
	return c.createTime
}

func (c *Connection) IsAlive() bool {
	if c.isClosed {
		return false
	}
	//检查连接最后一次活动时间,如果超过心跳间隔,则认为连接已经死亡
	return time.Now().Sub(c.lastActivityTime) < zconf.GlobalObject.HeartbeatMaxDuration()
}

func (c *Connection) updateActivity() {
	c.lastActivityTime = time.Now()
}

//SetHeartBeat 设置心跳检测器
func (c *Connection) SetHeartBeat(checker ziface.IHeartbeatChecker) {
	c.hc = checker
}

//GetHeartBeat 获取心跳检测器
func (c *Connection) GetHeartBeat() ziface.IHeartbeatChecker {
	return c.hc
}
