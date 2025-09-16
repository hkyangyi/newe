# Newe Makefile
# 使用方法: make [target]

.PHONY: help build run dev test clean docker docker-compose

# 默认目标
help:
	@echo "🚀 Newe 项目管理工具"
	@echo ""
	@echo "使用方法:"
	@echo "  make build     - 构建项目"
	@echo "  make run       - 运行项目"
	@echo "  make dev       - 开发模式运行"
	@echo "  make test      - 运行测试"
	@echo "  make clean     - 清理构建文件"
	@echo "  make docker    - 构建Docker镜像"
	@echo "  make compose   - 启动Docker Compose"
	@echo "  make lint      - 代码检查"
	@echo "  make fmt       - 代码格式化"

# 构建项目
build:
	@echo "🔨 构建项目..."
	go build -o newe-app .

# 运行项目
run: build
	@echo "🎯 运行项目..."
	./newe-app

# 开发模式
dev:
	@echo "🔧 开发模式启动..."
	go run examples/serve.go

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test ./... -v

# 清理
clean:
	@echo "🧹 清理构建文件..."
	rm -f newe-app
	rm -rf assets/runtime/*
	find . -name "*.log" -delete

# Docker构建
docker:
	@echo "🐳 构建Docker镜像..."
	docker build -t newe-app:latest .

# Docker Compose
compose:
	@echo "🚀 启动Docker Compose..."
	docker-compose up -d

# 代码检查
lint:
	@echo "📋 代码检查..."
	golangci-lint run ./...

# 代码格式化
fmt:
	@echo "🎨 代码格式化..."
	go fmt ./...

# 依赖检查
deps:
	@echo "📦 检查依赖..."
	go mod tidy
	go mod verify

# 生成文档
docs:
	@echo "📖 生成文档..."
	@echo "请使用swag或其他文档工具"

# 健康检查
health:
	@echo "🩺 健康检查..."
	curl http://localhost:8080/health

# 就绪检查
ready:
	@echo "✅ 就绪检查..."
	curl http://localhost:8080/ready

# 查看日志
logs:
	@echo "📝 查看日志..."
	tail -f assets/runtime/*/*.log

# 数据库迁移
migrate:
	@echo "🗄️  数据库迁移..."
	@echo "请实现数据库迁移脚本"

# 备份数据库
backup:
	@echo "💾 数据库备份..."
	@echo "请实现数据库备份脚本"