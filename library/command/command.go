package command

import (
	"github.com/spf13/cast"
	"math/rand"
	"reflect"
	"time"
	"unsafe"
)

//GenTraceID 生成链路追踪ID
func GenTraceID() string {
	traceId := []byte(cast.ToString(time.Now().UnixNano()))
	rand.Shuffle(len(traceId), func(i, j int) {
		traceId[i], traceId[j] = traceId[j], traceId[i]
	})

	return B2String(traceId)
}

//IsValueNil 值判空
func IsValueNil(v interface{}) bool {
	if v == nil {
		return true
	}

	// 判断值是否为空
	type eface struct {
		v   int64
		ptr unsafe.Pointer
	}
	efacePtr := (*eface)(unsafe.Pointer(&v))
	if efacePtr == nil {
		return true
	}

	// ok := efaceptr == nil || uintptr(efaceptr.ptr) == 0
	return uintptr(efacePtr.ptr) == 0x0
}

//B2String []byte 转 string
func B2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//S2Byte string 转 []byte
func S2Byte(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
