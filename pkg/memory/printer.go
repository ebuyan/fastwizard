package memory

import (
	"fmt"
	"runtime"
	"unsafe"
)

func Usage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB\n", bToMb(m.Alloc))
}

func Size(vars interface{}) {
	fmt.Printf("var: %d\n", unsafe.Sizeof(vars))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
