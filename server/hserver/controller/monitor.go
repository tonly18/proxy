package controller

import (
	"fmt"
	"net/http"
	"runtime"
)

//groution 监控
func MonitorGroutionController(w http.ResponseWriter, r *http.Request) {
	data := fmt.Sprintf(`goroutine数量: %v`, runtime.NumGoroutine())

	writeResponseBytes(w, []byte(data))
}

//memory 监控
func MonitorMemoryController(w http.ResponseWriter, r *http.Request) {
	//内存数据
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	data := fmt.Sprintf("TotalAlloc: %v(M), Alloc: %v(M), Mallocs: %v(次), Frees: %v(次), GCNumber: %v(次)", fmt.Sprintf("%.2f", float64(mem.TotalAlloc)/1024/1024), fmt.Sprintf("%.2f", float64(mem.Alloc)/1024/1024), mem.Mallocs, mem.Frees, mem.NumGC)

	writeResponseBytes(w, []byte(data))
}
