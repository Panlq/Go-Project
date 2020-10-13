/*
 * @Author: panlq01@mingyuanyun.com
 * @Date: 2020-10-12 16:44:49
 * @Description: Some desc
 * @LastEditors: panlq01@mingyuanyun.com
 * @LastEditTime: 2020-10-12 16:55:20
 */
package main

import (
	"fmt"
	"go-gin-ex/pkg/setting"
	"go-gin-ex/routers"
	"net/http"
)

func main() {
	router := routers.InitRouter()

	ser := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	ser.ListenAndServe()
}
