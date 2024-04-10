package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	lock      sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
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

	// 启动监听
	go s.ListenMessage()

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

// 当前客户端链接的业务
func (s *Server) Handler(conn net.Conn) {

	user := NewUser(conn)

	// 用户上线, 将用户添加到OnlineMap中
	s.lock.Lock()
	s.OnlineMap[user.Name] = user
	s.lock.Unlock()

	// 广播消息给所有用户
	s.Fanout(user, "已上线")
}

func (s *Server) Fanout(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]-%s: %s", user.Name, user.Addr, msg)
	s.Message <- sendMsg
}

// 监听服务端收到消息
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		s.lock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.lock.Unlock()
	}
}
