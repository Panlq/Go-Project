// package main


import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var once sync.Once
// 将双向通道转为单向通道时可以的，但是不能由单向通道转为双向通道

// 单向通道案例
func counter(out chan<- int) {
	// defer wg.Done()
	for i := 0; i < 100; i++ {
		out <- i
	}
	close(out)
}


func squarer(out chan<- int, in <-chan int) {
	// defer wg.Done()
	for i := range in {
		out <- i*i
	}

	once.Do(func() {close(out)}) // 确保某个操作只执行一次
}

func printer(in <-chan int){
	for i := range in {
		fmt.Println(i)
	}
}


func main() {
	// 
	ch1 := make(chan int)  //无缓冲通道 也被称为 同步通道
	ch2 := make(chan int)
	// wg.Add(2)
	go counter(ch1)
	go squarer(ch2, ch1)
	// wg.Wait()
	printer(ch2)
}