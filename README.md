# etcdv3-upsync-proxy

解决在`etcd3.*`注册的服务信息，使用v2版本的接口中读取不到的问题

该项目专门支持 `https://github.com/weibocom/nginx-upsync-module` 使用v2版本请求调用v3版本的数据

 - [x] 基本的读取功能，/v2/keys/path/to/you_service
 - [x] 支持go-zero的服务发现方式转换为upsync的服务发现数据格式
 - [x] 服务参数配置化
 - [x] gRPC 健康检查功能，/health-check/ip:port（只允许检测私有IP段）
 - [ ] 添加指标监控
 - [ ] 添加日志
 - [ ] 实现wait参数控制
 - [ ] 实现recursive参数控制
 - [ ] 实现waitIndex参数控制
