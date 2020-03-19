package main

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
)

// type Person struct {
// 	Name   string
// 	Age    int
// 	Gender string
// }

func main() {
	// p1 := Person{
	// 	Name:   "阿拉斯加",
	// 	Age:    19,
	// 	Gender: "male",
	// }
	var a [3]int = [...]int{1,2,3}
	b, err := msgpack.Marshal(a)
	if err != nil {
		fmt.Printf("msgpack marshal failed, err:%v", err)
		return
	}

	// var p2 Person
	var c [3]int
	err = msgpack.Unmarshal(b, &c)
	if err != nil {
		fmt.Printf("msgpack unmarshal failed, err:%v", err)
		return
	}
	fmt.Printf("p2:%#v\n", c)
}

