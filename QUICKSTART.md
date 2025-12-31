# 🚀 快速部署指南

这是最简化的部署步骤，让您快速将信臣健康互联网医院部署到服务器。

## 前置条件
- ✅ 一台Linux服务器（Ubuntu/CentOS/Debian）
- ✅ 服务器有公网IP
- ✅ DeepSeek API密钥

---

## 三步部署

### 步骤 1️⃣：在本地打包

```bash
cd /home/oliver/projects/medseek
./deploy.sh
```

这会生成一个 `medseek-deploy-日期时间.tar.gz` 文件。

### 步骤 2️⃣：传输到服务器

```bash
# 替换下面的user和服务器IP
scp medseek-deploy-*.tar.gz user@your-server-ip:/home/user/
```

### 步骤 3️⃣：在服务器上部署

```bash
# SSH登录服务器
ssh user@your-server-ip

# 解压
tar -xzf medseek-deploy-*.tar.gz
cd medseek-deploy-*

# 配置API密钥
nano .env
# 修改这一行：
# DEEPSEEK_API_KEY=your_actual_api_key_here

# 设置权限并运行
chmod +x medseek
./medseek
```

访问 `http://your-server-ip:8080` 测试！

---

## 生产环境配置（推荐）

### 1. 使用systemd管理服务

```bash
# 创建日志目录
sudo mkdir -p /var/log/medseek
sudo chown $USER:$USER /var/log/medseek

# 复制systemd服务文件（在服务器上）
sudo nano /etc/systemd/system/medseek.service
```

粘贴以下内容（**替换YOUR_USERNAME为实际用户名**）：
```ini
[Unit]
Description=信臣健康互联网医院服务
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME/medseek-deploy-XXXXX
EnvironmentFile=/home/YOUR_USERNAME/medseek-deploy-XXXXX/.env
ExecStart=/home/YOUR_USERNAME/medseek-deploy-XXXXX/medseek
Restart=always
RestartSec=10
StandardOutput=append:/var/log/medseek/access.log
StandardError=append:/var/log/medseek/error.log

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl start medseek
sudo systemctl enable medseek  # 开机自启动
sudo systemctl status medseek  # 查看状态
```

### 2. 使用Nginx反向代理（可选）

```bash
# 安装Nginx
sudo apt install nginx  # Ubuntu/Debian
# sudo yum install nginx  # CentOS

# 创建配置
sudo nano /etc/nginx/sites-available/medseek
```

粘贴以下内容（**替换YOUR_DOMAIN.com**）：
```nginx
server {
    listen 80;
    server_name YOUR_DOMAIN.com;
    
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

启用配置：
```bash
sudo ln -s /etc/nginx/sites-available/medseek /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 3. 配置HTTPS（推荐）

```bash
# 安装certbot
sudo apt install certbot python3-certbot-nginx

# 获取免费SSL证书
sudo certbot --nginx -d YOUR_DOMAIN.com
```

---

## 常用命令

```bash
# 查看服务状态
sudo systemctl status medseek

# 重启服务
sudo systemctl restart medseek

# 查看日志
sudo journalctl -u medseek -f

# 查看错误日志
tail -f /var/log/medseek/error.log
```

---

## 更新部署

当需要更新时：

```bash
# 1. 本地重新打包
cd /home/oliver/projects/medseek
./deploy.sh

# 2. 传输到服务器
scp medseek-deploy-*.tar.gz user@your-server-ip:/home/user/

# 3. 在服务器上
ssh user@your-server-ip
tar -xzf medseek-deploy-*.tar.gz
# 复制新的二进制文件和前端文件到旧目录
sudo systemctl restart medseek
```

---

## 故障排查

### 服务无法启动
```bash
# 查看详细日志
sudo journalctl -u medseek -n 100

# 检查端口是否被占用
sudo netstat -tuln | grep 8080
```

### 无法访问
```bash
# 检查防火墙
sudo ufw status
sudo ufw allow 8080/tcp

# 检查服务是否运行
ps aux | grep medseek
```

### API报错
```bash
# 检查环境变量
cat /home/user/medseek-deploy-*/.env

# 测试API密钥
curl https://api.deepseek.com/v1/models \
  -H "Authorization: Bearer YOUR_API_KEY"
```

---

## 🎉 完成！

现在您的信臣健康互联网医院已经部署成功！

- 访问地址：`http://your-server-ip:8080` 或 `https://your-domain.com`
- 查看完整文档：[DEPLOYMENT.md](DEPLOYMENT.md)

有问题？检查日志：`sudo journalctl -u medseek -f`
