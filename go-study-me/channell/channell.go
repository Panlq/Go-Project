package channell

import (
	"fmt"
	"testing"
	"time"
)

// 关闭的channel是可读的, 可读取直到通道为空, 在取值则返回通道类型零值
func TestCloseChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		v := <-ch
		fmt.Println(v)
		close(ch)
	}()

	ch <- 19

	closeVal := ch

	fmt.Println("channel close ", closeVal)
}
