# Service Init Tool

> 🚀 企业级 Go 微服务配置管理框架 - 解决多环境配置冲突、敏感信息泄露、配置漂移等生产痛点

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## 💡 核心特性

- **多源配置融合** - 智能合并环境配置文件、本地配置、环境变量，一键搞定
- **配置优先级** - 内置清晰的优先级机制，本地开发与生产环境无缝切换
- **自动模板生成** - 基于结构体自动生成环境变量模板，告别手写配置
- **开箱即用组件** - 内置 MySQL、Redis、Prometheus、Elasticsearch 等常用组件

## 📦 快速开始

### 安装

```bash
go get github.com/lie-flat-planet/service-init-tool@latest
```

### 基础用法

```go
package main

import (
    service_init_tool "github.com/lie-flat-planet/service-init-tool"
    "github.com/lie-flat-planet/service-init-tool/component/database"
    "github.com/lie-flat-planet/service-init-tool/component/redis"
)

type Config struct {
    Server *service_init_tool.Server
    Mysql  *database.Mysql
    Redis  *redis.Redis
    AppName string `env:""`
    Port    int    `env:""`
}

var Setting = &Config{
    Server: &service_init_tool.Server{Name: "my-service"},
    Mysql:  &database.Mysql{/* 默认配置 */},
    Redis:  &redis.Redis{/* 默认配置 */},
    AppName: "demo",
    Port:    8080,
}

func main() {
    // 自动注入配置并初始化组件
    if err := service_init_tool.Init("./", Setting); err != nil {
        panic(err)
    }
    // Setting 已完成配置注入，可直接使用
}
```

### 配置优先级

配置按以下优先级合并（数字越小优先级越高）：

```
1. local.yml      ← 最高优先级（本地开发）
2. 环境变量        ← 生产环境推荐
3. hot-fix.yml    ← 热修复配置
4. staging.yml    ← 预发布环境
5. test.yml       ← 测试环境
```

### 自动生成配置模板

执行 `Init()` 后会自动生成 `dev.yml` 模板文件：

```yaml
# dev.yml - 自动生成的环境变量参考
AppName: demo
Port: 8080
Mysql_Host: 127.0.0.1:3306
Mysql_User: root
Redis_Host: 127.0.0.1:6379
```

创建 `local.yml` 覆盖默认配置：

```yaml
# local.yml - 本地开发配置
Port: 3000
Mysql_Password: your_password
Redis_Password: your_password
```

## 🔧 内置组件

| 组件              | 功能                          | 说明                        |
| ----------------- | ----------------------------- | --------------------------- |
| **Database**      | MySQL、PostgreSQL、ClickHouse | 基于 GORM，支持连接池配置   |
| **Redis**         | Redis 客户端                  | 基于 go-redis，支持集群模式 |
| **Prometheus**    | 指标监控                      | 开箱即用的 metrics 端点     |
| **Elasticsearch** | ES 客户端                     | 支持索引管理和查询          |
| **Logger**        | 日志系统                      | 基于 Logrus 的结构化日志    |
| **HTTP Server**   | Web 服务                      | 基于 Gin 的 HTTP 服务器     |

## 🎯 设计理念

1. **约定优于配置** - 合理的默认值，最小化配置代码
2. **渐进式扩展** - 按需引入组件，不强制依赖
3. **生产级可靠** - 自动处理连接池、优雅关闭等生产细节
4. **开发体验优先** - 本地开发与生产环境配置分离，互不干扰

## ⚠️ 注意事项

- MySQL、Redis、Prometheus 等组件在 `Init()` 时会建立连接，请确保连接信息可用
- 生产环境推荐使用环境变量注入配置，避免敏感信息泄露
- `local.yml` 应添加到 `.gitignore`，避免提交到代码仓库

## 📚 示例

完整示例参见 [example/service](./example/service) 目录。

运行测试查看效果：

```bash
cd example/service
go test -v
```
