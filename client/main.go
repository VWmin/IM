package main

import (
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

// 初始化命令行参数
func init() {
	//绑定到变量，参数名，默认值，说明
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址")
	flag.IntVar(&serverPort, "port", 8888, "指定端口号")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>> 连接服务器失败...")
		return
	}

	fmt.Println(">>>>>>>>>> 连接服务器成功...")

	//启动客户单的业务
	select {}
}
