#!/bin/bash

# 设置错误时退出
set -e

# 定义变量
OUTPUT_DIR="build/output"
INIT_BINARY="$OUTPUT_DIR/init"

echo "============================================"
echo "LDH-OS Init检查工具"
echo "============================================"

# 检查init二进制是否存在
if [ ! -f "$INIT_BINARY" ]; then
    echo "错误: init二进制未找到: $INIT_BINARY"
    echo "请确保先编译init程序"
    echo "可以通过以下命令编译init:"
    echo "  cd init && go build -o ../build/output/init"
    exit 1
fi

# 显示init二进制信息
echo "Init二进制文件存在: $INIT_BINARY"
echo "文件大小: $(ls -lh "$INIT_BINARY" | awk '{print $5}')"
echo "编译时间: $(stat -c %y "$INIT_BINARY")"

# 检查init的依赖库
echo ""
echo "Init依赖库:"
ldd "$INIT_BINARY" | while read -r line; do
    echo "  $line"
done

echo "============================================"
echo "检查通过! 可以继续创建initramfs和测试系统"
echo "下一步:"
echo "1. 执行 './build/create_initramfs.sh' 创建initramfs"
echo "2. 执行 './build/test.sh' 测试系统"
echo "============================================" 