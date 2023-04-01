package model

import "fmt"

type AegntReportRequest struct {
	Hostname string
	IP		 string
	AgentVersion	string
	PluginVersion	string
}

func (res *AegntReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname:%s, IP:%s, AgentVersion:%s, PluginVersion:%s>",
		res.Hostname,
		res.IP,
		res.AgentVersion,
		res.PluginVersion,
	)
}

type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AegntReportRequest
}