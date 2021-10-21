package test

import (
	"context"
	"fmt"
	"math/rand"
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

func TestWithTimeout(t *testing.T) {
	subCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for {
		select {
		case <-subCtx.Done():
			fmt.Printf("query data status failed, %s", subCtx.Err())
			return
		case <-time.After(2 * time.Second):
			if ok := query_data_status(subCtx); ok {
				fmt.Println("query data status ok")
				return
			}
		}
	}
}

func query_data_status(ctx context.Context) bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(9) > 4 {
		return true
	}

	return false
}

func TestWithDeadline(t *testing.T) {
	dt := time.Now().Add(10 * time.Second)
	subCtx, cancel := context.WithDeadline(context.Background(), dt)

	defer cancel()

	go handler(subCtx)

	select {
	case <-subCtx.Done():
		fmt.Println("main", subCtx.Err())
	}
}

func handler(ctx context.Context) {
	duration := 2 * time.Second
	select {
	case <-ctx.Done():
		fmt.Println("handler", ctx.Err())
	case <-time.After(duration):
		// do somethind done
		fmt.Println("handler", "do something done with", duration)
	}
}
