package znet

import (
	"context"
	"errors"
	"proxy/core/zinx/ziface"
	"proxy/library/command"
	"proxy/utils"
	"sync"
	"time"
)

const (
	PRE_HANDLE  ziface.HandleStep = iota // PreHandle for pre-processing
	HANDLE                               // Handle for processing
	POST_HANDLE                          // PostHandle for post-processing

	HANDLE_OVER
)

// Request 请求
type Request struct {
	conn     ziface.IConnection //已经和客户端建立好的 链接
	msg      ziface.IMessage    //客户端请求的数据
	router   ziface.IRouter     //请求处理的函数
	steps    ziface.HandleStep  //用来控制路由函数执行
	stepLock *sync.RWMutex      //并发互斥
	traceId  string             //链路追踪ID
	args     map[string]any
}

func NewRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	return &Request{
		conn:     conn,
		msg:      msg,
		steps:    PRE_HANDLE,
		stepLock: new(sync.RWMutex),
		traceId:  command.GenTraceID(), //初始化链路追踪ID
		args:     make(map[string]any, 5),
	}
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetCmd()
}

// BindRouter 绑定路由
func (r *Request) BindRouter(router ziface.IRouter) {
	r.router = router
}

func (r *Request) next() {
	r.stepLock.Lock()
	r.steps++
	r.stepLock.Unlock()
}

func (r *Request) Call() error {
	if r.router == nil {
		return errors.New("router is empty")
	}

	for r.steps < HANDLE_OVER {
		var err error
		switch r.steps {
		case PRE_HANDLE:
			err = r.router.PreHandle(r)
		case HANDLE:
			err = r.router.Handle(r)
		case POST_HANDLE:
			err = r.router.PostHandle(r)
		}
		if err != nil {
			r.steps = PRE_HANDLE
			return err
		}

		r.next()
	}
	r.steps = PRE_HANDLE

	return nil
}

//Abort 终止处理函数的运行,但调用此方法的函数会执行完毕
//func (r *Request) Abort() {
//	r.stepLock.Lock()
//	r.steps = HANDLE_OVER
//	r.stepLock.Unlock()
//}

// SetAargs 数据
func (r *Request) SetAargs(key string, value any) {
	r.args[key] = value
}
func (r *Request) GetAargs(key string) any {
	return r.args[key]
}

// GetTraceId 链路追踪ID
func (r *Request) GetTraceId() string {
	return r.traceId
}

// GetCtx
func (r *Request) GetCtx() context.Context {
	return r.conn.Context()
}

// Deadline
func (r *Request) Deadline() (deadline time.Time, ok bool) {
	return r.conn.Context().Deadline()
}

// Done
func (r *Request) Done() <-chan struct{} {
	return r.conn.Context().Done()
}

// Err
func (r *Request) Err() error {
	return r.conn.Context().Err()
}

// Value
func (r *Request) Value(key any) any {
	if keywords, ok := key.(string); ok {
		if keywords == utils.ClientIP {
			return r.conn.GetRemoteIP()
		}
		if keywords == utils.TraceID {
			return r.GetTraceId()
		}
		return r.conn.GetProperty(keywords)
	}

	return nil
}
