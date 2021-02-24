package main

import (
	"net"
	"strings"
)

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
		Name:   userAddr,
		Addr:   userAddr,
		Chan:   make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听当前user channel的goroutine
	go user.ListenMessage()

	return user
}

//监听当前 User channel 的方法，一旦有消息就发送给客户端
func (user *User) ListenMessage() {
	for true {
		msg := <-user.Chan
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

//发送消息到当前用户所在客户端
func (user *User) SendMessage(msg string) {
	user.conn.Write([]byte(msg))
}

//处理用户消息的业务，不是这个用户收到的消息，而是他发出的消息
func (user *User) OnMessage(msg string) {
	if msg == "who" {
		user.server.mapLock.Lock()
		for _, onlineUser := range user.server.OnlineMap {
			onlineMsg := "[" + onlineUser.Addr + "]" + onlineUser.Name + ":" + "在线...\n"
			user.SendMessage(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename=" {
		//消息格式：rename=张三
		newName := strings.Split(msg, "=")[1]

		//判断当前名称是否被占用
		if _, ok := user.server.OnlineMap[newName]; ok {
			user.SendMessage("当前用户名称被使用\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.Name = newName
			user.server.OnlineMap[user.Name] = user
			user.server.mapLock.Unlock()

			user.SendMessage("更新用户名成功：" + user.Name + "\n")
		}

	} else if len(msg) > 3 && msg[:3] == "to=" {
		//消息格式：to|张三|消息内容

		//1. 获取用户名和消息内容
		remoteName := strings.Split(msg, "=")[1]
		if remoteName == "" {
			user.SendMessage("消息格式不正确，请使用 \"to|张三|消息内容\" 格式\n")
			return
		}
		content := strings.Split(msg, "=")[2]
		if content == "" {
			user.SendMessage("消息内容不能为空")
			return
		}

		//2. 通过用户名找到要发送到的用户，并发送
		remoteUser, ok := user.server.OnlineMap[remoteName]
		if !ok {
			user.SendMessage("没有找到用户：" + remoteName)
		} else {
			remoteUser.SendMessage(user.Name + " says: " + content)
		}
	} else {
		user.server.BroadCast(user, msg)
	}
}
