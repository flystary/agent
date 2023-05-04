package g

import (
	// "log"
	"runtime"
    //log "github.com/sirupsen/logrus"
)

var SystemType string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	SystemType = runtime.GOOS

	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
