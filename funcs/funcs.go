package funcs

import (
	"github.com/flystary/agent/g"
	"github.com/flystary/aiops/model"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval

	Mappers = []FuncsAndInterval{
		{
			Fs: []func() []*model.MetricValue{
				// CpuMetrics,
			},
			Interval: interval,
		},
	}
}
