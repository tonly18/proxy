package controller

import (
	"fmt"
	"github.com/spf13/cast"
	"net/http"
	"proxy/core/server"
	"proxy/server/global"
)

func TestController(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId != "" {
		conn, err := global.GetTCPServer().GetConnMgr().GetConnByUserId(cast.ToUint64(userId))
		if err != nil { //玩家不存在
			writeResponseData(w, &server.Response{
				Code: 2000,
				Msg:  fmt.Sprintf(`player not exist. error:%v`, err),
			})
		} else {
			writeResponseData(w, &server.Response{
				Code: 2005,
				Data: fmt.Sprintf(`connID:%d, userId:%d`, conn.GetConnID(), conn.GetUserId()),
			})
		}
		return
	}

	//l, err := w.Write([]byte("Hell World!\n"))
	//l, err := io.WriteString(req.ResponseWriter, "Hell World!\n")
	//fmt.Println("l-err:", l, err)

	//conn, err := tcpserver.GetTCPServer().GetConnMgr().Get(1)
	//fmt.Println("err: ", err)
	//fmt.Println("conn: ", conn)

	//conn.SendBuffMsg(20220704, []byte(`This is a message from http server!`))
	//return &server.Response{
	//	Code: 0,
	//	Data: ,
	//}

	//if err := conn.SendByteMsg([]byte("this is a test message")); err != nil {
	//	return &server.Response{
	//		Code: 1005,
	//		Msg:  fmt.Errorf(`SendByteMsg is error:%v`, err),
	//	}
	//}

	w.Header().Set("user_id", "111")
	connNum := global.GetTCPServer().GetConnMgr().Len()
	playerNum := global.GetTCPServer().GetConnMgr().PlayerLen()
	data := fmt.Sprintf(`在线玩家数量: %v-%v`, connNum, playerNum)

	writeResponseBytes(w, []byte(data))
}
