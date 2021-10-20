package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"internal/foo"
)

var key string = "name"

func TestInternalPkg(t *testing.T) {
	foo.Hello()
}

func TestWithValue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, key, "one")
	go watch(valueCtx)

	time.Sleep(10 * time.Second)
	fmt.Println("notify to cancel")
	cancel()
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Value(key), "monite exited")
			return
		default:
			fmt.Println(ctx.Value(key), "goroutine watching...")
			time.Sleep(2 * time.Second)
		}
	}
}
