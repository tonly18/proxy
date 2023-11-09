package command

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"unsafe"
)

// GenTraceID 生成链路追踪ID
func GenTraceID() string {
	traceId := GenRandom()
	return strconv.Itoa(int(traceId))
}

// IsValueNil 值判空
func IsValueNil(v any) bool {
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

// B2String []byte 转 string
func B2String(b []byte) string {
	return unsafe.String(&b[0], len(b))
}

// S2Byte string 转 []byte
func S2Byte(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// SliceJoin 链接为字符串
func SliceJoin[T comparable](s []T, sep string) string {
	var result bytes.Buffer
	for _, v := range s {
		result.WriteString(fmt.Sprintf(`%v`, v))
	}

	return result.String()
}

// GenRandom 生成随机数
func GenRandom() int64 {
	randomNum, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	return randomNum.Int64()
}
