package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
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

func main() {
	var cfg string
	if len(os.Args) > 1 {
		arg := os.Args
		cfg = arg[1]
	} else {
		return
	}
	doCommand(cfg)
}
func doCommand(filename string) {
	var err error
	var result []*MetricValue
	defer func() {
		if err != nil {
			result = append(result, GaugeValue("check.process.filename", -1, "file="+filename))
		}
	}()
	//tips:osx  does not support -b.
	//cmd := exec.Command("sh", "-c", "ps -ef | grep " + filename)
	//cmd := exec.Command("sh", "-c", "ps -ef | grep " + filename + " | grep -v grep")//todo:NOTWORK
	cmd := exec.Command("sh", "-c", "ps -ef | grep "+filename+" | grep -v grep | awk '{print $2}'")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	err, isTimeout := CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Second)
	if isTimeout {
		fmt.Fprintf(os.Stderr, "exec cmd : ps -ef %s timeout", filename)
		return
	}
	errStr := stderr.String()
	if errStr != "" {
		fmt.Fprintf(os.Stderr, errStr)
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ps -ef %s failed: %s", filename, err.Error())
		return
	}
	arr := strings.Fields(stdout.String())
	//fmt.Println(arr)
	if len(arr) > 0 {
		result = append(result, GaugeValue("check.process.filename", 1, "path="+filename))
	} else {
		result = append(result, GaugeValue("check.process.filename", 0, "path="+filename))
	}
	b, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	fmt.Fprintf(os.Stdout, string(b))
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
