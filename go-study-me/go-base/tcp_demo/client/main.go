package main

import (
	"fmt"
	// "os"
	"net"
	"selfcode.me/studygo/tcp_demo/proto"
	// "bufio"
	// "strings"
)

// func main() {
// 	// 1. 与server 建立连接
// 	conn, err := net.Dial("tcp", "127.0.0.1:20000")
// 	if err != nil {
// 		fmt.Println("connect err", err)
// 		return
// 	}

// 	defer conn.Close()
// 	buf := [512]byte{}
// 	inputReader := bufio.NewReader(os.Stdin)
// 	for {
// 		fmt.Printf("QQ1:")
// 		input, _ := inputReader.ReadString('\n')
// 		inputInfo := strings.Trim(input, "\r\n")
// 		if strings.ToUpper(inputInfo) == "Q" {
// 			break
// 		}
// 		_, err := conn.Write([]byte(inputInfo))
// 		if err != nil{
// 			return
// 		}
// 		n, err := conn.Read(buf[:])
// 		if err != nil {
// 			fmt.Println("recv failed, err", err)
// 			return
// 		}
// 		fmt.Println("recv from server msg:", string(buf[:n]))
// 	}
// }


// 测试粘包现象 
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		data, err := proto.Encode(msg)
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			return
		}
		fmt.Println("send %d msg", i)
		conn.Write(data)
	}
}