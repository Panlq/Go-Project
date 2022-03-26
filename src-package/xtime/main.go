package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	m1, _ := time.ParseDuration("-24h")
	t2 := now.Add(m1)
	fmt.Println(t2)

	duration := now.Sub(t2)
	fmt.Println(duration.Hours())
	fmt.Println(now)
}
