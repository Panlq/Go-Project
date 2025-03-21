package main


import (
	"fmt"
)

// 多路复用 同时从多个channel中取值

func main() {
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case x:= <-ch:
			fmt.Println(x)
		case ch <- i:
		}
	}
}