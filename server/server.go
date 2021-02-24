package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
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
		msg := <-server.MessageChan
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
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

	user := NewUser(conn, server)

	user.Online()

	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for true {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Printf("conn read err : %v\n", err)
				return
			}

			//提取用户的消息，去除\n
			msg := string(buf[:n-1])

			//用户针对msg进行处理
			user.OnMessage(msg)

			//用户的任意消息代表当前用户活跃
			isLive <- true
		}

	}()

	//当前handler阻塞，如果该方法结束，goroutine结束
	for true {
		select {
		case <-isLive:
			//当前用户活跃，应该重置触发器
			//不做操作，顺序实行下条
		case <-time.After(time.Minute * 15):
			//已经超时，将当前用户强制关闭
			user.SendMessage("你被踢了")

			//销毁用户的资源
			close(user.Chan)
			//关闭连接
			conn.Close()
			//退出当前的Handler
			runtime.Goexit() //return
		}
	}
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
