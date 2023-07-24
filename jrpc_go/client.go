package main

import (
	"net"
	"net/rpc"
)

type Conn struct {
	conn *rpc.Client
}

// NewClient 创建一个客户端
func NewClient(addr string) (*Conn, error) {
	// 建立连接
	conn, err := rpc.Dial("tpc", addr)
	if err != nil {
		return nil, err
	}

	// 返回连接
	return &Conn{conn: conn}, nil
}

// Call 发起 RPC 请求并返回响应结果
func (c *Conn) Call(serviceMethod string, req interface{}, resp *[]interface{}) error {
	// 发起 gRPC 请求
	err := c.conn.Call(serviceMethod, req, &resp)
	if err != nil {
		return err
	}

	// 返回响应结果
	return nil
}

func main() {
	co, _ := NewClient("localhost:1234")
	resp := []interface{}{}
	_ = co.Call("Arith.Add", net.Interface{}, &resp)

}
