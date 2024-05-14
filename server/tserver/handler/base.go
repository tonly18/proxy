package handler

import (
	"fmt"
	"proxy/core/zinx/ziface"
	"proxy/core/zinx/znet"
	"proxy/library/logger"
	"proxy/server/utils/pack"
	"time"
)

// BaseHandler Struct
type BaseHandler struct {
	znet.BaseRouter
	startTime time.Time
}

func (bh *BaseHandler) PreHandle(request ziface.IRequest) error {
	bh.startTime = time.Now()
	return nil
}

func (bh *BaseHandler) Handle(request ziface.IRequest) error {
	return nil
}

func (bh *BaseHandler) PostHandle(request ziface.IRequest) error {
	endTime := time.Since(bh.startTime).Milliseconds() //结束时间(毫秒)

	//执行时长
	logger.Infof(request, `[MsgID:%d | Execution Time:%d(ms)]`, request.GetMsgID(), endTime)

	return nil
}

// 向客户端推送消息
func PushMessage(req ziface.IRequest, cmd, code uint32, data []byte) error {
	downMsg := pack.NewMessageDown(cmd, code, data)
	dp := pack.NewDataPackDown()
	msgPack, err := dp.Pack(downMsg)
	if err != nil {
		return fmt.Errorf(`base push message unpack error:%v`, err)
	}
	if err := req.GetConnection().SendByteMsg(msgPack); err != nil {
		return fmt.Errorf(`base push message conn send error:%v`, err)
	}

	//return
	return nil
}
