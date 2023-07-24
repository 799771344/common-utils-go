package server

import (
	jrpc "e.coding.net/new-journey/journey/jrpc_go/jrpc_go"
	"net"
)

type Args struct {
	A, B int
}

type ArithImpl struct{}

func (t *ArithImpl) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func (t *ArithImpl) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	s, _ := jrpc.NewServer(":1234")
	s.Start(new(ArithImpl))

	co, _ := jrpc.NewClient("localhost:1234")
	resp := []interface{}{}
	_ = co.Call("Arith.Add", net.Interface{}, &resp)

}
