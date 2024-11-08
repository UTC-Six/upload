set -e
set -x

TEAMGRAMIO="$GOPATH/src"

echo "切换到项目目录: ${TEAMGRAMIO}/order"
cd ${TEAMGRAMIO}/order

echo "重置本地更改"
git checkout .
echo "切换到 develop 分支"
git checkout develop
echo "拉取最新代码"
git pull

# 确保在启动服务前构建成功
echo "开始构建项目..."
make build || { echo "构建失败"; exit 1; }
echo "构建成功"

# 结束 order-api 进程
echo "尝试终止 order-api 进程"
pkill order-api || { echo "未找到 order-api 进程"; }

# 等待 order-api 进程完全终止
echo "等待 order-api 进程终止..."
TIMEOUT=30
COUNT=0
while pgrep order-api > /dev/null; do
    if [ "$COUNT" -ge "$TIMEOUT" ]; then
        echo "order-api 进程未能在 $TIMEOUT 秒内终止，强制退出脚本"
        exit 1
    fi
    echo "等待中... ($COUNT/$TIMEOUT)"
    sleep 1
    COUNT=$((COUNT + 1))
done
echo "order-api 进程已终止"

cd /opt/data/teamgram/bin

# 检查配置文件是否存在
if [ ! -f ../etc/order.yaml ]; then
    echo "配置文件 ../etc/order.yaml 不存在"
    exit 1
fi

echo "运行 order ..."
nohup ./order-api -f=../etc/order.yaml >> ../logs/order.log 2>&1 &
echo "order-api 已启动，日志输出到 ../logs/order.log"