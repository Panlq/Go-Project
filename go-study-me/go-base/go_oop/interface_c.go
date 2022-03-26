package go_oop

import (
	"fmt"
)

type animal interface {
	move()
}

type cat struct {
	name string
	feet int8
}

func (c cat) move() {
	fmt.Println("走猫步")
}

type dog struct {
	name string
	feet int8
}

func (d dog) move() {
	fmt.Println("狗跑....")
}

func main() {
	var a1 animal

	b := cat{
		name: "tom",
		feet: 4,
	}

	c := dog{
		name: "jerry",
		feet: 4,
	}

	a1 = b
	a1 = c
	fmt.Println(a1)
}
