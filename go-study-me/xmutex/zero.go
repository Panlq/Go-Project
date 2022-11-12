package xmutex

func cas(val *int32, old, new int32) bool
func semacquire(*int32)
func semrelease(*int32)

type Mutex struct {
	key  int32 // 锁是否被持有
	sema int32 // 信号量专用，用以阻塞/唤醒goroutine
}

func xadd(val *int32, delta int32) (new int32) {
	for {
		v := *val
		if cas(val, v, v+delta) {
			return v + delta
		}
	}
	panic("unreached")
}

func (m *Mutex) Lock() {
	// 如果等于1， 成功获取到锁
	if xadd(&m.key, 1) == 1 {
		return
	}

	semacquire(&m.sema) // 否则阻塞等待
}

func (m *Mutex) Unlock() {
	if xadd(&m.key, -1) == 0 {
		return
	}

	semrelease(&m.sema)
}
