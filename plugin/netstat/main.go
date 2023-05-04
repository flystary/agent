package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/toolkits/nux"
	"github.com/toolkits/slice"
)

var timeout = 30

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}
type PortConfig struct {
	TcpPort []int64 `json:"tcp"`
	UdpPort []int64 `json:"udp"`
	All     []int64 `json:"all"`
}

func main() {
	var cfg PortConfig
	if len(os.Args) > 1 {
		arg := os.Args
		s := arg[1]
		err := json.Unmarshal([]byte(s), &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}
	if len(cfg.TcpPort) == 0 && len(cfg.UdpPort) == 0 && len(cfg.All) == 0 {
		return
	}
	result := PortMetrics(&cfg)
	b, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	fmt.Fprintf(os.Stdout, string(b))
}
func PortMetrics(config *PortConfig) (L []*MetricValue) {
	allTcpPorts, err := nux.TcpPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	allUdpPorts, err := nux.UdpPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	for _, tcp := range config.TcpPort {
		tags := fmt.Sprintf("port=%d", tcp)
		if slice.ContainsInt64(allTcpPorts, tcp) {
			L = append(L, GaugeValue("port.status.plugin.tcp", 1, tags))
		} else {
			L = append(L, GaugeValue("port.status.plugin.tcp", 0, tags))
		}
	}
	for _, udp := range config.UdpPort {
		tags := fmt.Sprintf("port=%d", udp)
		if slice.ContainsInt64(allUdpPorts, udp) {
			L = append(L, GaugeValue("port.status.plugin.udp", 1, tags))
		} else {
			L = append(L, GaugeValue("port.status.plugin.udp", 0, tags))
		}
	}
	for _, all := range config.All {
		tags := fmt.Sprintf("port=%d", all)
		if slice.ContainsInt64(allTcpPorts, all) && slice.ContainsInt64(allUdpPorts, all) {
			L = append(L, GaugeValue("port.status.plugin.all", 1, tags))
		} else {
			L = append(L, GaugeValue("port.status.plugin.all", 0, tags))
		}
	}
	return
}
func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *MetricValue {
	mv := MetricValue{
		Metric: metric,
		Value:  val,
		Type:   dataType,
	}
	size := len(tags)
	if size > 0 {
		mv.Tags = strings.Join(tags, ",")
	}
	return &mv
}
func GaugeValue(metric string, val interface{}, tags ...string) *MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}
