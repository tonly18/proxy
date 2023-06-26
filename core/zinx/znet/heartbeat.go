package znet

import (
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/zlog"
	"time"
)

type HeartbeatChecker struct {
	interval time.Duration      //心跳检测时间间隔
	quitChan chan bool          //退出信号
	conn     ziface.IConnection //绑定的链接
}

//HeatBeatDefaultRouter 收到remote心跳消息的默认回调路由业务
type HeatBeatDefaultRouter struct{}

func NewHeartbeatChecker(interval time.Duration) ziface.IHeartbeatChecker {
	heartbeat := &HeartbeatChecker{
		interval: interval,
	}

	return heartbeat
}

func (h *HeartbeatChecker) start() {
	zlog.Infof("[HeartBeat Start] Goroutine is Running! ConnID=%+v, Address:%v", h.conn.GetConnID(), h.conn.GetRemoteAddr())
	ticker := time.NewTicker(h.interval)
	for {
		select {
		case <-ticker.C:
			h.check()
		case <-h.quitChan:
			ticker.Stop()
			zlog.Infof("[HeartBeat Start] Goroutine Exit! ConnID=%+v, Address:%v", h.conn.GetConnID(), h.conn.GetRemoteAddr())
			return
		}
	}
}

func (h *HeartbeatChecker) Start() {
	go h.start()
}

func (h *HeartbeatChecker) Stop() {
	zlog.Infof("[HeartBeat Stop] Stop! ConnID=%+v, Address:%v", h.conn.GetConnID(), h.conn.GetRemoteAddr())
	h.quitChan <- true
}

func (h *HeartbeatChecker) check() error {
	if h.conn == nil {
		return nil
	}
	if !h.conn.IsAlive() {
		zlog.Infof(`[HeartBeat check] remote connection %s is not alive, stop it`, h.conn.GetRemoteAddr())
		h.conn.Stop()
	}
	return nil
}

func (h *HeartbeatChecker) BindConn(conn ziface.IConnection) {
	h.conn = conn
	conn.SetHeartBeat(h)
}

//Clone 克隆到一个指定的链接上
func (h *HeartbeatChecker) Clone() ziface.IHeartbeatChecker {
	heartbeat := &HeartbeatChecker{
		interval: h.interval,
		quitChan: make(chan bool, 1),
		conn:     nil,
	}

	return heartbeat
}
