# Newe - Go语言企业级开发框架

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Gin](https://img.shields.io/badge/Gin-v1.9.1-00ADD8?style=for-the-badge&logo=go)

Newe是一个基于Gin+GORM+Redis构建的企业级Go语言开发框架，提供完整的后台管理系统基础功能。

## ✨ 特性

- 🚀 **高性能**: 基于Gin框架，支持高并发请求
- 🗄️ **ORM支持**: 集成GORM，支持MySQL数据库
- 🔐 **身份认证**: JWT令牌认证，支持RBAC权限控制
- 📊 **数据字典**: 动态数据字典管理
- 📝 **日志系统**: 多级别日志记录，支持文件分割
- 🔄 **WebSocket**: 实时通信支持
- 📁 **文件管理**: 文件上传下载管理
- ⚡ **Redis缓存**: 高性能缓存支持
- 🩺 **健康检查**: 应用健康状态监控
- 🔧 **配置管理**: 灵活的配置文件管理

## 🚀 快速开始

### 环境要求

- Go 1.20+
- MySQL 5.7+
- Redis 6.0+
- Docker (可选，用于容器化部署)

### 安装部署

#### 方式一：传统部署

1. **克隆项目**
```bash
git clone https://github.com/hkyangyi/newe.git
cd newe
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**
```sql
CREATE DATABASE newe CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. **启动应用**
```bash
# 使用Makefile
make run

# 或开发模式
make dev

# 或直接运行
go run examples/serve.go
```

#### 方式二：Docker部署

1. **使用Docker Compose一键部署**
```bash
# 启动所有服务
make compose

# 或手动启动
docker-compose up -d
```

2. **仅构建应用镜像**
```bash
make docker
```

3. **查看服务状态**
```bash
docker-compose ps
```

#### 方式三：使用脚本

```bash
# 开发环境
./scripts/start.sh dev

# 生产环境
./scripts/start.sh prod
```

### 配置说明

首次运行会自动生成默认配置文件 `assets/config/config.json`，您可以根据需要修改配置。

### 配置文件

首次运行会自动生成默认配置文件 `assets/config/config.json`：

```json
{
  "http_run_mode": "debug",
  "http_port": 8080,
  "db_host": "127.0.0.1:3306",
  "db_user": "root",
  "db_password": "newe123",
  "db_name": "newe",
  "redis_host": "127.0.0.1:6379"
}
```

## 📖 API文档

### 健康检查
- `GET /health` - 应用健康状态
- `GET /ready` - 服务就绪状态
- `GET /metrics` - 应用指标监控

### 认证接口
- `POST /admin/login` - 用户登录
- `GET /admin/logout` - 用户登出
- `GET /admin/verifysole` - Token验证

### 系统管理
- 菜单管理
- 角色管理
- 用户管理
- 字典管理
- 部门管理
- API管理

## 🏗️ 项目结构

```
newe/
├── assets/                 # 静态资源和配置文件
├── common/                # 公共模块
│   ├── basebank/          # 基础银行数据
│   ├── config/            # 配置管理
│   ├── db/                # 数据库模块
│   ├── redis/             # Redis缓存
│   ├── utils/             # 工具函数
│   ├── worklog/           # 工作日志
│   └── ws/                # WebSocket
├── model/                 # 数据模型
├── router/                # 路由层
│   ├── app/               # 应用路由
│   ├── middle/            # 中间件
│   ├── route/             # 路由注册
│   └── v1/                # API版本1
├── scripts/               # 脚本文件
└── examples/              # 示例代码
```

## 🔧 配置说明

### 数据库配置
```json
{
  "db_type": "mysql",
  "db_host": "127.0.0.1:3306",
  "db_user": "root",
  "db_password": "your_password",
  "db_name": "newe",
  "db_table_prefix": ""
}
```

### Redis配置
```json
{
  "redis_host": "127.0.0.1:6379",
  "redis_password": "",
  "redis_max_idle": 10,
  "redis_max_active": 100,
  "redis_idle_timeout": 300
}
```

### 文件上传配置
```json
{
  "img_prefix_url": "/images",
  "img_save_path": "upload/images",
  "img_max_size": 2097152,
  "img_allow_exts": "jpg,jpeg,png,gif",
  "file_prefix_url": "/files",
  "file_save_path": "upload/files",
  "file_max_size": 5242880,
  "file_allow_exts": "pdf,doc,docx,xls,xlsx,txt"
}
```

## 🛠️ 开发指南

### 添加新功能

1. **创建数据模型** - 在 `model/` 目录下定义数据结构
2. **添加路由** - 在 `router/v1/` 目录下创建API路由
3. **实现业务逻辑** - 在相应的路由文件中实现处理函数
4. **注册路由** - 在 `router/route/route.go` 中注册新路由

### 自定义中间件

在 `router/middle/` 目录下创建新的中间件文件：

```go
package middle

import "github.com/gin-gonic/gin"

func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}
```

## 📊 性能优化

- 数据库连接池配置
- Redis连接池优化
- GORM性能调优
- Gin中间件优化
- 静态资源缓存

## 🐛 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查MySQL服务是否启动
   - 确认数据库配置正确

2. **Redis连接失败**
   - 检查Redis服务是否启动
   - 确认Redis配置正确

3. **端口冲突**
   - 修改配置文件中的端口号

### 日志查看

应用日志位于 `assets/runtime/` 目录：
- `error/` - 错误日志
- `trace/` - 跟踪日志
- `sql/` - SQL查询日志

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📞 支持

如有问题请提交Issue或联系维护者。

---

⭐ 如果这个项目对你有帮助，请给它一个Star！