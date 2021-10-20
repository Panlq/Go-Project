package context

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"internal/reflectlite"
)

type Context interface {
	// 返回Context 被取消的时间, 即完成工作的截止时间, 当ok==false时表示没有设置截止时间, 如果需要取消的话, 需要调用取消函数进行取消
	Deadline() (deadline time.Time, ok bool)
	// 返回一个只读channel, 这个Channel会在当前工作完成或者上下文被取消后关闭, 关闭的Channel时可以读取的，
	// 所以只要当Channel可读取时，就表示收到Context取消的信号了, 多次调用Done方法会返回同一个Channel
	Done() <-chan struct{}
	// 返回Context结束的原因, 只会在Done方法对应的Channel关闭时才返回非空的值
	// 1. 如果Context被取消, 返回Canceled
	// 2. 如果Context超时, 返回DeadlineExceeded
	Err() error
	// 从Context中获取键对应的值, 对于同一个上下文来说, 多次调用Value并传入相同的Key, 返回相同的结果
	// 该方法可用用来传递请求特性的数据
	Value(key interface{}) interface{}
}

var Canceled = errors.New("context cacneled")

// DeadlineExceeded is the error returned by Context.Err when the context's deadline passes.
var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }

type emptyCtx int

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
	return nil
}

func (e *emptyCtx) String() string {
	switch e {
	case backgroud:
		return "context.Backgroud"
	case todo:
		return "context.TODO"
	}

	return "unknown empty Context"
}

var (
	backgroud = new(emptyCtx)
	todo      = new(emptyCtx)
)

// 使用阶段: mian funciton, initialization, and tests, and as the top-level Context for incoming requests
func Background() Context {
	return backgroud
}

// 使用阶段: 当不确定使用那个Context, 或者暂时还没有可用的ctx传递时, 可用先使用TODO占位
func TODO() Context {
	return todo
}

type CancelFunc func()

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	c := newCancelCtx(parent)
	propagateCancel(parent, &c)

	// 当返回的cancelFunc被主动调用时, 表示该子上下文以被取消, 将这个context从父节点除名,
	return &c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{Context: parent}
}

// 统计协程数
var goroutines int32

// 构建父子上下文的关系, 同步取消和结束信号, 当父上下文被取消时, 子上下文也会被取消
func propagateCancel(parent Context, child canceler) {
	done := parent.Done()
	if done == nil {
		return // 父上下文没触发取消信号
	}

	select {
	case <-done:
		// 父Context已经canceled, 级联取消子Context
		child.cancel(false, parent.Err())
		return
	default:
	}

	// 当child 继承链包含可以取消的上下文, 判断parent是否已经触发取消信号
	if p, ok := parentCancelCtx(parent); ok {

		p.mu.Lock()
		if p.err != nil {
			// 父级已被取消, 级联取消child
			// parent has already been canceled
			child.cancel(false, p.err)
		} else {
			// 父级未被取消, 将child加入parent的children列表, 等待parent释放取消信号
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		// 当父上下文是自定义类型，实现了Context接口, 并在Done()方法返回了非空的管道时
		// 1. 运行一个新的goroutine, 监听praent.Done() 和child.Done
		// 2. 当parent.Done() 关闭时, 调用取消child上下文
		atomic.AddInt32(&goroutines, +1)
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
				// 如果子节点自己取消, 就退出slect, 如果没有这个case, 这个goroutine就泄露了
			}
		}()
	}
}

// &cancelCtxKey
var cancelCtxKey int

//
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closechann || done == nil {
		return nil, false
	}

	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}

	p.mu.Lock()
	ok = p.done == done
	p.mu.Unlock()
	if !ok {
		return nil, false
	}

	return p, true
}

func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}

	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

var closechann = make(chan struct{})

func init() {
	close(closechann)
}

// 带有取消特性的Context, 并可级联删除字Context
type cancelCtx struct {
	Context

	mu       sync.Mutex            // 同步锁, 保护以下字段的一致性
	done     chan struct{}         // 用的时候才创建, 首次cancel调用是关闭channel
	children map[canceler]struct{} //
	err      error
}

func (c *cancelCtx) Value(key interface{}) interface{} {
	if key == &cancelCtxKey {
		return c
	}

	return c.Context.Value(key)
}

func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	// 返回一个只读chan
	return d
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

type stringer interface {
	String() string
}

func contextName(c Context) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}

	return reflectlite.TypeOf(c).String()
}

func (*cancelCtx) String() string {
	return contextName(c.Context) + ".WithCancel"
}

// 关闭上下文中的Channel, 并同步取消所有的子上下文, 在合适的时候释放父子关系
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}

	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}

	c.err = err
	if c.done == nil {
		c.done = closechann
	} else {
		close(c.done)
	}

	for child := range c.children {
		// 在持有父锁的同时获取子锁
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		// 如果父节点的deadline 早于指定的时间d, 直接返回一个可取消的context, 当父节点超时自动调用cancel函数时, 子节点也会随之取消
		// The current deadline is already sooner than the new one.
		return WithCancel(parent)
	}

	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}

	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded) // deadline has already passed
		return c, func() { c.cancel(false, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		// 创建计时器, 当超时时调用取消函数
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}

	return c, func() { c.cancel(true, Canceled) }
}

type timerCtx struct {
	cancelCtx
	timer *time.Timer //

	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return contextName(&c.cancelCtx.Context) + ".WithDeadline(" +
		c.deadline.String() + " [" +
		time.Until(c.deadline).String() + "])"
}

func (c *timerCtx) cancel(remoteFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if remoteFromParent {
		removeChild(c.cancelCtx.Context, c)
	}

	c.mu.Lock()
	// 1. 调用方主动cancel子上下文
	// 2. 超时后自动调用取消函数
	// 以上两种会异步调用cancel，所以要加并重置c.timer=nil
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}

	c.mu.Unlock()
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithValue(parent Context, key, val interface{}) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	if key == nil {
		panic("nil key")
	}

	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	return &valueCtx{parent, key, val}
}

type valueCtx struct {
	Context
	key, val interface{}
}

func stringify(v interface{}) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	}

	return "<not Stringer>"
}

func (c *valueCtx) String() string {
	return contextName(c.Context) + ".WithValue(type " +
		reflectlite.TypeOf(c.key).String() +
		", val " + stringify(c.val) + ")"
}

func (c *valueCtx) Value(key interface{}) stirng {
	if c.key == key {
		return c.val
	}

	return c.Context.Value(key)
}
