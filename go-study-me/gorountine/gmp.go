package gmp

import "crypto/tls"

// g 代表一个goroutine对象, 每次go调用的时候 会创建一个G对象 
type g struct {
	stack stack // 描述真实的栈内存, 包括上下界
	m 	*m // 当前的M
	sched   gobuf  //goroutine切换时, 用于保存g的上下文
	param   unsafe.Pointer　// 用于传递参数，睡眠时其他gorountine可以设置param, 唤醒时该goroutine可以获取
	atomicstatus  uint32
	stackLock     uint32
	goid		int64 //goroutine的ID
	waitsince   int64 // g被阻塞的大致时间
	lockedm     *m    // G被锁定只在这个m上运行
}

// 其中最主要的是sched 保存了goroutine的上下文, goroutine切换的时候不同于线程有OS来负责上下文切换,　而是由一个gobuf对象来保存, 这样可以更加轻量级

type gobuf struct {
	sp uintptr    // 栈指针
	pc uintptr	  // 计数器
	g  guintptr   // 记录自身g 的指针是为了能更快的访问goroutine中的信息
	ctxt unsafe.Pointer
	ret sys.Uintreg
	lr  uintptr
	bp  uintptr
}


type m struct {
	g0 		*g    // 带有调度栈的 goroutine
	gsignal *g    // 处理信号的goroutine
	tls		[6]uintptr // thread-local storage
	mstartfn func()
	curg 	*g    // 当前绑定的结构体G, 即运行的goroutine
	caughtsig guintptr
	p 		puintptr // 关联p和执行的go代码
	nextp   puintptr
	id      int32
	mallocing int32  // 状态
	spinning  bool  // m是否out of work
	blocked   bool  // m是否被阻塞
	inwb	  bool  // m是否被执行写屏蔽
	printlock  int8
	incgo     bool  // m在执行cgo吗
	fastrand  uint32 
	ncgocall  uint64  // cgo调用总数
	ncgo      int32
	park 	  note
	alllink   *m   // 用于连接allm
	schedlink  muintptr
	mcache    *mcache   // 当前m的内存缓存
	lockedg   *g  // 锁定g在当前m上执行， 而不会切换到其他m
	createstack [32]uintptr   // thread 创建的栈
}

//普通的goroutine 的栈是在堆上分配的可增长的栈, 而g0的栈是M对应的线程的栈, 所有调度相关代码, 会先切换到该goroutine的栈中再执行
//也就是说线程的栈也是用g实现, 而不是使用的OS的
