# Go 并发实战

## Mutex

同步原语适用场景

- 共享资源：并发地读写共享资源，会出现数据竞争(data race)的问题，所以需要`Mutex`,`RWMutex` 这样的并发原语来保护
- 任务编排：需要goroutine按照一定的规律运行，而goroutine之间有相互等待或者依赖的顺序关系，我们常常使用WaitGroup或者channel来实现
- 消息传递：信息交流以及不同的goroutine之间的线程安全数据交流，常常使用channel实现



从复杂度，性能，结构设计来认识golang mutex 同步原语的演进



1. 为什么mutex声明后不需要初始化，零值是 还没有goroutine等待的未加锁状态？
2. 如果 Mutex 已经被一个 goroutine 获取了锁，其它等待中的 goroutine 们只能一直等待。那么，等这个锁释放后，等待中的 goroutine 中哪一个会优先获取 Mutex 呢？





unsafe模块学习