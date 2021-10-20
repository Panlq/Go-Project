## go src package learning

- [x] context




## some issue
#### use of internal package internal/unsafeheader not allowed

阅读源码的过程如果直接引用源码中的库会遇到

在 Go 1.4 及后续版本中，可以通过创建 [Internal packages](https://golang.google.cn/doc/go1.4#internalpackages) 代码包让一些程序实体只能被当前模块中的其他代码引用
**规则如下:**

- internal 代码包中声明的公开程序实体仅能被该代码包的**直接父包及其子包中的代码引用**
- 对于其他代码包，导入该 internal 包都是非法的，无法通过编译
- 名称必须是internal

如果想引用internal文件夹中的内容
1. 需要拷贝一份, 并修改相应的导入路径
2. 拷贝一份, 并重命名相应的文件夹, 改为当前目录的私有模块

示例：

```shell
├── context
│   ├── README.md
│   ├── context.go
│   ├── internal
│   │   └── ctx
│   │       └── ctx.go
│   └── test
│       └── context_test.go
├── internal
│   ├── bar
│   │   ├── bar.go
│   │   └── go.mod
│   └── foo
│       ├── foo.go
│       └── go.mod
├── go.mod
├── main.go
├── README.md
```

如上包结构的程序，`internal`文件可被当前包所有子包导入，但是`context/internal/ctx`只能被`contxt`包及其子包中的代码导入, 不能被`main.go` 导入调用

>  **注：**如果想直接写导入`internal/xxx`, 则需要放入`GOPATH` 或者修改`go.mod` [方法](https://stackoverflow.com/questions/33351387/how-to-use-internal-packages)

##### 思考：这个internal包有什么好处吗？

如果项目包含多个包，可能有一些公共的函数，这些函数旨在供项目中的其他包使用，但不打算成为项目的公共API的一部分。 如果你发现是这种情况，那么 `go tool` 会识别一个特殊的文件夹名称 - 而非包名称 - `internal/` 可用于放置对项目公开的代码，但对其他项目是私有的。比如go源码中`internal` 该包就是仅供标准库内的包使用，不可被外部调用。

目前go生态系统中比较常见的项目包布局形式就包含 `cmd`, `internal`, `pkg` 三个基础目录。

参考: [Go 面向包的设计和架构分层](https://github.com/danceyoung/paper-code/blob/master/package-oriented-design/packageorienteddesign.md)



