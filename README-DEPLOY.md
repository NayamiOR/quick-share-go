# 快速文件分享应用 - 容器化部署指南

## 项目简介

这是一个基于Go语言和Gin框架开发的快速文件分享应用，支持文件上传、分享链接生成、文件下载和管理功能。

## 系统要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少 512MB 可用内存
- 至少 1GB 可用磁盘空间

## 快速部署

### 方法一：使用部署脚本（推荐）

```bash
# 给脚本执行权限
chmod +x deploy.sh

# 运行部署脚本
./deploy.sh
```

### 方法二：手动部署

```bash
# 1. 构建并启动容器
docker-compose up -d

# 2. 查看应用状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f
```

## 访问应用

部署成功后，通过以下地址访问：

- **主页**: http://localhost:8080
- **上传页面**: http://localhost:8080/upload
- **管理员页面**: http://localhost:8080/admin

## 管理员功能

- 管理员密码: `admin123`
- 可以查看所有分享的文件
- 可以删除不需要的分享

## 常用命令

```bash
# 启动应用
docker-compose up -d

# 停止应用
docker-compose down

# 查看日志
docker-compose logs -f

# 重启应用
docker-compose restart

# 重新构建镜像
docker-compose build --no-cache

# 查看容器状态
docker-compose ps
```

## 数据持久化

应用会将上传的文件保存在 `./uploads` 目录中，该目录已通过Docker卷映射到容器内，确保数据不会丢失。

## 自定义配置

### 修改端口

编辑 `docker-compose.yml` 文件中的端口映射：

```yaml
ports:
  - "你的端口:8080"
```

### 修改管理员密码

需要修改 `main.go` 文件中的 `adminPassword` 变量，然后重新构建镜像。

## 生产环境部署建议

1. **使用反向代理**: 建议使用Nginx作为反向代理
2. **配置HTTPS**: 使用Let's Encrypt等证书
3. **设置防火墙**: 只开放必要端口
4. **监控和日志**: 配置日志收集和监控
5. **备份策略**: 定期备份uploads目录

### 生产环境docker-compose.yml示例

```yaml
version: '3.8'

services:
  quick-share:
    build: .
    container_name: quick-share-app
    ports:
      - "127.0.0.1:8080:8080"  # 只允许本地访问
    volumes:
      - ./uploads:/app/uploads
      - ./logs:/app/logs
    environment:
      - GIN_MODE=release
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    container_name: quick-share-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - quick-share
    restart: unless-stopped
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 查看端口占用
   netstat -tulpn | grep 8080
   
   # 修改docker-compose.yml中的端口映射
   ```

2. **权限问题**
   ```bash
   # 确保uploads目录有正确权限
   chmod 755 uploads
   ```

3. **容器启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs
   
   # 检查磁盘空间
   df -h
   ```

## 更新应用

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 卸载应用

```bash
# 停止并删除容器
docker-compose down

# 删除镜像
docker rmi quick-share_quick-share

# 删除数据（谨慎操作）
rm -rf uploads/
``` 