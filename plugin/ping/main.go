package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var timeout = 25

type MetricValue struct {
	Endpoint  string      `json:"endpoint"`
	Metric    string      `json:"metric"`
	Value     interface{} `json:"value"`
	Step      int64       `json:"step"`
	Type      string      `json:"counterType"`
	Tags      string      `json:"tags"`
	Timestamp int64       `json:"timestamp"`
}

var err error

func main() {
	var cfg []string
	if len(os.Args) > 1 {
		arg := os.Args
		s := arg[1]
		err := json.Unmarshal([]byte(s), &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		BatchPing(cfg)
	} else {
		return
	}
}
func BatchPing(addrs []string) {
	result := make([]*MetricValue, 0)
	defer func() {
		if err != nil {
			result = append(result, GaugeValue("check.process.ping", -1, ""))
		}
	}()
	chs := make([]chan []*MetricValue, len(addrs))
	for i, addr := range addrs {
		chs[i] = make(chan []*MetricValue)
		go run(addr, chs[i])
	}
	for _, ch := range chs {
		temp := <-ch
		if temp != nil {
			result = append(result, temp...)
		}
	}
	b, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	fmt.Fprintf(os.Stdout, string(b))
	return
}
func run(addr string, ch chan []*MetricValue) {
	ch_run := make(chan []*MetricValue)
	go doCommand(addr, ch_run)
	select {
	case re := <-ch_run:
		ch <- re
	case <-time.After(time.Duration(timeout) * time.Second):
		ch <- nil
	}
}
func doCommand(addr string, ch chan []*MetricValue) []*MetricValue {
	cmd := exec.Command("bash", "-c", "ping "+addr+" -c 20 -i 0.5")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		ch <- nil
		return nil
	}
	err, isTimeout := CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Second)
	if isTimeout {
		fmt.Fprintf(os.Stderr, "exec cmd : ping %s timeout", addr)
		ch <- nil
		return nil
	}
	errStr := stderr.String()
	if errStr != "" {
		fmt.Fprintf(os.Stderr, errStr)
		ch <- nil
		return nil
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ping %s failed: %s", addr, err.Error())
		ch <- nil
		return nil
	}
	arr := strings.Fields(stdout.String())
	var loss string
	var rtt string
	for i := 0; i < len(arr); i++ {
		if arr[i] == "received," {
			loss = arr[i+1][0 : len(arr[i+1])-1]
		}
		if arr[i] == "rtt" {
			rtt = strings.Split(arr[i+3], "/")[1]
		}
	}
	_, err = strconv.ParseFloat(loss, 32)
	_, err = strconv.ParseFloat(rtt, 32)
	result := make([]*MetricValue, 2)
	result[0] = GaugeValue("ping.result.loss", loss, addr)
	result[1] = GaugeValue("ping.result.rtt", rtt, addr)
	ch <- result
	return nil
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
func CmdRunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	var err error
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(timeout):
		go func() {
			<-done // allow goroutine to exit
		}()
		//IMPORTANT: cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} is necessary before cmd.Start()
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		return err, true
	case err = <-done:
		return err, false
	}
}
