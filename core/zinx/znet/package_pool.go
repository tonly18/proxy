package znet

import (
	"sync"
)

const maxSize = 8

// 包头pool对象池
var pkgHeadPool = &sync.Pool{
	New: func() any {
		return make([]byte, maxSize)
	},
}

func pkgHeadGet() []byte {
	return pkgHeadPool.Get().([]byte)
}

func pkgHeadPut(buffer []byte) {
	if len(buffer) <= maxSize && cap(buffer) <= maxSize {
		clear(buffer)
		pkgHeadPool.Put(buffer)
	}
}
