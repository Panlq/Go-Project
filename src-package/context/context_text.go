package context

import (
	"fmt"
	"testing"
	"time"
)

var key string = "name"

func WithValueTest(t *testing.T) {
	ctx, cancel := WithCancel(Background())
	valueCtx := WithValue(ctx, key, "one")
	go watch(valueCtx)

	time.Sleep(10 * time.Second)
	fmt.Println("notify to cancel")
	cancel()
	time.Sleep(5 * time.Second)
}

func watch(ctx Context) {
}
