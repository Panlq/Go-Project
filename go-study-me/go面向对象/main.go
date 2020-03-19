// package main

import (
	"fmt"
)

//Animal 动物
type Animal struct {
	name string
}


func (a *Animal) move() {
	fmt.Printf("%s会动!\n", a.name)
}


//Dog 狗
type Dog struct {
	Feet int8
	*Animal // 通过嵌套匿名结构体实现继承
}

func (d *Dog) wang() {
	fmt.Printf("%s 汪汪\n", d.name)
}


func main() {
	d1 := &Dog{
		Feet: 4,
		Animal: &Animal{   // 这里嵌套的是结构体指针
			name: "lele",
		},
	}
	d1.Animal.move()
	d1.wang()
	d1.move()
}