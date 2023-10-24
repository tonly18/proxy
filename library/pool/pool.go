package pool

import (
	"sync"
)

const maxSize = 8

var poolBuffer = &sync.Pool{
	New: func() any {
		return make([]byte, maxSize)
	},
}

func PoolGet() []byte {
	return poolBuffer.Get().([]byte)
}

func PoolPut(buffer []byte) {
	clear(buffer)
	if len(buffer) == maxSize && cap(buffer) == maxSize {
		poolBuffer.Put(buffer)
	}
}
