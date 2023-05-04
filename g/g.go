package g

import (
	"log"
	"runtime"
)

var SystemType string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	SystemType = runtime.GOOS

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
