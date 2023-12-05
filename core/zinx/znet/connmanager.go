package znet

import (
	"errors"
	"fmt"
	"proxy/core/zinx/zconf"
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/zlog"
	"sync"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint64]ziface.IConnection //map[connID]conn
	players     map[uint64]uint64             //map[userId]connID
	connLock    sync.RWMutex
}

// NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]ziface.IConnection, zconf.GlobalObject.MaxConn),
		players:     make(map[uint64]uint64, zconf.GlobalObject.MaxConn),
	}
}

// Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	//将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	//如果conn(已登录),则添加到players中
	if conn.GetUserId() > 0 {
		connMgr.players[conn.GetUserId()] = conn.GetConnID()
	}
	connMgr.connLock.Unlock()

	zlog.Infof(`[Conn Manager Add] Connection Add To ConnManager Successfully! Conn Number:%v, Address:%v`, connMgr.Len(), conn.GetRemoteAddr())
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()

	//删除players
	if cid, ok := connMgr.players[conn.GetUserId()]; ok {
		if conn.GetConnID() == cid {
			delete(connMgr.players, conn.GetUserId())
		}
	}
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	connMgr.connLock.Unlock()

	zlog.Infof(`[Conn Manager Remove] Conn Remove ConnID:%v Successfully! Conn Number:%v, Address:%v`, conn.GetConnID(), connMgr.Len(), conn.GetRemoteAddr())
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

// Len 获取当前连接
func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

// ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for uid, cid := range connMgr.players {
		delete(connMgr.players, uid)
		if conn, ok := connMgr.connections[cid]; ok {
			conn.Stop()                      //停止
			delete(connMgr.connections, cid) //删除
		}
	}
	//停止并删除全部的连接信息
	for cid, conn := range connMgr.connections {
		conn.Stop()                      //停止
		delete(connMgr.connections, cid) //删除
	}

	connMgr.connLock.Unlock()

	zlog.Info("[Conn Manager ClearConn] Clear All Connections successfully: conn num = ", connMgr.Len())
}

// GetConnByUserId 根据userId获取链接
func (connMgr *ConnManager) GetConnByUserId(userId uint64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if connID, ok := connMgr.players[userId]; ok {
		if conn, exist := connMgr.connections[connID]; exist {
			return conn, nil
		}
	}

	return nil, errors.New("connection not found")
}

// AddConnByUserId 添加到players
func (connMgr *ConnManager) AddConnByUserId(conn ziface.IConnection) error {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//未登录
	if conn.GetUserId() < 1 {
		return fmt.Errorf(`conn manager player not login, donot add to players. user id:%v`, conn.GetUserId())
	}
	//已经添加到players
	if _, ok := connMgr.players[conn.GetUserId()]; ok {
		return fmt.Errorf(`conn manager player already exists. user id:%v`, conn.GetUserId())
	}
	//已经登录时,则添加到players
	connMgr.players[conn.GetUserId()] = conn.GetConnID()

	return nil
}

// PlayerLen 获取当前连接
func (connMgr *ConnManager) GetOnLinePlayer() int {
	connMgr.connLock.Lock()
	//在线人数
	length := len(connMgr.players)

	//修正数据不一致
	if length != len(connMgr.connections) {
		for uid, connId := range connMgr.players {
			if _, ok := connMgr.connections[connId]; !ok {
				delete(connMgr.players, uid)
			}
		}
		for connId, conn := range connMgr.connections {
			if conn.GetUserId() > 0 {
				if _, ok := connMgr.players[connId]; !ok {
					conn.Stop()
					delete(connMgr.players, conn.GetUserId())
				}
			}
		}
		length = len(connMgr.players)
	}

	connMgr.connLock.Unlock()

	return length
}
