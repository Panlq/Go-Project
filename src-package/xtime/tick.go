package xtime

import (
	"fmt"
	"time"
)

func tick() {
	cleanup := time.NewTicker(3 * time.Second)
	defer cleanup.Stop()

	go func() {
		for {
			t := <-cleanup.C
			fmt.Println("current time: ", t)
		}
	}()

	time.Sleep(30 * time.Second)
}
