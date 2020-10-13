/*
 * @Author: panlq01@mingyuanyun.com
 * @Date: 2020-10-12 16:44:49
 * @Description: Some desc
 * @LastEditors: panlq01@mingyuanyun.com
 * @LastEditTime: 2020-10-13 18:44:39
 */
package main

import (
	"fmt"
	"go-gin-ex/models"
	"go-gin-ex/pkg/logging"
	"go-gin-ex/pkg/setting"
	"go-gin-ex/routers"
	"net/http"
)

func main() {

	// endless.DefaultReadTimeOut = setting.ReadTimeout
	// endless.DefaultWriteTimeOut = setting.WriteTimeout
	// endless.DefaultMaxHeaderBytes = 1 << 20
	// endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

	// server := endless.NewServer(endPoint, routers.InitRouter())

	// server.BeforeBegin = func(add string) {
	// 	log.Printf("Actual pid is %d", syscall.Getpid())
	// }

	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Printf("Server err: %v", err)
	// }
	setting.Setup()
	models.Setup()
	logging.Setup()

	router := routers.InitRouter()

	ser := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	ser.ListenAndServe()
}
