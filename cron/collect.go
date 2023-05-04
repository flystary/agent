package cron

import (
	"time"

	"github.com/flystary/agent/funcs"
	"github.com/flystary/agent/g"
	"github.com/flystary/aiops/model"
)

func InitDataHistory() {
	for {
		funcs.UpdateCpuStat()
		time.Sleep(g.COLLECT_INTERVAL)
	}
}

func Collect() {
	if !g.Config().Transfer.Enabled {
		return
	}

	if len(g.Config().Transfer.Addrs) == 0 {
		return
	}

	for _, v := range funcs.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func()  []*model.MetricValue) {

}
