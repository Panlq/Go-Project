package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type job struct {
	value int64
}

type result struct {
	job *job
	sum int64
}

var wg sync.WaitGroup
var once sync.Once

var jobChan = make(chan *job, 100)
var resultChan = make(chan *result, 100)

func genJob(out chan<- *job) {
	// 循环生成int64类型随机数
	for {
		x := rand.Int63()
		newJob := &job{
			value: x,
		}
		out <- newJob
		time.Sleep(time.Millisecond * 500)
	}
}

func handle(in <-chan *job, out chan<- *result) {
	// 从jobChan中去除随机数计算各位数的和，将结果发送 resultChan
	for j := range in {
		sum := int64(0)
		n := j.value
		for n > 0 {
			sum += n % 10
			n = n / 10
		}
		newResult := &result{
			job: j,
			sum: sum,
		}
		out <- newResult
	}
}

func main() {
	wg.Add(1)
	go genJob(jobChan)
	for i := 0; i <= 24; i++ {
		go handle(jobChan, resultChan)
	}

	// 主goroutine从resultChan中获取结果并打印到终端输出
	for result := range resultChan {
		fmt.Printf("value:%d sum:%d\n", result.job.value, result.sum)
	}
	wg.Wait()
}
