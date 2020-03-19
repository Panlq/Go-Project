package main

import (
	"fmt"
	"time"
	"selfcode.me/studygo/日志/mylogger"
)

func main() {

	// log := mylogger.NewConsoleLog("fatal")
	// log := mylogger.NewFileHandle("info", "./", "test.log", 3*1024)
	log := mylogger.NewFileAsyncHandle("info", "D:/ZIYUAN/Go/gph/src/selfcode.me/studygo/日志/", "test.log", 8*1024)
	for i := 1; i < 1000; i++ {
		fmt.Printf("执行第%d次----", i)
		log.Info("这是一条简单Info的测试")
		log.Debug("这是一条简单Debug的测试")
		log.Warnning("这是一条简单Warnning的测试")
		log.Error("这是一条简单Error的测试")
		log.Fatal("这是一条简单Fatal的测试")
		time.Sleep(time.Second)
	}
}