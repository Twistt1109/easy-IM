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
}

func NewUser(conn net.Conn) *User {

	addr := conn.RemoteAddr().String()
	user := &User{
		Name: addr,
		Addr: addr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前user 的 channel
	go user.ListenMessage()

	return user
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
