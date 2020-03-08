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

### Go gin FileSystem

go 对filesystem进行了封装只需要在 router层增加一行代码即可开启

> r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))