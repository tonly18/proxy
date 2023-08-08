package znet

import (
	"errors"
	"fmt"
	"proxy/core/zinx/zconf"
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/zlog"
	"sync"
)

//ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint64]ziface.IConnection //map[connID]conn
	players     map[uint64]uint64             //map[userId]connID
	connLock    sync.RWMutex
}

//NewConnManager 创建一个链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint64]ziface.IConnection, zconf.GlobalObject.MaxConn),
		players:     make(map[uint64]uint64, zconf.GlobalObject.MaxConn),
	}
}

//Add 添加链接
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

//Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	//删除players
	if conn.GetUserId() > 0 {
		if conn.GetConnID() == connMgr.players[conn.GetUserId()] {
			delete(connMgr.players, conn.GetUserId())
		}
	}

	connMgr.connLock.Unlock()

	zlog.Infof(`[Conn Manager Remove] Conn Remove ConnID:%v Successfully! Conn Number:%v, Address:%v`, conn.GetConnID(), connMgr.Len(), conn.GetRemoteAddr())
}

//Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}
	for k, v := range connMgr.players {
		if v == connID {
			delete(connMgr.players, k)
			break
		}
	}

	return nil, errors.New("connection not found")
}

//Len 获取当前连接
func (connMgr *ConnManager) Len() int {
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}

//ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
		//删除players
		delete(connMgr.players, conn.GetUserId())
	}
	for k, _ := range connMgr.players {
		delete(connMgr.players, k)
	}

	connMgr.connLock.Unlock()

	zlog.Info("[Conn Manager ClearConn] Clear All Connections successfully: conn num = ", connMgr.Len())
}

//GetConnByUserId 根据userId获取链接
func (connMgr *ConnManager) GetConnByUserId(userId uint64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if connID, ok := connMgr.players[userId]; ok {
		if conn, ok := connMgr.connections[connID]; ok {
			return conn, nil
		} else {
			delete(connMgr.players, userId)
		}
	}

	return nil, errors.New("connection not found")
}

//AddConnByUserId 添加到players
func (connMgr *ConnManager) AddConnByUserId(conn ziface.IConnection) error {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	//未登录
	if conn.GetUserId() < 1 {
		conn.Stop()
		return fmt.Errorf(`conn manager player not login, donot add to players. user id:%v`, conn.GetUserId())
	}
	//已经添加到players
	if _, ok := connMgr.players[conn.GetUserId()]; ok {
		conn.Stop()
		return fmt.Errorf(`conn manager player already exists. user id:%v`, conn.GetUserId())
	}
	//已经登录时,则添加到players
	connMgr.players[conn.GetUserId()] = conn.GetConnID()
	//是否已添加到connections
	if _, ok := connMgr.connections[conn.GetConnID()]; !ok {
		connMgr.connections[conn.GetConnID()] = conn
	}

	return nil
}

//PlayerLen 获取当前连接
func (connMgr *ConnManager) PlayerLen() int {
	connMgr.connLock.RLock()
	//在线人数
	length := len(connMgr.players)

	//修正数据不一致
	//if length > len(connMgr.connections) {
	//	for uid, connId := range connMgr.players {
	//		if _, ok := connMgr.connections[connId]; !ok {
	//			delete(connMgr.players, uid)
	//		}
	//	}
	//	length = len(connMgr.players)
	//}

	connMgr.connLock.RUnlock()

	return length
}
