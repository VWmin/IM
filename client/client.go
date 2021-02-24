package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前菜单选择
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {}

		//根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式选择...")
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式选择...")
			break
		case 3:
			//更改同户名
			fmt.Println("更新用户名选择...")
			break

		}
	}
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag	= flag
		return true
	} else {
		fmt.Println(">>>>>>请输入合法范围内的数字<<<<<")
		return false
	}

}

func NewClient(serverIp string, serverPort int) *Client {
	//创建client对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	//连接Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Printf("net.Dial err : %V\n", err)
		return nil
	}

	client.conn = conn

	//返回对象
	return client
}
