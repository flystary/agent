package g

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var SystemType string

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    SystemType = runtime.GOOS

    // log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type LogFormatter struct{}

// 格式详情
func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05.000") //Golang MagicNumber,DO NOT TOUCH!
	var file string
	var len int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}
	//主要控制代码
	msg := fmt.Sprintf("%s [GOID:%d] %s %s:%d - %s\n", timestamp, getGID(), strings.ToUpper(entry.Level.String()), file, len, entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine"))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

