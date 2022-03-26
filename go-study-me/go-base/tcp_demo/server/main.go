package main
// TCP Server

import (
	"io"
	"fmt"
	"net"
	"bufio"
	"selfcode.me/studygo/tcp_demo/proto"
)

// func process(conn net.Conn) {
// 	defer conn.Close()  // 关闭连接
// 	var buf [128]byte
// 	addr := conn.RemoteAddr()
// 	reader := bufio.NewReader(conn)
// 	for {
// 		n, err := reader.Read(buf[:]) // 读取数据
// 		if err == io.EOF {
// 			fmt.Println("client deadline.")
// 			break 
// 		}
// 		if err != nil {
// 			fmt.Println("read from client failed, err:", err)
// 			break
// 		}
// 		recvStr := string(buf[:n])
// 		fmt.Printf("addr %T, %v\n", addr, addr)
// 		fmt.Printf("recv from client:%s msg:%s\n", addr.String(), recvStr)
// 		conn.Write([]byte(recvStr))  // 发送数据
// 	}
// }

// 使用自定义协议来发送和接收包, 避免粘包现象
func process(conn net.Conn) {
	defer conn.Close()  // 关闭连接
	addr := conn.RemoteAddr()
	reader := bufio.NewReader(conn)
	for {
		recvStr, err := proto.Decode(reader)
		if err == io.EOF {
			fmt.Println("client deadline.")
			break 
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}

		fmt.Printf("recv from client:%s msg:%s\n", addr.String(), recvStr)
	}
}

func main() {
	// 1. 本地端口启动服务
	listener, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("start tcp server on 127.0.0.1:20000")
		return
	}

	defer listener.Close()
	// 2. 等待别人来跟我建立连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("read from conn failed, err", err)
			return
		}
		go process(conn)  // 启动一个goroutine处理连接
	}
}