#!/bin/bash

# 一键更新部署脚本

# 配置部分，请根据实际情况修改
APP_NAME="image-upload"                         # 应用名称
TEAMGRAMIO="$GOPATH/src"                        # 新增：团队代码目录
APP_DIR="${TEAMGRAMIO}/order"                    # 修改：应用所在目录
GIT_REPO="https://github.com/your/repo.git"    # Git 仓库地址
SERVICE_NAME="image-upload.service"             # systemd 服务名称
BIN_PATH="/usr/local/bin/$APP_NAME"            # 二进制文件路径

echo "=============================="
echo "开始一键更新部署流程"
echo "=============================="

# 1. 进入应用目录
echo "进入应用目录: $APP_DIR"
cd "$APP_DIR" || { echo "无法进入目录: $APP_DIR"; exit 1; }

# 2. 拉取最新代码
echo "拉取最新代码..."
git pull origin main || { echo "拉取代码失败"; exit 1; }

# 3. 编译应用
echo "编译应用..."
go build -o $APP_NAME . || { echo "编译失败"; exit 1; }

# 4. 停止现有服务
echo "停止现有服务: $SERVICE_NAME"
sudo systemctl stop $SERVICE_NAME || { echo "停止服务失败"; exit 1; }

# 5. 备份旧的二进制文件（可选）
if [ -f "$BIN_PATH" ]; then
    echo "备份旧的二进制文件..."
    sudo mv "$BIN_PATH" "${BIN_PATH}.bak_$(date +%Y%m%d%H%M%S)" || { echo "备份失败"; exit 1; }
fi

# 6. 部署新的二进制文件
echo "部署新的二进制文件..."
sudo mv "$APP_DIR/$APP_NAME" "$BIN_PATH" || { echo "部署失败"; exit 1; }
sudo chmod +x "$BIN_PATH" || { echo "设置执行权限失败"; exit 1; }

# 7. 启动服务
echo "启动服务: $SERVICE_NAME"
sudo systemctl start $SERVICE_NAME || { echo "启动服务失败"; exit 1; }

# 8. 检查服务状态
echo "检查服务状态..."
sudo systemctl status $SERVICE_NAME

echo "=============================="
echo "一键更新部署完成"
echo "==============================" 