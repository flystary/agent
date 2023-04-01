package funcs

import (
	"strings"

	"github.com/flystary/aiops/model"
)

func NewMetricValue(metric string, val interface{}, dataType string, tages ...string) *model.MetricValue {
	mmv := model.MetricValue{
		Metric: metric,
		Value:  val,
		Type:   dataType,
	}

	size := len(tages)

	if size > 0 {
		mmv.Tags = strings.Join(tages, ",")
	}

	return &mmv
}

func GaugeValue(metric string, val interface{}, tags ...string) *model.MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val interface{}, tags ...string) *model.MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}
