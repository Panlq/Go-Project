#### learn from eddycjy
[煎蛋大佬博客](https://book.eddycjy.com/golang)

#### docker
- FROM 指定基础镜像 必须有的指令, 且必须是第一条指令
- WORKDIR <工作目录路径> 
- COPY 
    格式:
    COPY <source> <dst>
    COPY ["<源路径1>",..."<目标路径>"]

- RUN <命令>  // 执行命令行命令
- EXPOSE [<端口1>, <端口2>] 指令是声明运行时容器提供服务端口，这只是一个声明，在运行时并不会因为这个声明应用就会开启这个端口的服务
    - 帮助镜像使用者理解这个镜像服务的守护端口，以方便配置映射
    - 运行时使用随机端口映射时，也就是 docker run -P 时，会自动随机映射 EXPOSE 的端口
- ENTRYPOINT 指定容器启动程序及参数
    exec格式: <ENTRYPOINT> "<CMD>"
    shell 格式: ENTRYPOINT ["curl", "-s", "http://ip.cn"]

### GORM 
> GORM itself is powered by Callbacks, so you could fully customize GORM as you want
- 注册一个新的回调
- 删除现有的回调
- 替换现有的回调
- 注册回调顺序

在go中 当存在多个init函数时，执行顺序为:

- 相同包下的init函数: 按照源文件编译顺序决定执行顺序（默认按文件名排序）
- 不同包下的init函数: 按照包导入的依赖关系决定先后顺序

### 编写models callbacks
gorm的Callbacks 可以将回调方方法定义为模型结构的指针，在创建，更新，查询，删除时将被调用

如果任何回掉返回错误，gorm将停止未来操作并回滚所有更改。

gorm支持的回调方法:

- create: BeforeSave, BeforeCreate, AfterCreate, AfterSvae
- update: BeforeSave, BeforeUpdate, AfterUpdate, AfterSave
- delete: BeforeDelete, AfterDelete
- select: AfterFind

gorm 关联属性方式(外键)

- gorm 会通过类名+ID的方式去找到两个类之间的关联关系， 可以通过Related 进行查询

```golang
var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

type Tag struct {
	Model

	Name       string `josn:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

type Article struct {
    Model

    TagID int `json:"tag_id" gorm:"index"`
    Tag   Tag `json:"tag"`

    Title string `json:"title"`
    Desc string `json:"desc"`
    Content string `json:"content"`
    CreatedBy string `json:"created_by"`
    ModifiedBy string `json:"modified_by"`
    State int `json:"state"`
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)
}
```

Preload 是一个预加载器， 会执行两条SQL, 

1. SELECT * FROM blog_articles;
2. SELECT * FROM blog_tag WHERE id IN (1, 2, 3, 4)

查询出结构后，gorm内部处理对应映射逻辑，将其填充到对应的结构体中，可避免循环查询

一般方法:

- gorm的Join

- 循环Related

#### 文件操作

- os.Stat: 返回文件信息结构描述文件，如果出现错误，会返回*PathError

  ```golang
  type PathError struct {
  	Op string
  	Path string
  	Err error
  }
  ```

- os.IsNotExist: 可以判断目录或者文件不存在

- os.IsPermission: 查看是否有权限操作

- os.OpenFile: 调用文件 传入文件名称, 指定的模式调用文件，文件权限，返回的文件的方法可以用I/O

  ```golang
  const (
      // Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
      O_RDONLY int = syscall.O_RDONLY // 以只读模式打开文件
      O_WRONLY int = syscall.O_WRONLY // 以只写模式打开文件
      O_RDWR   int = syscall.O_RDWR   // 以读写模式打开文件
      // The remaining values may be or'ed in to control behavior.
      O_APPEND int = syscall.O_APPEND // 在写入时将数据追加到文件中
      O_CREATE int = syscall.O_CREAT  // 如果不存在，则创建一个新文件
      O_EXCL   int = syscall.O_EXCL   // 使用O_CREATE时，文件必须不存在
      O_SYNC   int = syscall.O_SYNC   // 同步IO
      O_TRUNC  int = syscall.O_TRUNC  // 如果可以，打开时
  )
  ```

- os.Getwd:  返回当前目录对应的根路径名

- os.MkdirAll: 创建对应的目录以及所需的子目录，成功返回nil， 否则返回error

- os.ModePerm: const定义 ModePerm FileMode = 0777

### Go gin FileSystem

go 对filesystem进行了封装只需要在 router层增加一行代码即可开启

> r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

Go 读写excel 第三方库

- [tealeg/xlsx](https://github.com/tealeg/xlsx)
- [360EntSecGroup-Skylar/excelize](https://github.com/360EntSecGroup-Skylar/excelize)



#### Go 生成二维码，合并海报

没做错误日志记录



#### Makefile

规则

Makefile 由多条规则组成，每条规则都以一个 target（目标）开头，后跟一个 : 冒号，冒号后是这一个目标的 prerequisites（前置条件）

```makefile
[target] ...: [prerequisites]...
	[command]
	...
	...
```

- target: 一个目标代表一条规则，可以是一个或多个文件名，也可以是某个操作的名字(标签)，称为伪目标（phony）
- prerequisites: 前置条件，这一项是可选参数，通常是多个文件名、伪目标.。它的作用是 target 是否需要重新构建的标准，如果前置条件不存在或有过更新（文件的最后一次修改时间）则认为 target 需要重新构建
- command: 构建一个target的具体命令集

```makefile
.PHONY: build clean tool lint help

all: build

build:
	go build -v .
	
tool:
	go tool vet . |& grep -v vendor; true
	gofmt -w .
	
lint:
	golint ./ ...
	
clean:
	rm -rf go-gin-example
	go clean -i .
	
help:
	@echo "make: copile packages and dependencies"
	@echo "make tool: run sepcitie
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"
```

