package main

import "net"

type User struct {
	Name string
	Addr string
	Chan chan string
	conn net.Conn
}

//创建一个用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		Chan: make(chan string),
		conn: conn,
	}

	//启动监听当前user channel的goroutine
	go user.ListenMessage()

	return user
}

//监听当前 User channel 的方法，一旦有消息就发送给客户端
func (user *User) ListenMessage() {
	for true {
		msg := <- user.Chan
		user.conn.Write([]byte(msg + "\n"))
	}
}
