# 大文件分片上传实现

本项目实现了基于Go Fiber框架的大文件分片上传功能，支持断点续传。

## 实现原理

### 前端实现

1. **文件唯一标识**
   - 使用SparkMD5计算整个文件的hash值作为唯一标识
   - 通过hash值判断文件是否已上传，实现秒传功能

2. **文件分片处理**
   - 使用Blob.slice()方法将文件切分成固定大小的分片
   - 每个分片包含：分片序号、文件hash、分片数据

3. **分批并发上传**
   - 控制并发上传的分片数量，避免服务器压力过大
   - 使用Promise.all管理多个分片的上传任务

4. **断点续传**
   - 上传前先调用check接口获取已上传的分片信息
   - 仅上传未完成的分片，节省带宽和时间
   - 支持上传过程中断后继续上传

### 后端实现

1. **分片管理**
   - 使用文件hash作为目录名，存储分片文件
   - 每个分片使用序号命名，方便合并时排序
   - 使用status文件记录每个分片的上传状态

2. **状态记录**
   - status文件使用二进制格式记录分片状态
   - 每个bit表示一个分片的上传状态
   - 通过状态文件实现断点续传的分片校验

3. **分片合并**
   - 所有分片上传完成后进行合并
   - 合并时按序号顺序读取分片
   - 使用bufio提升合并效率
   - 合并完成后校验文件hash确保完整性

## 目录结构

```
├── data/
│   └── chunks/     # 存储分片文件
├── upload/
│   └── upload.go   # 上传相关处理逻辑
├── web/           # 前端代码
└── main.go       # 服务入口
```

## API接口

### 1. 检查文件状态

```
GET /upload/check

参数：
- hash: 文件hash值
- filename: 文件名
- totalChunks: 总分片数

返回：
{
    "exists": boolean,        // 文件是否已存在
    "uploadedChunks": []int  // 已上传的分片序号
}
```

### 2. 上传分片

```
POST /upload/chunk

参数：
- chunk: 分片文件
- hash: 文件hash值
- index: 分片序号

返回：
{
    "success": boolean
}
```

### 3. 合并分片

```
POST /upload/merge

参数：
- hash: 文件hash值
- filename: 文件名
- size: 文件大小

返回：
{
    "success": boolean
}
```

## 启动服务

1. 启动后端服务
```bash
go run main.go
```

2. 启动前端开发服务器
```bash
cd web
npm install
npm run dev
```

访问 http://localhost:3080 即可使用上传功能。