#!/bin/bash

# Newe 项目启动脚本
# 使用方法: ./scripts/start.sh [dev|prod]

set -e

ENV=${1:-dev}
CONFIG_DIR="assets/config"
CONFIG_FILE="$CONFIG_DIR/config.json"

echo "🚀 启动 Newe 项目 ($ENV 环境)"

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "📝 配置文件不存在，正在创建默认配置..."
    mkdir -p "$CONFIG_DIR"
    
    # 运行Go程序生成默认配置
    go run -tags=init_config main.go || echo "⚠️  如果Go运行失败，请手动创建配置文件"
fi

# 根据环境设置变量
if [ "$ENV" = "prod" ]; then
    echo "🏗️  生产环境模式"
    export GIN_MODE=release
else
    echo "🔧 开发环境模式"
    export GIN_MODE=debug
fi

# 检查依赖
echo "📦 检查Go模块依赖..."
go mod tidy

# 构建项目
echo "🔨 构建项目..."
go build -o newe-app .

# 启动项目
echo "🎯 启动应用..."
./newe-app

# 如果启动失败，显示错误信息
if [ $? -ne 0 ]; then
    echo "❌ 启动失败"
    exit 1
fi