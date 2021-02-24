package main

import "fmt"

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>>>> 连接服务器失败...")
		return
	}

	fmt.Println(">>>>>>>>>> 连接服务器成功...")

	//启动客户单的业务
	select {}
}
