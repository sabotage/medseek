#!/bin/bash

# 信臣健康互联网医院 - 快速部署打包脚本
# 用途：在本地构建并打包所有需要部署到服务器的文件

set -e  # 遇到错误立即退出

echo "======================================"
echo "信臣健康互联网医院 - 部署打包工具"
echo "======================================"
echo ""

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

# 部署包目录
DEPLOY_DIR="$HOME/medseek-deploy-$(date +%Y%m%d_%H%M%S)"
ARCHIVE_NAME="medseek-deploy-$(date +%Y%m%d_%H%M%S).tar.gz"

echo "📦 步骤 1/5: 清理旧的构建文件..."
rm -rf frontend/dist
rm -f medseek

echo "✅ 清理完成"
echo ""

echo "🔨 步骤 2/5: 构建前端..."
cd frontend
npm run build
if [ $? -ne 0 ]; then
    echo "❌ 前端构建失败！"
    exit 1
fi
echo "✅ 前端构建完成"
cd ..
echo ""

echo "🔨 步骤 3/5: 构建后端（Linux x64）..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o medseek cmd/server/main.go
if [ $? -ne 0 ]; then
    echo "❌ 后端构建失败！"
    exit 1
fi
echo "✅ 后端构建完成"
echo ""

echo "📂 步骤 4/5: 准备部署文件..."
mkdir -p "$DEPLOY_DIR"

# 复制必要文件
cp medseek "$DEPLOY_DIR/"
cp -r frontend/dist "$DEPLOY_DIR/"
cp .env.example "$DEPLOY_DIR/.env"

# 创建README
cat > "$DEPLOY_DIR/README.txt" << 'EOF'
信臣健康互联网医院 - 部署包
===============================

部署步骤：

1. 解压此文件到服务器目录：
   tar -xzf medseek-deploy-*.tar.gz
   cd medseek-deploy-*

2. 编辑.env文件，配置DeepSeek API密钥：
   nano .env
   
   必填项：
   DEEPSEEK_API_KEY=你的API密钥
   PORT=8080

3. 设置执行权限：
   chmod +x medseek

4. 测试运行：
   ./medseek
   
   访问 http://服务器IP:8080 测试

5. 配置systemd服务（推荐生产环境）：
   参考完整部署文档：DEPLOYMENT.md

文件说明：
- medseek: 后端可执行文件
- dist/: 前端静态文件
- .env: 环境配置文件（需要编辑）

更多详细信息请参考项目中的 DEPLOYMENT.md 文档。
EOF

echo "✅ 部署文件准备完成"
echo ""

echo "📦 步骤 5/5: 打包部署文件..."
cd "$(dirname "$DEPLOY_DIR")"
tar -czf "$ARCHIVE_NAME" "$(basename "$DEPLOY_DIR")"
mv "$ARCHIVE_NAME" "$PROJECT_DIR/"
rm -rf "$DEPLOY_DIR"

cd "$PROJECT_DIR"
echo "✅ 打包完成"
echo ""

echo "======================================"
echo "✅ 部署包创建成功！"
echo "======================================"
echo ""
echo "📦 部署包位置: $PROJECT_DIR/$ARCHIVE_NAME"
echo "📊 文件大小: $(du -h "$ARCHIVE_NAME" | cut -f1)"
echo ""
echo "下一步："
echo "1. 将部署包传输到服务器："
echo "   scp $ARCHIVE_NAME user@your-server:/home/user/"
echo ""
echo "2. 在服务器上解压并配置："
echo "   ssh user@your-server"
echo "   tar -xzf $ARCHIVE_NAME"
echo "   cd $(basename "$ARCHIVE_NAME" .tar.gz)"
echo "   nano .env  # 配置API密钥"
echo "   chmod +x medseek"
echo "   ./medseek  # 测试运行"
echo ""
echo "详细部署步骤请参考: DEPLOYMENT.md"
echo "======================================"
