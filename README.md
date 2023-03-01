# etcdv3_upsync_proxy

etcd 3.5 注册的服务在v2版本的接口中读取不到

该项目专门支持 `https://github.com/weibocom/nginx-upsync-module` 使用v2版本调用v3的数据

 - [x] 基本的读取功能，/v2/keys/path/to/you_service
 - [x] 支持go-zero的服务发现方式转换为upsync的服务发现数据格式
 - [ ] 服务参数配置化
 - [ ] 实现wait参数控制
 - [ ] 实现recursive参数控制
 - [ ] 实现waitIndex参数控制
