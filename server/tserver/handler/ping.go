package handler

import (
	"fmt"
	"proxy/core/zinx/ziface"
	"proxy/library/logger"
	"proxy/server/global"
	"proxy/server/utils/pack"
	"time"
)

// PingRouter Struct
type PingRouter struct {
	BaseHandler
}

func (h *PingRouter) Handle(request ziface.IRequest) error {
	//开始时间
	start := time.Now()

	//params
	data := request.GetData()
	fmt.Println("ping-data::::::", string(data))

	downMsg := pack.NewMessageDown(global.CMD_DOWN_PONG, 0, nil)
	dp := pack.NewDataPackDown()
	data, err := dp.Pack(downMsg)
	if err != nil {
		return fmt.Errorf(`[ping handler] unpack error:%v`, err)
	}
	if err := request.GetConnection().SendByteMsg(data); err != nil {
		return fmt.Errorf(`[ping handler] conn send error:%v`, err)
	}

	//结束时间(毫秒)
	end := time.Since(start).Milliseconds()
	logger.Infof(request, `[Status Code:%d | MsgID:%d | Execution Time:%d(ms)]`, 200, request.GetMsgID(), end)

	//return
	return nil
}
