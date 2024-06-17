package znet

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"io"
	"net"
	"proxy/core/zinx/zconf"
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/zlog"
	"proxy/library/pool"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Connection 链接
type Connection struct {
	//当前Conn属于哪个Server
	tcpServer ziface.IServer
	//当前连接的socket TCP套接字
	conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	connID uint64
	//消息管理MsgID和对应处理方法的消息管理模块
	msgHandler ziface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	//链接属性
	property map[string]any
	//保护当前property的锁
	propertyLock sync.Mutex

	//当前链接是属于哪个Connection Manager
	connManager ziface.IConnManager
	//当前连接创建时Hook函数
	onConnStart func(ziface.IConnection)
	//当前连接断开时的Hook函数
	onConnStop func(ziface.IConnection)
	//当前链接的远程地址
	remoteAddr net.Addr
	//当前链接的本地地址
	localAddr net.Addr
	//conn创建时间
	createTime time.Time
	//最后一次活动时间
	lastActivityTime time.Time
	//心跳检测器
	hc ziface.IHeartbeatChecker
}

// NewConnection 创建连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint64) ziface.IConnection {
	//初始化Conn属性
	c := &Connection{
		tcpServer:   server,
		conn:        conn,
		connID:      connID,
		connManager: server.GetConnMgr(),
		msgHandler:  server.GetMsgHandler(),
		msgBuffChan: make(chan []byte, zconf.GlobalObject.MaxMsgChanLen),
		property:    make(map[string]any),
		onConnStart: server.GetOnConnStart(),
		onConnStop:  server.GetOnConnStop(),
		remoteAddr:  conn.RemoteAddr(),
		localAddr:   conn.LocalAddr(),
		createTime:  time.Now(),
	}

	//将新创建的Conn添加到链接管理中
	c.connManager.Add(c)

	//return
	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	zlog.Info("[Conn Write] Goroutine is Running!", c.GetRemoteAddr())

	defer func() {
		zlog.Info("[Conn Write] Goroutine is Exit!", c.GetRemoteAddr())
		if err := recover(); err != nil {
			zlog.Error("[Conn Write] Goroutine is Exit Error!", c.GetRemoteAddr(), err)
		}
		c.Stop()
	}()

	for {
		select {
		case <-c.ctx.Done():
			zlog.Info("[Conn Write] context cancel!", c.GetRemoteAddr())
			return
		case data, ok := <-c.msgBuffChan:
			if ok {
				if err := c.Send(data); err != nil {
					zlog.Errorf("[Conn Writer] Send error:%v", err)
				}
			} else {
				zlog.Error("[Conn Write] msgBuffChan is Closed!")
			}
			break
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	zlog.Info("[Conn Read] Goroutine is Running!", c.GetRemoteAddr())

	defer func() {
		zlog.Info("[Conn Read] Goroutine is Exit!", c.GetRemoteAddr())
		if err := recover(); err != nil {
			zlog.Error("[Conn Read] Goroutine is Exit Error!", c.GetRemoteAddr(), err)
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
				c.GetTCPConnection().SetReadDeadline(time.Now().Add(time.Second * time.Duration(zconf.GlobalObject.MaxConnReadTime)))
			}

			//读取message head
			msgHeadBuffer := pool.PoolGet()
			if _, err := io.ReadFull(c.GetTCPConnection(), msgHeadBuffer); err != nil {
				pool.PoolPut(msgHeadBuffer)
				zlog.Errorf(`[Conn Read] Read Msg Head Error:%v, Address:%v`, err, c.GetRemoteAddr())
				return
			}
			//拆包:得到datalen、cmd并放在msg中
			msg, err := c.GetTCPServer().Packet().UnPack(msgHeadBuffer)
			pool.PoolPut(msgHeadBuffer)
			if err != nil {
				zlog.Errorf(`[Conn Read] Unpack Error:%v, Address:%v`, err, c.GetRemoteAddr())
				return
			}
			//根据dataLen读取data,放在msg.Data中
			if msg.GetMsgLen() > c.GetTCPServer().Packet().GetHeadLen() {
				msgBodyBuffer := make([]byte, msg.GetMsgLen()-c.GetTCPServer().Packet().GetHeadLen())
				if _, err := io.ReadFull(c.GetTCPConnection(), msgBodyBuffer); err != nil {
					zlog.Error("[Conn Read] Read Msg Data Error:", err)
					return
				}
				msg.SetData(msgBodyBuffer) //设置message body
			}

			//正常读取到对端数据,更新心跳检测Active状态
			if c.hc != nil {
				c.UpdateActivity()
			}

			//Request 得到当前客户端请求的Request数据
			req := NewRequest(c, msg)

			//执行request
			if zconf.GlobalObject.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.msgHandler.SendMsgToTaskQueue(req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.msgHandler.DoMsgHandler(req)
			}
		}
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	zlog.Infof(`[Conn Start] Goroutine is Running! Addr:%v`, c.GetRemoteAddr())

	defer func() {
		zlog.Infof(`[Conn Start] Goroutine is Exit! Addr:%v`, c.GetRemoteAddr())
		if err := recover(); err != nil {
			zlog.Errorf(`[Conn Start] Goroutine is Exit! Addr:%v, Error:%v`, c.GetRemoteAddr(), err)
			for i := 1; i < 20; i++ {
				if pc, file, line, ok := runtime.Caller(i); ok {
					function := runtime.FuncForPC(pc).Name() //获取函数名
					zlog.Errorf(`[Conn Start] goroutine:%v, file:%s, func:%s, line:%d`, pc, file, function, line)
				}
			}
		}
	}()

	//context
	c.ctx, c.cancel = context.WithCancel(context.Background())

	//开启心跳检测器
	if c.hc != nil {
		c.UpdateActivity()
		c.hc.Start()
	}

	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	//按照用户传递进来的创建连接时需要处理的业务,执行钩子方法
	c.callOnConnStart()

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

// Stop 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	c.cancel()
}

// GetTCPServer 获取TCPServer
func (c *Connection) GetTCPServer() ziface.IServer {
	return c.tcpServer
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint64 {
	return c.connID
}

func (c *Connection) GetConnMgr() ziface.IConnManager {
	return c.connManager
}

// GetMsgHandler 获取消息处理器
func (c *Connection) GetMsgHandler() ziface.IMsgHandle {
	return c.msgHandler
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.remoteAddr
}

// GetLocalAddr 获取服务端地址信息
func (c *Connection) GetLocalAddr() net.Addr {
	return c.localAddr
}

// GetRemoteIP ip
func (c *Connection) GetRemoteIP() string {
	return strings.Split(c.remoteAddr.String(), ":")[0]
}

// GetRemotePort port
func (c *Connection) GetRemotePort() string {
	return strings.Split(c.remoteAddr.String(), ":")[1]
}

// Send 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) Send(data []byte) error {
	if c.isClosed() == true {
		return errors.New("[Conn Send] connection closed when send msg")
	}

	//写回客户端: 设置写入数据流超时时间
	startTime := time.Now()
	if zconf.GlobalObject.MaxConnWriteTime > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(time.Duration(zconf.GlobalObject.MaxConnWriteTime) * time.Millisecond))
	}
	if n, err := c.conn.Write(data); err != nil {
		zlog.Errorf("[Conn Send] writed length:%v, raw data length:%v, time:%v, error:%v", n, len(data), time.Since(startTime).Milliseconds(), err)
		return err
	}

	return nil
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed() == true {
		return errors.New("connection closed when send msg")
	}

	//将data封包，并且发送
	dp := c.GetTCPServer().Packet()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		zlog.Error("[Conn SendMsg] Pack error msg ID:", msgID, ", err:", err)
		return errors.New("pack error msg")
	}

	if err := c.Send(msg); err != nil {
		zlog.Errorf("[Conn SendMsg] SendMsg err msg ID:%d, data:%+v, error:%+v", msgID, string(msg), err)
		return err
	}

	return nil
}

// SendByteMsg 发生BuffMsg
func (c *Connection) SendByteMsg(msgID uint32, data []byte) error {
	if c.isClosed() == true {
		return errors.New("connection closed when send buff msg")
	}

	//将data封包，并且发送
	dp := c.tcpServer.Packet()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		zlog.Error("[Conn SendByteMsg] Pack error msg ID = ", msgID, " Err: ", err)
		return errors.New("pack error msg")
	}

	if err := c.SendBuffMsg(msg); err != nil {
		zlog.Errorf("[conn SendByteMsg] error:%v", err)
		return err
	}

	return nil
}

// SendBuffMsg 发生BuffMsg
func (c *Connection) SendBuffMsg(data []byte) error {
	if c.isClosed() == true {
		return errors.New("connection closed when send byte msg")
	}

	//time out
	idleTimeout := time.NewTimer(10 * time.Millisecond)
	defer idleTimeout.Stop()

	//发送超时
	select {
	case <-c.ctx.Done():
		return errors.New("[conn SendBuffMsg] connection closed when send buff msg")
	case <-idleTimeout.C:
		zlog.Error("[conn SendBuffMsg] send buff msg timeout")
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- data:
		return nil
	}

	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value any) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) any {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	return c.property[key]
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// Context 返回ctx,用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) finalizer() {
	//将链接从连接管理器中删除
	if c.connManager != nil {
		c.connManager.Remove(c)
	}
	//如果用户注册了该链接的关闭回调业务,那么在此刻应该显示调用
	c.callOnConnStop()

	//停止心跳检测器
	if c.hc != nil {
		c.hc.Stop()
	}

	//关闭socket链接
	_ = c.conn.Close()

	//关闭该链接全部管道
	if c.msgBuffChan != nil {
		close(c.msgBuffChan)
	}

	//logger
	zlog.Infof(`[Conn Finalizer] Conn Stop ConnID:%v, Address:%v`, c.GetConnID(), c.GetRemoteAddr())
}

// Deadline
func (c *Connection) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done
func (c *Connection) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err
func (c *Connection) Err() error {
	return c.ctx.Err()
}

// Value
func (c *Connection) Value(key any) any {
	if k, ok := key.(string); ok {
		if k == "client_ip" {
			return c.GetRemoteIP()
		}
	}

	return c.GetProperty(cast.ToString(key))
}

// GetCreateTime 链接创建时间
func (c *Connection) GetCreateTime() time.Time {
	return c.createTime
}

func (c *Connection) callOnConnStart() {
	if c.onConnStart != nil {
		zlog.Info(fmt.Sprintf("callOnConnStart, remote Addr:%v, conn id:%v", c.GetRemoteAddr(), c.GetConnID()))
		c.onConnStart(c)
	}
}

func (c *Connection) callOnConnStop() {
	if c.onConnStop != nil {
		zlog.Info(fmt.Sprintf("callOnConnStop, remote Addr:%v, conn id:%v", c.GetRemoteAddr(), c.GetConnID()))
		c.onConnStop(c)
	}
}

func (c *Connection) IsAlive() bool {
	if c.isClosed() {
		return false
	}

	//检查连接最后一次活动时间,如果超过心跳间隔,则认为连接已经死亡
	return time.Now().Sub(c.lastActivityTime) < zconf.GlobalObject.HeartbeatMaxDuration()
}

// GetActivity 获取conn最后活跃时间
func (c *Connection) GetActivity() time.Time {
	return c.lastActivityTime
}

// UpdateActivity 更新conn最后活跃时间
func (c *Connection) UpdateActivity() {
	c.lastActivityTime = time.Now()
}

// SetHeartBeat 设置心跳检测器
func (c *Connection) SetHeartBeat(checker ziface.IHeartbeatChecker) {
	c.hc = checker
}

// GetHeartBeat 获取心跳检测器
func (c *Connection) GetHeartBeat() ziface.IHeartbeatChecker {
	return c.hc
}

// isClosed 链接是否已关闭
func (c *Connection) isClosed() bool {
	return c.ctx.Err() != nil
}
