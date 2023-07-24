package main

import (
	"log"
	"net"
	"net/rpc"
)

type AAA struct {
}

type Server struct {
	listener net.Listener
}

func NewServer(addr string) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: lis,
	}, nil
}

func (s *Server) Start(serverStruct interface{}) {
	rpc.Register(serverStruct)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("accept error: %v", err)
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	s, _ := NewServer(":1234")
	s.Start(new(AAA))
}
