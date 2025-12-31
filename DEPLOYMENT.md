# 信臣健康互联网医院 - 生产环境部署指南

本指南将帮助您将信臣健康互联网医院系统部署到生产服务器。

## 目录
- [服务器要求](#服务器要求)
- [部署前准备](#部署前准备)
- [方式一：手动部署](#方式一手动部署)
- [方式二：Docker部署](#方式二docker部署)
- [Nginx反向代理配置](#nginx反向代理配置)
- [进程管理](#进程管理)
- [SSL/HTTPS配置](#sslhttps配置)
- [监控和日志](#监控和日志)

---

## 服务器要求

### 最低配置
- **CPU**: 2核
- **内存**: 2GB
- **硬盘**: 10GB
- **操作系统**: Linux (Ubuntu 20.04+, CentOS 7+, Debian 10+)
- **网络**: 公网IP，开放80和443端口

### 推荐配置
- **CPU**: 4核
- **内存**: 4GB
- **硬盘**: 20GB SSD
- **带宽**: 5Mbps+

---

## 部署前准备

### 1. 在本地构建项目

```bash
# 进入项目目录
cd /home/oliver/projects/medseek

# 构建前端
cd frontend
npm run build
cd ..

# 构建后端（针对Linux服务器）
GOOS=linux GOARCH=amd64 go build -o medseek cmd/server/main.go
```

### 2. 准备部署文件

创建部署包目录：
```bash
mkdir -p ~/medseek-deploy
cp medseek ~/medseek-deploy/
cp -r frontend/dist ~/medseek-deploy/
cp .env.example ~/medseek-deploy/.env
```

### 3. 打包文件

```bash
cd ~
tar -czf medseek-deploy.tar.gz medseek-deploy/
```

现在您有了 `medseek-deploy.tar.gz` 文件，可以传输到服务器。

---

## 方式一：手动部署

### 1. 传输文件到服务器

```bash
# 使用scp传输（替换your-server-ip和用户名）
scp medseek-deploy.tar.gz user@your-server-ip:/home/user/

# 或使用rsync
rsync -avz medseek-deploy.tar.gz user@your-server-ip:/home/user/
```

### 2. 在服务器上解压并配置

```bash
# SSH连接到服务器
ssh user@your-server-ip

# 解压文件
cd /home/user
tar -xzf medseek-deploy.tar.gz
cd medseek-deploy

# 配置环境变量
nano .env
```

编辑 `.env` 文件：
```bash
# DeepSeek API密钥（必填）
DEEPSEEK_API_KEY=your_deepseek_api_key_here

# 服务器端口（默认8080）
PORT=8080
```

### 3. 设置权限

```bash
chmod +x medseek
```

### 4. 测试运行

```bash
# 前台运行测试
./medseek

# 在浏览器访问 http://your-server-ip:8080
# 确认服务正常后按Ctrl+C停止
```

### 5. 使用systemd管理服务（推荐）

创建systemd服务文件：
```bash
sudo nano /etc/systemd/system/medseek.service
```

添加以下内容：
```ini
[Unit]
Description=信臣健康互联网医院服务
After=network.target

[Service]
Type=simple
User=user
WorkingDirectory=/home/user/medseek-deploy
Environment="PATH=/usr/local/bin:/usr/bin:/bin"
EnvironmentFile=/home/user/medseek-deploy/.env
ExecStart=/home/user/medseek-deploy/medseek
Restart=always
RestartSec=10

# 日志
StandardOutput=append:/var/log/medseek/access.log
StandardError=append:/var/log/medseek/error.log

# 安全设置
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

创建日志目录：
```bash
sudo mkdir -p /var/log/medseek
sudo chown user:user /var/log/medseek
```

启动服务：
```bash
# 重新加载systemd配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start medseek

# 查看状态
sudo systemctl status medseek

# 设置开机自启动
sudo systemctl enable medseek

# 查看日志
sudo journalctl -u medseek -f
```

### 6. 常用管理命令

```bash
# 启动服务
sudo systemctl start medseek

# 停止服务
sudo systemctl stop medseek

# 重启服务
sudo systemctl restart medseek

# 查看状态
sudo systemctl status medseek

# 查看实时日志
sudo journalctl -u medseek -f

# 查看最近100行日志
sudo journalctl -u medseek -n 100
```

---

## 方式二：Docker部署

### 1. 创建Dockerfile

在项目根目录创建 `Dockerfile`:
```bash
cd /home/oliver/projects/medseek
nano Dockerfile
```

```dockerfile
# 多阶段构建
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o medseek cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-builder /app/medseek .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./medseek"]
```

### 2. 创建docker-compose.yml

```bash
nano docker-compose.yml
```

```yaml
version: '3.8'

services:
  medseek:
    build: .
    container_name: medseek
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - DEEPSEEK_API_KEY=${DEEPSEEK_API_KEY}
      - PORT=8080
    env_file:
      - .env
    volumes:
      - ./logs:/var/log/medseek
    networks:
      - medseek-network

networks:
  medseek-network:
    driver: bridge
```

### 3. 构建并运行

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重启服务
docker-compose restart
```

---

## Nginx反向代理配置

### 1. 安装Nginx

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx

# CentOS/RHEL
sudo yum install nginx
```

### 2. 配置Nginx

创建配置文件：
```bash
sudo nano /etc/nginx/sites-available/medseek
```

添加以下内容：
```nginx
# HTTP配置
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    # 访问日志
    access_log /var/log/nginx/medseek-access.log;
    error_log /var/log/nginx/medseek-error.log;

    # 客户端最大上传大小
    client_max_body_size 10M;

    # WebSocket支持
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket超时设置
        proxy_read_timeout 86400;
        proxy_send_timeout 86400;
    }

    # API请求
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 静态文件
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. 启用配置

```bash
# 创建符号链接
sudo ln -s /etc/nginx/sites-available/medseek /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重载Nginx
sudo systemctl reload nginx
```

---

## SSL/HTTPS配置

### 使用Let's Encrypt免费SSL证书

```bash
# 安装certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书（替换your-domain.com）
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# 自动续期测试
sudo certbot renew --dry-run
```

Certbot会自动修改Nginx配置，添加HTTPS支持。

### 手动HTTPS配置

如果使用自己的证书，修改Nginx配置：
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 其他配置同上...
}

# HTTP重定向到HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

---

## 监控和日志

### 1. 查看应用日志

```bash
# systemd日志
sudo journalctl -u medseek -f

# 应用日志文件
tail -f /var/log/medseek/access.log
tail -f /var/log/medseek/error.log
```

### 2. 查看Nginx日志

```bash
tail -f /var/log/nginx/medseek-access.log
tail -f /var/log/nginx/medseek-error.log
```

### 3. 系统资源监控

```bash
# 查看进程资源使用
top
htop

# 查看内存使用
free -h

# 查看磁盘使用
df -h

# 查看网络连接
netstat -tuln | grep 8080
```

### 4. 设置日志轮转

创建日志轮转配置：
```bash
sudo nano /etc/logrotate.d/medseek
```

```
/var/log/medseek/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    missingok
    create 0640 user user
}
```

---

## 更新部署

当需要更新应用时：

```bash
# 1. 在本地重新构建
cd /home/oliver/projects/medseek
git pull  # 如果使用git
GOOS=linux GOARCH=amd64 go build -o medseek cmd/server/main.go
cd frontend && npm run build && cd ..

# 2. 传输到服务器
scp medseek user@your-server-ip:/home/user/medseek-deploy/
scp -r frontend/dist/* user@your-server-ip:/home/user/medseek-deploy/dist/

# 3. 在服务器上重启服务
ssh user@your-server-ip
sudo systemctl restart medseek
```

---

## 故障排查

### 服务无法启动
```bash
# 查看详细错误日志
sudo journalctl -u medseek -n 50 --no-pager

# 检查端口占用
sudo netstat -tuln | grep 8080

# 检查环境变量
sudo systemctl show medseek | grep Environment
```

### WebSocket连接失败
```bash
# 检查Nginx配置
sudo nginx -t

# 查看Nginx错误日志
sudo tail -f /var/log/nginx/error.log

# 测试WebSocket连接
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  http://localhost:8080/ws?session_id=test&user_id=test
```

### API响应慢
```bash
# 检查服务器资源
top
free -h
df -h

# 检查网络延迟
ping api.deepseek.com
```

---

## 安全建议

1. **防火墙配置**
```bash
# 只开放必要端口
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

2. **定期更新系统**
```bash
sudo apt update && sudo apt upgrade -y
```

3. **备份数据**
```bash
# 定期备份配置和日志
tar -czf backup-$(date +%Y%m%d).tar.gz /home/user/medseek-deploy
```

4. **限制API密钥权限**
   - 确保`.env`文件权限为600
   - 不要将API密钥提交到版本控制

5. **监控异常访问**
```bash
# 查看访问日志中的异常
sudo tail -f /var/log/nginx/medseek-access.log | grep -v "200\|301\|304"
```

---

## 性能优化

1. **Nginx缓存配置**
```nginx
# 在http块中添加
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=medseek_cache:10m max_size=100m inactive=60m;
```

2. **启用Gzip压缩**
```nginx
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/javascript application/json;
```

3. **Go服务优化**
   - 使用`GOMAXPROCS`设置合适的CPU核心数
   - 启用生产模式编译：`go build -ldflags="-s -w"`

---

## 联系支持

如有部署问题，请检查：
1. 服务日志：`sudo journalctl -u medseek -f`
2. Nginx日志：`/var/log/nginx/medseek-error.log`
3. 系统资源：`top`, `free -h`, `df -h`

---

**部署成功后，您的信臣健康互联网医院将在以下地址访问：**
- HTTP: `http://your-domain.com`
- HTTPS: `https://your-domain.com`

祝部署顺利！🎉
