// 

package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

var redisdb *redis.Client

func initClient() (err error) {
	redisdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6739",
		Password: "",
		DB: 0,   // user default db
	})

	_, err = redisdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

//go env -w GOPROXY=direct,https://goproxy.cn,https://goproxy.io,https://mirrors.aliyun.com/goproxy,https://athens.azurefd.net,https://proxy.golang.org

