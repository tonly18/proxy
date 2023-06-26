// 通用数据类型定义
package utils

//ResultJSON 返回数据类型
type ResultJSON struct {
	Code uint32 `json:"code"`
	Data any    `json:"data"`
	Msg  any    `json:"msg"`
}
