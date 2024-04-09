package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

func (s *Server) Start() {
	lintener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen失败, err:", err)
		return
	}

	fmt.Println("net.Listen成功")

	defer lintener.Close()

	for {
		conn, err := lintener.Accept()
		if err != nil {
			fmt.Println("conn失败, err:", err)
			return
		}

		fmt.Println("conn成功")

		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	fmt.Println("链接成功")
}
