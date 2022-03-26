package main

import (
	"fmt"
	"io/ioutil"
	"github.com/golang/protobuf/proto"
	"selfcode.me/studygo/series/protobuf_demo/address"
)

func main() {
	var cb address.ContactBook
	
	p1 := address.Person{
		Name: "阿良",
		Gender: address.GenderType_MALE,
		Number: "1315831512",
	}

	fmt.Println(p1)

	cb.Persons = append(cb.Persons, &p1)
	data, err := proto.Marshal(&p1)
	if err != nil {
		fmt.Printf("marshal failed, err:%v\n", err)
		return
	}
	ioutil.WriteFile("./proto.dat", data, 0644)
	data2, err := ioutil.ReadFile("./proto.dat")

	if err != nil {
		fmt.Printf("read file failed, err: %v\n", err)
		return
	}

	var p2 address.Person
	proto.Unmarshal(data2, &p2)
	fmt.Println(p2)
}