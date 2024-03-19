# JOJ
Online Judge


## 开发环境
- centos 7.9
- go 1.20.2

## 项目结构图
```txt
.
├── api                 api层
│   └── v1
├── common
│   ├── jwt             jwt模块
│   ├── log             日志模块
│   ├── request         各种request结构定义
│   └── response        各种response结构定义
├── configs             各种组件配置
├── deployments
│   └── docker          docker
├── docs                swaggo文档
├── internal
│   ├── dao             dao层
│   ├── middleware      中间件
│   ├── model           模型层
│   ├── router          路由配置
│   └── service         服务层
└── logs
```

## 功能模块

1. user模块
    - 用户注册
    - 用户登录
    - 用户管理
        - 密码修改

2. problem模块

3. Judge模块

4. contest模块


## 技术栈

1. docker部署
2. 后端框架：Gin Gorm
3. 传输协议：protobuf
4. 日志：zap+lumberjack
5. 配置管理：viper
6. 数据库：MySQL