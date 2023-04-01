package model

import "fmt"

type SimpleRpcResponse struct {
	Code int `json:"code"`
}

func (sresp *SimpleRpcResponse) String() string {
	return fmt.Sprintf("<Code: %d>", sresp.Code)
}

type NullRpcRequest struct{}
