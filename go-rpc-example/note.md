### RPC
Remote Procedure Call 远程过程调用 

完整的RPC框架应包含，负载均衡、服务注册和发现、服务治理等功能，并具有可拓展性，便于流量监控系统等接入。


### Protobuf

Protocol Buffers 是一种与语言、平台无关，可扩展的序列化结构化数据的方法，常用于通信协议，数据存储等，相较于json, xml, 它更小，快捷，简单。
1. IDL编写
2. 生成指定语言的代码
3. 序列化和反序列化

##### protobuf 语法
```protobuf
syntax = "proto3"   // 第一行声明使用proto3语法，不声明默认使用proto2

service SerchService {
    rpc Search (SearchRequest) returns (SearchResponse);
}

message SearchRequest {
    string query = 1;
    int32 page_number = 2;
    int32 result_per_page = 3;
}

message SearchResponse {
    ...
}
```