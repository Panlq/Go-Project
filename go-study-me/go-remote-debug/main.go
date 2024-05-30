package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// 初始化一个http server
	mu := http.NewServeMux()

	// 注册一个简单的路由
	mu.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("hello world: %s\n", os.Getenv("LOGNAME"))))
		return
	})

	// 启动http server
	http.ListenAndServe(":8080", mu)
}
