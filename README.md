# GinGateway

**共琢** 使用Gin框架实现的API网关系统

## 目录

- [GinGateway](#GoZGateway)
    - [目录](#目录)
    - [日志](#日志)
        - [6.3 记录](#6.3-记录)
        - [6.5 记录](#6.5-记录)

## 日志

### 6.3 记录
现在只实现了基本的鉴权、限流和路由转发功能。路由转发通过yaml配置文件进行配置，而限流则有乐观锁(CAS)和悲观锁sync.Lock两种实现。接下来会找时间实现服务注册和发现中心

### 6.5 记录
目前实现的功能有
1. 路由转发
2. 服务发现和注册
3. 限流
4. 权限鉴定

微服务网关的雏形基本上完成了😎
