package server

//Response 响应
type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
	Msg  any `json:"msg"`

	Type int8 `json:"-"`
}
