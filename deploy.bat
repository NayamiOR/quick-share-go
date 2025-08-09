@echo off
chcp 65001 >nul
echo 🚀 开始部署快速文件分享应用...

REM 检查Docker是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker未安装，请先安装Docker Desktop
    pause
    exit /b 1
)

REM 检查Docker Compose是否安装
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker Compose未安装，请先安装Docker Compose
    pause
    exit /b 1
)

REM 创建上传目录
echo 📁 创建上传目录...
if not exist uploads mkdir uploads

REM 构建并启动容器
echo 🔨 构建Docker镜像...
docker-compose build

echo 🚀 启动应用...
docker-compose up -d

REM 等待应用启动
echo ⏳ 等待应用启动...
timeout /t 5 /nobreak >nul

REM 检查应用状态
echo 🔍 检查应用状态...
curl -f http://localhost:8080/ >nul 2>&1
if errorlevel 1 (
    echo ❌ 应用启动失败，请检查日志:
    docker-compose logs
    pause
    exit /b 1
) else (
    echo ✅ 应用部署成功！
    echo 🌐 访问地址: http://localhost:8080
    echo 📊 查看日志: docker-compose logs -f
    echo 🛑 停止应用: docker-compose down
)

pause 