## Context的本质
在Goroutine构成的树形结构中对**信号**进行**同步**（**同步阻塞通道**）以减少计算资源的浪费是Context最大的作用。
go服务每一个请求都会起一个goroutine，处理逻辑中也会起新的goroutine访问数据库或其他服务
Context就可用在不同的goroutine之间同步特定的数据，取消信号以及处理请求的截至日期，类似一个全局变量。

![image-20211021141957283](https://gitee.com/jonpan/mypic/raw/master/mic/202110211726421.png)


## Context 的继承衍生
context.Backgroud, context.TODO 都是通过new(emptyCtx)初始化的, 指向私有结构体`context.emptyCtx`的指针, 两个变量互为别名，没有实际的功能，在使用和语义有一点不同
- context.Background: 一般作为main函数, 等上层，初始上下文向下传递
- context.TODO: 当不确定使用那个Context, 或者暂时还没有可用的ctx传递时, 可用先使用TODO占位


```golang
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)    // 可用来当作多个goroutine中的同步channel
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key, val interface{}) Context

```

## 应用场景

1. WithTimeout

超时函数

```golang
func TestWithTimeout(t *testing.T) {
	subCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for {
		select {
		case <-subCtx.Done():
			fmt.Printf("query data status failed, %s", subCtx.Err())
			return
		case <-time.After(2 * time.Second):
			if ok := query_data_status(subCtx); ok {
				fmt.Println("query data status ok")
				return
			}
		}
	}
}

func query_data_status(ctx context.Context) bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(9) > 4 {
		return true
	}

	return false
}
```

2. WithDeadline

延时函数

```golang
func TestWithDeadline(t *testing.T) {
	dt := time.Now().Add(10 * time.Second)
	subCtx, cancel := context.WithDeadline(context.Background(), dt)

	defer cancel()

	go handler(subCtx)

	select {
	case <-subCtx.Done():
		fmt.Println("main", subCtx.Err())
	}
}

func handler(ctx context.Context) {
	duration := 2 * time.Second
	select {
	case <-ctx.Done():
		fmt.Println("handler", ctx.Err())
	case <-time.After(duration):
		// do somethind done
		fmt.Println("handler", "do something done with", duration)
	}
}
```

3. WithValue

   这个上下文平时比较少用，一般都是用的框架里封装后的ctx，比如`gin.Context` 

    在真正使用`WithValue` 传值的功能时我们也应该非常谨慎，使用 [`context.Context`](https://draveness.me/golang/tree/context.Context) 传递请求的所有参数是一种非常差的设计，比较常见的使用场景是传递请求对应用户的认证令牌以及用于进行分布式追踪的请求 ID等必要参数.

```golang
func TestWithValue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, key, "one")
	go watch(valueCtx)

	time.Sleep(10 * time.Second)
	fmt.Println("notify to cancel")
	cancel()
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Value(key), "monite exited")
			return
		default:
			fmt.Println(ctx.Value(key), "goroutine watching...")
			time.Sleep(2 * time.Second)
		}
	}
}
```



## WHY

1. 结构体内嵌接口有什么用？

```go
type cancelCtx struct {
	Context

	mu       sync.Mutex            // 同步锁, 保护以下字段的一致性
	done     chan struct{}         // 用的时候才创建, 首次cancel调用是关闭channel
	children map[canceler]struct{} //
	err      error
}
```

- **结构体嵌套结构体**

  Go 中在 **结构体 A** **嵌套**另外一个 **结构体 B** 见的很多，通过嵌套，可以扩展A的能力。

  A不仅拥有了B的属性，还拥有了B的方法，这里面有一个[**字段提升**](https://gfw.go101.org/article/type-embedding.html)的概念。

  

- **结构体嵌套接口**

  结构体里嵌套接口的**目的**：

  当前结构体实例**可以用所有实现了该接口的其他结构体来初始化**（即使他们的属性不完全一致）



2. 针对非衍生类型的context处理, 两个case分别有什么用, 为什么要监听`child.Done()` [code]()

```go
atomic.AddInt32(&goroutines, +1)
go func() {
    select {
        case <-parent.Done():
        	child.cancel(false, parent.Err())
        case <-child.Done():
    }
}()
```

- 监听`parent.Done()` : 如父节点是`ctx := gin.Context` 当`ctx.Done()` 返回了非空的Channel时，则需要级联取消子节点

- 监听`child.Done()`: 如果派生的类似`WithCancel` 的子节点主动取消了，则需要退出`select` ，**避免`goroutine`泄露**

## 参考

1. [类型内嵌](https://gfw.go101.org/article/type-embedding.html)
2. [GO语言设计与实现-上下文](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-context/#61-%E4%B8%8A%E4%B8%8B%E6%96%87-context)
3. [飞雪无情-context](https://www.flysnow.org/2017/05/12/go-in-action-go-context.html#context%E6%8E%A5%E5%8F%A3)