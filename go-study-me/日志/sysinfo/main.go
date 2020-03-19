package main

import (
	"fmt"
	"path"
	"runtime"
)


func getStackInfo(skip int) (funcName, fileName, string, linno int){
	pc, file, linno, ok := runtime.Caller(skip)
	funcName := runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	if !ok {
		fmt.Printf("runtime.Caller() fai ")
	}
	return
}

func main() {
	getStackInfo(1)
}