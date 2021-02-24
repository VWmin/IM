package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前菜单选择
}

//处理Server返回的消息，直接显示
func (client *Client) DealResponse() {
	//一旦有数据，直接拷贝到标准输出
	io.Copy(os.Stdout, client.conn)
	//等价于
	//for true {
	//	buf := make()
	//	client.conn.Read(buf)
	//	fmt.Println(string(buf))
	//}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {
		}

		//根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			client.PublicChat()
			break
		case 2:
			//私聊模式
			client.PrivateChat()
			break
		case 3:
			//更改同户名
			for !client.UpdateName() {
			}
			break

		}
	}
}

func (client *Client) PublicChat() {
	// 提示用户输入消息
	var msg string
	fmt.Println("请输入要发送的消息，输入\"!q\"以退出：")
	for true {
		fmt.Scanln(&msg)
		if msg == "!q" {
			return
		}
		//发送给服务器
		_, err := client.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("conn.Write err : ", err)
			return
		}
	}
}

func (client *Client) SearchOnline(){
	//向服务器查询当前在线用户列表
	_, err := client.conn.Write([]byte("who\n"))
	if err != nil {
		fmt.Println("conn.Write err : ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	client.SearchOnline()

	var user string
	var msg string

	fmt.Println("请输入要发送到的用户，输入\"!q\"以退出：")
	fmt.Scanln(&user)
	if user == "!q" {
		return
	}
	fmt.Println("请输入要发送的消息，输入\"!q\"以退出：")
	fmt.Scanln(&msg	)
	if msg == "!q" {
		return
	}

	//发送给服务器
	_, err := client.conn.Write([]byte("to" + "=" + user + "=" + msg + "\n"))
	if err != nil {
		fmt.Println("conn.Write err : ", err)
		return
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename=" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err : ", err)
		return false
	}
	return true

}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
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
