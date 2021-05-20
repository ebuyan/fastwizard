package memory

import (
	"fmt"
	"runtime"
)

func Usage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB\n", bToMb(m.Alloc))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
