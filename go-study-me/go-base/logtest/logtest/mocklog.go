package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fileObj, err := os.OpenFile("./t.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("open file failed err: %v\n", err)
		return
	}
	log.SetOutput(fileObj)
	for {
		log.Println("这是一个简单的日志")
		time.Sleep(time.Second * 3)
	}
}
