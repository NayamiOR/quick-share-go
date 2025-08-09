#!/bin/bash

# 快速文件分享应用部署脚本

set -e

echo "🚀 开始部署快速文件分享应用..."

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装，请先安装Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 创建上传目录
echo "📁 创建上传目录..."
mkdir -p uploads

# 构建并启动容器
echo "🔨 构建Docker镜像..."
docker-compose build

echo "🚀 启动应用..."
docker-compose up -d

# 等待应用启动
echo "⏳ 等待应用启动..."
sleep 5

# 检查应用状态
if curl -f http://localhost:8080/ > /dev/null 2>&1; then
    echo "✅ 应用部署成功！"
    echo "🌐 访问地址: http://localhost:8080"
    echo "📊 查看日志: docker-compose logs -f"
    echo "🛑 停止应用: docker-compose down"
else
    echo "❌ 应用启动失败，请检查日志:"
    docker-compose logs
    exit 1
fi 