package main

import "net"

type User struct {
	Name string
	Addr string
	Chan chan string
	conn net.Conn

	//当前用户所属的Server
	server *Server
}

//创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		Chan: make(chan string),
		conn: conn,
		server: server,
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

// 用户上线业务
func (user *User) Online() {
	//用户上线，加入map
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	//广播用户上线
	user.server.BroadCast(user, "已上线")
}

//用户下线业务
func (user *User) Offline() {
	//用户下线，从map除去
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	//广播用户下线
	user.server.BroadCast(user, "下线")
}

//用户处理消息的业务
func (user *User) OnMessage(msg string){
	user.server.BroadCast(user, msg)
}