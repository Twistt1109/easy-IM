package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {

	addr := conn.RemoteAddr().String()
	user := &User{
		Name:   addr,
		Addr:   addr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前user 的 channel
	go user.ListenMessage()

	return user
}

func (u *User) Online() {
	u.server.lock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.lock.Unlock()

	// 广播消息给所有用户
	u.server.Fanout(u, "已上线")
}

func (u *User) Offline() {
	u.server.lock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.lock.Unlock()

	// 广播消息给所有用户
	u.server.Fanout(u, "已下线")
}

func (u *User) DoMessage(msg string) {
	u.server.Fanout(u, msg)
}

// 服务端 -> 客户端
// 监听当前user的 channel, 有消息就取值发送给当前user的客户端(conn)
func (u *User) ListenMessage() {

	for {
		msg := <-u.C
		count, err := u.conn.Write([]byte(msg + "\n"))
		fmt.Println("count:", count, "err:", err)
	}
}
