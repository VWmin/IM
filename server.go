package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	MessageChan chan string
}

//创建一个server对象
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:          ip,
		Port:        port,
		OnlineMap:   make(map[string]*User),
		MessageChan: make(chan string),
	}
}

//监听Message Channel广播消息的goroutine，一旦有消息，就发送给全部的在线用户
func (server *Server) ListenMessageChan() {
	for true {
		msg := <- server.MessageChan
		server.mapLock.Lock()
		for _, user := range server.OnlineMap{
			user.Chan <- msg
		}
		server.mapLock.Unlock()
	}
}

//广播消息的方法
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.MessageChan <- sendMsg
}

func (server *Server) Handler(conn net.Conn) {
	//当前连接的业务
	//fmt.Printf("连接建立成功\n")

	user := NewUser(conn)

	//用户上线，加入map
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()

	//广播用户上线
	server.BroadCast(user, "已上线")

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for true {
			n, err := conn.Read(buf)
			if n == 0 {
				server.BroadCast(user, "下线")
				return
			}

			if err != nil && err != io.EOF{
				fmt.Printf("conn read err : %v\n", err)
				return
			}

			//提取用户的消息，去除\n
			msg := string(buf[:n-1])
			//将得到的消息广播
			server.BroadCast(user, msg)
		}

	}()

	//当前handler阻塞，如果该方法结束，goroutine结束
	//select {}
}

//启动服务器的方法 启动一个socket监听在ip:port
func (server *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Printf("net.Listen err : %v\n", err)
		return
	}

	//close listener socket
	defer listener.Close()

	//启动监听message channel的goroutine
	go server.ListenMessageChan()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("listener accept err : %v\n", err)
			continue
		}
		//do handler
		go server.Handler(conn)
	}

}