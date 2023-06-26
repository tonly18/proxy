package znet

import (
	"context"
	"errors"
	"fmt"
	"net"
	"proxy/core/zinx/zconf"
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/zlog"
	"sync/atomic"
	"time"
)

//Server 接口实现，定义一个Server服务类
type Server struct {
	//服务器ID
	ID uint32
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)

	ctx   context.Context
	canel context.CancelFunc

	packet ziface.Packet

	//connection id
	connId uint64
	//心跳检测器
	hc ziface.IHeartbeatChecker
}

//NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	s := &Server{
		ID:         zconf.GlobalObject.ServerID,
		Name:       zconf.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         zconf.GlobalObject.Host,
		Port:       zconf.GlobalObject.TCPPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		packet:     NewMessagePack(),
		ctx:        nil,
		canel:      nil,
	}
	s.ctx, s.canel = context.WithCancel(context.Background())

	return s
}

func (s *Server) StartConn(conn ziface.IConnection) {
	//HeartBeat check
	if s.hc != nil {
		s.hc.Clone().BindConn(conn)
	}

	//Start Conn
	conn.Start()
}

//Start 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[TCP SERVER START] Server Name: %s, Listenner at IP: %s, Port %d is Starting\n", s.Name, s.IP, s.Port)

	//0 启动worker工作池机制
	s.msgHandler.StartWorkerPool()

	//1 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		panic(fmt.Errorf(`[TCP SERVER START] Resolve Tcp Addr Error: %w`, err))
	}

	//2 监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		panic(fmt.Sprintf(`[TCP SERVER START] net Listen Tcp Error: %v`, err))
	}
	fmt.Println("[TCP SERVER START] Server:", s.Name, "Successful, Now Listening...")

	//开启一个go去做服务端Linster业务
	go func() {
		//3 启动server网络连接业务
		for {
			//3.1 设置服务器最大连接控制,如果超过最大连接,那么则关闭此新的连接
			if s.ConnMgr.Len() >= zconf.GlobalObject.MaxConn {
				zlog.Infof(`[TCP SERVER START] Exceeded the maxConnNum:%d, Wait:%d`, zconf.GlobalObject.MaxConn, AcceptDelay.duration)
				AcceptDelay.Delay()
				continue
			}

			//3.2 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				//Go 1.16+
				if errors.Is(err, net.ErrClosed) {
					zlog.Error("[TCP SERVER START] Accept Err(1):", err)
					return
				}
				zlog.Error("[TCP SERVER START] Accept Err(2):", err)
				AcceptDelay.Delay()
				continue
			}
			AcceptDelay.Reset()
			zlog.Info("[TCP SERVER START] A New Conn Remote Addr:", conn.RemoteAddr().String())

			//3.3 处理该新连接请求的业务方法,此时应该有handler和conn是绑定的
			cid := atomic.AddUint64(&s.connId, 1)
			dealConn := NewConnection(s, conn, cid, s.msgHandler)

			//3.4 启动当前链接的处理业务
			go s.StartConn(dealConn)
		}
	}()

	select {
	case <-s.ctx.Done():
		zlog.Info("[TCP SERVER START] listener close")
		if err := listener.Close(); err != nil {
			zlog.Errorf(`[TCP SERVER START] listener close err: %v`, err)
		}
	}
}

//Stop 停止服务
func (s *Server) Stop() {
	zlog.Info("[TCP SERVER STOP] Server Name:", s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()

	//退出
	s.canel()
}

//Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//阻塞,否则主Go退出,listener的go将会退出
	select {
	case <-s.ctx.Done():
		zlog.Info("[TCP SERVER SERVE] Context Cancel")
	}
}

//AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

//GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		//fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

//CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() ziface.Packet {
	return s.packet
}

func (s *Server) GetID() uint32 {
	return s.ID
}

//StartHeartBeat 启动心跳检测interval 每次发送心跳的时间间隔
func (s *Server) StartHeartBeat(interval time.Duration) {
	checker := NewHeartbeatChecker(interval)

	//server绑定心跳检测器
	s.hc = checker
}

func (s *Server) GetHeartBeat() ziface.IHeartbeatChecker {
	return s.hc
}
