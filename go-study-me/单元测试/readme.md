
golang 单元测试主要依赖于testing这个内置包
一般在工程 的包目录内 以**_test.go为后缀的源代码文件就是go test 测试的代码, 不会被go build 编译到最终的可执行文件中

*_test.go 文件中一般包含三种类型的函数
- 测试函数  
    **函数名前缀为Test..**, 测试程序的一些逻辑行为是否正确
- 基准函数  
    **函数名前缀为Benchmark..**, 测试函数的性能
- 示例函数  
    **函数名前缀为Example..**, 提供示例文档

### 单元测试
#### 格式
```go
func TestName(t *testing.T){
    //...
}
```
>函数名必须以Test开头，后缀必须以大写字母开头, 最好使用双驼峰法命名函数

### 基准测试用例(检测性能)
##### 使用: 
> go test -bench=Parttern -benchmem

-benchmem 表示开启内存统计  也可以函数内调用b.ReportAllocs() 来表示仅对该测试用例进行内存统计  
##### 格式:
```go 
func BenchmemSplitE(b *testing.B) {
    b.ReportAllocs()
    // .....
}
```

##### 结果说明:
```go
λ go test -v -bench=Split -benchmem
goos: windows
goarch: amd64
pkg: selfcode.me/studygo/split
BenchmarkSplit-4          934220              1276 ns/op          496 B/op    5 allocs/op
PASS
ok      selfcode.me/studygo/split       3.274s
```

- BenchmarkSplit-4 表示电脑有4核, 即b.N == GOMAXPROCS == 4  
- 934220次调用 1276ns/op 表示每次调用Split函数耗时
- 496B/op  表示每次操作内存分配了496字节
- 5allocs/op  表示每次操作进行了5次内存分配

**Benchmark  结构体提供的信息**
```go 
type BenchmarkResult struct {
    N         int           // The number of iterations. 即 b.N
    T         time.Duration // The total time taken. 基准测试花费的时间
    Bytes     int64         // Bytes processed in one iteration. 一次迭代处理的字节数，通过 b.SetBytes 设置
    MemAllocs uint64        // The total number of memory allocations. 总的分配内存的次数
    MemBytes  uint64        // The total number of bytes allocated. 总的分配内存的字节数
}
```

#### 示例函数

##### 格式

```go
func ExampleSplit() {
    fmt.Println(split.Split("a:b:c", ":"))
    fmt.Println(split.Split("锅中有肉中有油中有菜", "中"))
}
// Output:
// [a b c]
// [锅中 有肉 有油 有菜]
```


#### 常用命令参数

- -v   			查看测试函数名称和运行耗时
- -run          -run=Partern 对应一个正则表达式,  指定测试的函数 go test -v -run="More" 包含More的测试函数
- -cover         查看测试覆盖率   即 在测试中至少被运行一次的代码占总代码的比例
- -bench        执行基准测试 测试函数性能
- -benchmem 配合-bench 使用 输出更详细的信息，获取内存分配的统计数据