package server

import (
	"github.com/spf13/cast"
	"net/http"
	"proxy/core/zinx/ziface"
	"proxy/library/command"
	"time"
)

// Request 请求
type Request struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	UserID         uint64 //当前玩家ID
	Conn           ziface.IConnection
	data           map[string]any
}

func (r *Request) GetData(key string) any {
	return r.data[key]
}

func (r *Request) SetData(key string, value any) {
	if len(r.data) == 0 {
		r.data = make(map[string]any, 10)
	}
	r.data[key] = value
}

// Deadline
func (r *Request) Deadline() (deadline time.Time, ok bool) {
	return r.Request.Context().Deadline()
}

// Done
func (r *Request) Done() <-chan struct{} {
	return r.Request.Context().Done()
}

// Err
func (r *Request) Err() error {
	return r.Request.Context().Err()
}

// Value
func (r *Request) Value(key any) any {
	value := r.GetData(cast.ToString(key))
	if command.IsValueNil(value) {
		value = r.Request.Context().Value(key)
		if command.IsValueNil(value) {
			if r.Conn != nil {
				value = r.Conn.GetProperty(cast.ToString(key))
			}
		}
	}

	return value
}
