## Context的本质
在Goroutine构成的树形结构中对信号进行同步以减少计算资源的浪费是Context最大的作用。
go服务每一个请求都会起一个goroutine，处理逻辑中也会起新的goroutine访问数据库或其他服务
Context就可用在不同的goroutine之间同步特定的数据，取消信号以及处理请求的截至日期，类似一个全局变量。


## Usage
context.Backgroud, context.TODO 都是通过new(emptyCtx)初始化的, 指向私有结构体`context.emptyCtx`的指针, 两个变量互为别名，没有实际的功能，在使用和语义有一点不同
- context.Background: 一般作为main函数, 等上层，初始上下文向下传递
- context.TODO: 当不确定使用那个Context, 或者暂时还没有可用的ctx传递时, 可用先使用TODO占位