#!/bin/bash

# 设置错误时退出
set -e

# 定义变量
KERNEL="kernel/arch/x86_64/boot/bzImage"
INITRD="build/output/initrd.img"
QEMU_PID_FILE="/tmp/ldh-os-qemu.pid"

# 清理函数
cleanup() {
    if [ -f "$QEMU_PID_FILE" ]; then
        QEMU_PID=$(cat "$QEMU_PID_FILE")
        if ps -p "$QEMU_PID" > /dev/null; then
            echo "正在停止QEMU (PID: $QEMU_PID)..."
            kill -TERM "$QEMU_PID" 2>/dev/null || true
            sleep 1
            if ps -p "$QEMU_PID" > /dev/null; then
                echo "强制停止QEMU..."
                kill -9 "$QEMU_PID" 2>/dev/null || true
            fi
        fi
        rm -f "$QEMU_PID_FILE"
    fi
}

# 设置信号处理
trap cleanup EXIT INT TERM

# 检查必要文件
if [ ! -f "$KERNEL" ]; then
    echo "错误: 内核镜像未找到: $KERNEL"
    echo "请确保先执行了 './build.sh' 构建内核"
    exit 1
fi

if [ ! -f "$INITRD" ]; then
    echo "错误: initrd未找到: $INITRD"
    echo "请确保先执行了 './build/create_initramfs.sh' 创建initramfs"
    exit 1
fi

# 显示版本信息
echo "============================================"
echo "LDH-OS 测试环境"
echo "============================================"
echo "内核: $KERNEL"
echo "Initrd: $INITRD"
echo "QEMU PID文件: $QEMU_PID_FILE"
echo "当前时间: $(date)"
echo "============================================"

# 启动QEMU (使用更详细的控制台配置)
echo "正在启动QEMU..."
qemu-system-x86_64 \
    -kernel "$KERNEL" \
    -initrd "$INITRD" \
    -append "console=ttyS0 root=/dev/ram0 init=/init earlyprintk=serial,ttyS0,115200 debug loglevel=8" \
    -nographic \
    -no-reboot \
    -m 2G \
    -smp 2 \
    -device virtio-rng-pci \
    -pidfile "$QEMU_PID_FILE" &

QEMU_PID=$!
echo "QEMU启动，PID: $QEMU_PID"
echo "按Ctrl+C退出"
echo ""
echo "等待系统启动中..."
echo "如果长时间没有输出，可能是内核或init系统出现问题"
echo "============================================"

# 等待QEMU进程
wait $QEMU_PID
