package g

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

func InitLog(level string) (err error) {

	var logSize int
	var logPath string

	if config.Log.Size <= 0 {
		logSize = 10
	} else {
		logSize = config.Log.Size
	}

	if config.Log.Path == "" {
		logPath = "/var/log/"
	} else {
		logPath = config.Log.Path
	}
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.Fatal("log conf only allow [info, debug, warn], please check your confguire")
	}

	/* 日志切割参数
	   `Filename` 日志文件名
	   `MaxSize` 单个日志文件的大小，单位M
	   `MaxBackups` 设置备份日志的文件数
	   `Compress` 是否压缩日志
	*/
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath + "/agent-" + time.Now().Format("2006-01-02") + ".log",
		MaxSize:    logSize,
		MaxBackups: 1,
		Compress:   false,
	}
	mw := io.MultiWriter(os.Stdout, lumberJackLogger)
	log.SetOutput(mw)
	log.SetFormatter(new(LogFormatter))
	return
}
