package main

/*
当一个线程获取读锁之后, 其他线程如果是获取读锁会继续获取锁，如果是获取写锁就会等待，
当一个线程获取写锁之后，其他的线程无论是获取读锁还是写锁都会等待
*/

import (
	"fmt"
	"sync"
	"time"
)

var (
	x int64
	wg sync.WaitGroup
	lock sync.Mutex
	rwlock sync.RWMutex
)

func write() {
	// lock.Lock()  // 加互斥锁
	rwlock.Lock()  // 写锁
	x = x + 1
	time.Sleep(time.Millisecond * 10)  // 模拟读操作耗时10ms
	rwlock.Unlock()
	wg.Done()
}

func read() {
	rwlock.RLock()  // 读锁
	time.Sleep(time.Millisecond)
	rwlock.RUnlock()
	wg.Done()
}

func main() {
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go read()
	}

	wg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
}