#!/bin/bash

# 信臣健康互联网医院 - 快速更新脚本
# 用于在服务器上更新应用

echo "======================================"
echo "信臣健康互联网医院 - 快速更新"
echo "======================================"
echo ""

# 检查是否提供了tar.gz文件
if [ -z "$1" ]; then
    # 如果没有提供文件，查找最新的
    LATEST_FILE=$(ls -t medseek-deploy-*.tar.gz 2>/dev/null | head -1)
    if [ -z "$LATEST_FILE" ]; then
        echo "❌ 错误：未找到部署包文件"
        echo "用法: $0 [medseek-deploy-*.tar.gz]"
        exit 1
    fi
    DEPLOY_FILE="$LATEST_FILE"
else
    DEPLOY_FILE="$1"
fi

if [ ! -f "$DEPLOY_FILE" ]; then
    echo "❌ 错误：文件不存在: $DEPLOY_FILE"
    exit 1
fi

echo "📦 使用部署包: $DEPLOY_FILE"
echo ""

# 解压新文件
echo "1️⃣ 解压部署包..."
tar -xzf "$DEPLOY_FILE"
NEW_DIR=$(tar -tzf "$DEPLOY_FILE" | head -1 | cut -f1 -d"/")

if [ -z "$NEW_DIR" ]; then
    echo "❌ 解压失败"
    exit 1
fi

cd "$NEW_DIR"
echo "✅ 已解压到: $NEW_DIR"
echo ""

# 修复目录结构
echo "2️⃣ 检查目录结构..."
if [ -d "dist" ] && [ ! -d "frontend/dist" ]; then
    mkdir -p frontend
    mv dist frontend/dist
    echo "✅ 目录结构已修复"
else
    echo "✅ 目录结构正常"
fi

# 验证文件
if [ ! -f "frontend/dist/index.html" ]; then
    echo "❌ 错误：未找到前端文件"
    exit 1
fi
echo ""

# 配置环境变量
echo "3️⃣ 配置环境变量..."
if [ -f ".env" ]; then
    echo "📝 请编辑 .env 文件配置 DEEPSEEK_API_KEY"
    echo "   nano .env"
    echo ""
    read -p "是否现在编辑? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        nano .env
    fi
else
    echo "❌ 未找到.env文件"
    exit 1
fi

# 停止旧服务
echo ""
echo "4️⃣ 停止旧服务..."
pkill medseek || echo "没有运行中的服务"
sleep 2
echo "✅ 已停止"
echo ""

# 启动新服务
echo "5️⃣ 启动服务..."
chmod +x medseek
nohup ./medseek > medseek.log 2>&1 &
MEDSEEK_PID=$!
echo "✅ 服务已启动 (PID: $MEDSEEK_PID)"
sleep 3
echo ""

# 测试服务
echo "6️⃣ 测试服务..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ 2>/dev/null)

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ 服务测试成功 (HTTP $HTTP_CODE)"
else
    echo "⚠️  警告：服务响应异常 (HTTP $HTTP_CODE)"
    echo "   查看日志: tail -f medseek.log"
fi

echo ""
echo "======================================"
echo "✅ 更新完成！"
echo "======================================"
echo ""
echo "📍 部署目录: $HOME/$NEW_DIR"
echo "📊 访问地址: http://$(hostname -I | awk '{print $1}'):8080"
echo "📝 查看日志: tail -f $HOME/$NEW_DIR/medseek.log"
echo "🛑 停止服务: pkill medseek"
echo ""
echo "💡 移动端访问优化已启用："
echo "   - 响应式布局"
echo "   - 触摸优化"
echo "   - 防止缩放"
echo "   - 输入框自动适配"
echo ""
echo "======================================"
