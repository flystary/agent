package funcs

import "github.com/flystary/aiops/model"

func AgentMetrics() (mms []*model.MetricValue) {
	return []*model.MetricValue{GaugeValue("agent.alive", 1)}
}
