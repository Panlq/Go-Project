package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const testScript = `local ttl_window = tonumber(ARGV[1])
if redis.call("EXISTS", KEYS[1]) == 1 then
    redis.call("expire", KEYS[1], ttl_window)
	local eid = redis.call("hget", KEYS[1], "name")
	redis.call("hset", KEYS[1], "age", 18)
	return eid
else
	return ""
end`

func main() {
	host := "localhost"
	port := 6379
	db := 0
	pwd := ""
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Password:     pwd,
		PoolSize:     4,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	hashMap := map[string]interface{}{
		"name":   "张三",
		"gender": "male",
		"age":    15,
		"job":    "dba",
	}

	key := "test:lua:script"
	if _, err := client.HMSet(key, hashMap).Result(); err != nil {
		log.Fatal(err)
	}

	val, err := client.HGet(key, "name").Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(val)

	age, err := client.HGet(key, "age").Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v, %T", age, age)

	resp, err := client.Eval(testScript, []string{key}, strconv.Itoa(80)).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)

	ttl, err := client.TTL(key).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ttl)
}
