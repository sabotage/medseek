# Nginx子路径部署指南 (/medseek)

当您没有独立域名时，可以通过Nginx将应用部署在服务器的子路径下。

## 场景
- 服务器IP: `http://your-server-ip`
- 访问路径: `http://your-server-ip/medseek`
- 所有请求都路由到后端: `http://localhost:8080`

---

## 快速部署步骤

### 1️⃣ 安装Nginx
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install nginx

# CentOS/RHEL
sudo yum install nginx
```

### 2️⃣ 配置Nginx

复制配置文件：
```bash
sudo cp nginx.conf.template /etc/nginx/sites-available/medseek
```

编辑配置（可选，默认已配置好）：
```bash
sudo nano /etc/nginx/sites-available/medseek
```

### 3️⃣ 启用配置
```bash
# 创建符号链接
sudo ln -s /etc/nginx/sites-available/medseek /etc/nginx/sites-enabled/

# 禁用默认配置（如果需要）
sudo rm /etc/nginx/sites-enabled/default

# 测试配置
sudo nginx -t

# 重载Nginx
sudo systemctl reload nginx
```

### 4️⃣ 启动后端服务

```bash
# 进入部署目录
cd /home/user/medseek-deploy

# 后台运行（选项1：使用nohup）
nohup ./medseek > /tmp/medseek.log 2>&1 &

# 或后台运行（选项2：使用systemd - 推荐）
sudo systemctl start medseek
```

### 5️⃣ 访问应用

打开浏览器访问：
```
http://your-server-ip/medseek
```

---

## Nginx配置说明

当前配置将处理以下路由：

| 请求路径 | 代理目标 | 说明 |
|---------|---------|------|
| `/medseek/ws` | `http://localhost:8080/ws` | WebSocket连接 |
| `/medseek/api/*` | `http://localhost:8080/api/*` | API请求 |
| `/medseek/` | `http://localhost:8080/` | 静态文件和其他请求 |
| `/` | 重定向到 `/medseek/` | 根路径重定向 |

### 关键配置项

```nginx
# WebSocket代理需要这些头部
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";

# 长连接超时设置
proxy_read_timeout 86400;
proxy_send_timeout 86400;

# 前缀头部（可选，用于后端识别子路径）
proxy_set_header X-Forwarded-Prefix /medseek;
```

---

## 常见问题

### Q: 页面加载但样式/脚本不显示？

**A:** 这通常是因为前端资源路径问题。如果出现404，按以下步骤检查：

1. **检查Nginx日志**
```bash
sudo tail -f /var/log/nginx/medseek-error.log
```

2. **检查请求路径**
在浏览器开发者工具(F12) → Network标签中，查看资源请求的URL

3. **解决方案**
- 检查后端是否正确提供静态文件
- 确认`frontend/dist`文件已复制到服务器

### Q: WebSocket连接失败？

**A:** 检查以下配置：

```bash
# 确认Nginx配置正确
sudo nginx -t

# 检查Nginx日志
sudo tail -f /var/log/nginx/medseek-error.log

# 检查后端服务是否运行
ps aux | grep medseek
netstat -tuln | grep 8080
```

### Q: API请求返回404？

**A:** 检查代理路径是否正确：

```nginx
# 错误 ❌
location /medseek/api/ {
    proxy_pass http://localhost:8080/api;  # 末尾少了斜杠
}

# 正确 ✅
location /medseek/api/ {
    proxy_pass http://localhost:8080/api/;  # 末尾有斜杠
}
```

### Q: 如何查看访问日志？

```bash
# 实时查看
sudo tail -f /var/log/nginx/medseek-access.log

# 查看最后100行
sudo tail -n 100 /var/log/nginx/medseek-access.log

# 查看特定错误
sudo tail -f /var/log/nginx/medseek-error.log
```

---

## 监控和维护

### 查看Nginx状态
```bash
# 检查Nginx是否运行
sudo systemctl status nginx

# 重启Nginx
sudo systemctl restart nginx

# 重载配置（不中断服务）
sudo systemctl reload nginx
```

### 查看后端服务
```bash
# 检查服务是否运行
ps aux | grep medseek

# 查看进程占用的端口
sudo netstat -tuln | grep 8080

# 查看后端日志（如果使用systemd）
sudo journalctl -u medseek -f
```

### 性能优化

**1. 启用Gzip压缩**
```bash
sudo nano /etc/nginx/nginx.conf
```

在`http`块中添加：
```nginx
gzip on;
gzip_types text/plain text/css text/xml text/javascript application/javascript application/json;
gzip_min_length 1024;
gzip_vary on;
```

**2. 增加Nginx进程**
```nginx
# 在nginx.conf中修改
worker_processes auto;
worker_connections 2048;
```

---

## 防火墙配置

```bash
# 允许HTTP和HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 允许SSH（如果需要）
sudo ufw allow 22/tcp

# 启用防火墙
sudo ufw enable

# 查看规则
sudo ufw status
```

---

## 升级到HTTPS（可选）

### 使用Let's Encrypt获取免费证书（针对服务器IP）

```bash
# 安装certbot
sudo apt install certbot python3-certbot-nginx

# 如果只有IP地址，可以手动指定证书配置
# 或使用其他方式生成自签名证书
```

### 自签名证书（测试用）
```bash
# 生成自签名证书
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/ssl/private/medseek.key \
  -out /etc/ssl/certs/medseek.crt

# 编辑Nginx配置，启用SSL部分
sudo nano /etc/nginx/sites-available/medseek
# 取消注释并填写证书路径

# 测试和重载
sudo nginx -t
sudo systemctl reload nginx
```

---

## 故障排查清单

| 问题 | 排查命令 | 常见原因 |
|------|---------|---------|
| 无法访问应用 | `curl http://localhost:8080` | 后端未启动 |
| Nginx报错 | `sudo nginx -t` | 配置语法错误 |
| 页面加载但无样式 | `F12 → Network` | 资源路径错误 |
| WebSocket失败 | `sudo tail -f /var/log/nginx/medseek-error.log` | 代理配置不完整 |
| 响应慢 | `top`, `free -h` | 服务器资源不足 |

---

## 示例部署流程

```bash
# 1. 服务器上准备部署目录
mkdir -p ~/medseek-deploy
cd ~/medseek-deploy

# 2. 上传文件（从本地执行）
scp medseek user@your-server-ip:~/medseek-deploy/
scp -r dist user@your-server-ip:~/medseek-deploy/
scp .env user@your-server-ip:~/medseek-deploy/

# 3. 配置Nginx（服务器上执行）
sudo cp nginx.conf.template /etc/nginx/sites-available/medseek
sudo ln -s /etc/nginx/sites-available/medseek /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# 4. 启动后端服务
chmod +x medseek
./medseek &

# 5. 访问
# 打开浏览器：http://your-server-ip/medseek
```

---

## 需要帮助？

检查这些日志文件：
1. **Nginx访问日志**: `/var/log/nginx/medseek-access.log`
2. **Nginx错误日志**: `/var/log/nginx/medseek-error.log`
3. **后端日志**: `sudo journalctl -u medseek -f`
4. **系统日志**: `dmesg` 或 `journalctl -f`

祝部署顺利！🚀
