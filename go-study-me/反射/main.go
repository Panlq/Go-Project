package main

import (
	"fmt"
	"reflect"
)

func reflectType(x interface{}) {
	v := reflect.TypeOf(x)
	// switch res := x.(type) {
	// case float32:
	// 	fmt.Printf("type: %v\n", res)
	// }
	fmt.Printf("type: %v\nkind: %v\n", v.Name(), v.Kind())
}


type student struct {
	Name string `json:"name"`
	Score int 	`json:"score"`
}

// type student struct {
// 	Name string
// 	Score int 	
// }

func main() {
	// var a float32 = 3.14
	// reflectType(a)

	// var b int64 = 1221
	// reflectType(b)

	// var c *int
	// fmt.Println("var a *int IsNil: ", reflect.ValueOf(c).IsNil())
	// fmt.Println("nil IsValid: ", reflect.ValueOf(nil).IsValid())

	// d := struct{}{}
	// fmt.Println("不存在的结构体成员", reflect.ValueOf(d).FieldByName("abc").IsValid())
	// fmt.Println("结构体中不存在该方法", reflect.ValueOf(d).MethodByName("abc").IsValid())

	// e := map[string]int{}
	// //
	// h := reflect.ValueOf("nafsfd")
	// fmt.Printf("h %T\n%v\n", h, h)
	// fmt.Println("map中不存在该key", reflect.ValueOf(e).MapIndex(reflect.ValueOf("lf")).IsValid())

	stu1 := student{
		Name: "小王子",
		Score: 90,
	}
	t := reflect.TypeOf(stu1)
	fmt.Println(t.Name(), t.Kind())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("name:%s index:%d type:%v json:%v\n", field.Name, field.Index, field.Type, field.Tag.Get("json"))
	}
	if scoreField, ok := t.FieldByName("Score"); ok{
		fmt.Printf("name:%s index:%d type:%v json:%v\n", scoreField.Name, scoreField.Index, scoreField.Type, scoreField.Tag.Get("json"))
	}
}