package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

//创建一个server对象
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip: ip,
		Port: port,
	}
}

func (server *Server) Handler(conn net.Conn)  {
	//当前连接的业务
	fmt.Printf("连接建立成功\n")
}

//启动服务器的方法 启动一个socket监听在ip:port
func (server *Server) Start(){
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%server:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Printf("net.Listen err : %v\n", err)
		return
	}

	//close listener socket
	defer listener.Close()

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