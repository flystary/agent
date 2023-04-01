package model

import "fmt"


type TransferResponse struct {
	Message 	string
	Total		int
	Invalid		int
	Latency		int64
}

func (res *TransferResponse) String() string {
	return fmt.Sprintf(
		"<Total=%v, Invalid:%v, Latency=%vms, Message:%s>",
		res.Total,
		res.Invalid,
		res.Latency,
		res.Message,
	)
}