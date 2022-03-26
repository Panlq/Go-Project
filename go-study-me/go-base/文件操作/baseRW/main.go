package main

// 文件操作

import (
	"io"
	// "io/ioutil"
	"os"
	"fmt"
	"bufio"
)

func main() {

	// 使用io/ioutil读取文件
	// content, err := ioutil.ReadFile("C:/Users/asus/Desktop/Go/编程基础/funcion.go")
	// if err != nil {
	// 	fmt.Println("open file failed ,err:", err)
	// 	return
	// }
	// fmt.Println(string(content))
	// os模块读取文件
	fmt.Println(os.Getwd())
	file, err := os.Open("C:/Users/asus/Desktop/Go/编程基础/funcion.go")
	if err != nil {
		fmt.Println("open file failed ,err:", err)
		return
	}

	defer file.Close()

	// // 使用bufio包读取文件
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(line) != 0 {
				fmt.Println(line)
			}
			fmt.Println("finally")
			break
		}

		if err != nil {
			fmt.Println("read file failed, err:", err)
			return
		}
		fmt.Print(line)
	}
	// // 循环读取
	// var content []byte
	// var tmp = make([]byte, 128)
	// for {
	// 	n, err := file.Read(tmp)
	// 	if err == io.EOF {
	// 		fmt.Println("文件读完了")
	// 		break
	// 	}

	// 	if err != nil {
	// 		fmt.Println("read file failed, err", err)
	// 		return
	// 	}
	// 	content = append(content, tmp[:n]...)
	// }

	// fmt.Println(string(content))
}