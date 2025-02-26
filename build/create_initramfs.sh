#!/bin/bash

# 设置错误时退出
set -e

# 获取当前工作目录的绝对路径
CURRENT_DIR="$(pwd)"

# 定义变量（使用绝对路径）
OUTPUT_DIR="$CURRENT_DIR/build/output"
INITRAMFS_DIR="$OUTPUT_DIR/initramfs"
INIT_BINARY="$OUTPUT_DIR/init"
INITRD_FILE="$OUTPUT_DIR/initrd.img"

echo "============================================"
echo "LDH-OS Initramfs创建工具"
echo "============================================"
echo "使用以下绝对路径:"
echo "当前目录: $CURRENT_DIR"
echo "输出目录: $OUTPUT_DIR"
echo "Initramfs目录: $INITRAMFS_DIR"
echo "Init二进制: $INIT_BINARY"
echo "Initrd文件: $INITRD_FILE"
echo "============================================"

# 检查init二进制是否存在
if [ ! -f "$INIT_BINARY" ]; then
    echo "错误: init二进制未找到: $INIT_BINARY"
    echo "请确保先编译init程序"
    echo "可以通过以下命令编译init:"
    echo "  cd init && go build -o ../build/output/init"
    exit 1
fi

echo "清理旧的initramfs目录..."
rm -rf "$INITRAMFS_DIR"

# 确保输出目录存在
echo "确保输出目录存在..."
mkdir -p "$OUTPUT_DIR"

echo "创建initramfs目录结构..."
mkdir -p "$INITRAMFS_DIR"/{bin,sbin,proc,sys,dev,etc,lib,lib64,mcp,llm,usr/bin,usr/sbin,usr/lib,usr/lib64,tmp,var/log,run}
mkdir -p "$INITRAMFS_DIR/etc/ldh-os"

# 设置目录权限
chmod 755 "$INITRAMFS_DIR"/{bin,sbin,proc,sys,dev,etc,lib,lib64,mcp,llm,usr,usr/bin,usr/sbin,usr/lib,usr/lib64,tmp,var,var/log,run}
chmod 1777 "$INITRAMFS_DIR/tmp"

echo "复制init二进制..."
cp "$INIT_BINARY" "$INITRAMFS_DIR/init"
chmod 755 "$INITRAMFS_DIR/init"

echo "创建默认配置文件..."
cat > "$INITRAMFS_DIR/etc/ldh-os/services.yaml" << 'EOF'
# LDH-OS 默认服务配置
syslog:
  description: "系统日志服务"
  type: "daemon"
  exec: "/bin/echo 'Syslog服务启动模拟'"
  restart: "always"

test:
  description: "测试服务"
  type: "oneshot"
  exec: "/bin/echo 'LDH-OS启动成功!'"
  dependencies: []
  restart: "never"
EOF

echo "创建/etc/inittab文件..."
cat > "$INITRAMFS_DIR/etc/inittab" << 'EOF'
# LDH-OS inittab
::sysinit:/init
::respawn:/bin/sh
::ctrlaltdel:/bin/echo "Ctrl-Alt-Del pressed"
EOF

echo "复制必要的共享库..."
ldd "$INIT_BINARY" 2>/dev/null | grep -v "statically linked" | while read -r line; do
    if [[ $line =~ [^[:space:]]+$ ]]; then
        lib="${BASH_REMATCH[0]}"
        if [ -f "$lib" ]; then
            mkdir -p "$INITRAMFS_DIR$(dirname "$lib")"
            cp -L "$lib" "$INITRAMFS_DIR$lib"
            echo "  已复制: $lib"
        fi
    fi
done

# 复制必要的设备节点(基本的设备节点，实际系统会自动创建)
echo "创建基本设备节点..."
# /dev/null
mknod -m 666 "$INITRAMFS_DIR/dev/null" c 1 3
# /dev/console
mknod -m 600 "$INITRAMFS_DIR/dev/console" c 5 1
# /dev/tty
mknod -m 666 "$INITRAMFS_DIR/dev/tty" c 5 0
# /dev/ttyS0
mknod -m 666 "$INITRAMFS_DIR/dev/ttyS0" c 4 64
# /dev/random
mknod -m 444 "$INITRAMFS_DIR/dev/random" c 1 8
# /dev/urandom
mknod -m 444 "$INITRAMFS_DIR/dev/urandom" c 1 9

# 添加必要的Shell脚本和工具
echo "添加基本Shell脚本..."
cat > "$INITRAMFS_DIR/bin/sh" << 'EOF'
#!/bin/busybox sh
echo "LDH-OS Shell"
echo "这是一个简单的Shell环境"
echo "当前路径: $(pwd)"
echo "可用命令: echo, ls"
while true; do
    printf "# "
    read cmd
    case "$cmd" in
        exit) break ;;
        *) echo "$cmd" ;;
    esac
done
EOF
chmod 755 "$INITRAMFS_DIR/bin/sh"

# 创建一个简单的测试脚本
cat > "$INITRAMFS_DIR/bin/test_system" << 'EOF'
#!/bin/sh
echo "LDH-OS系统测试"
echo "----------------------------"
echo "系统已成功启动!"
echo "内核版本: $(uname -r)"
echo "系统架构: $(uname -m)"
echo "主机名: $(uname -n)"
echo "----------------------------"
EOF
chmod 755 "$INITRAMFS_DIR/bin/test_system"

echo "创建initramfs..."

# 使用临时文件夹的方式来创建initramfs
echo "第一步：进入initramfs目录"
cd "$INITRAMFS_DIR"

echo "第二步：创建cpio归档"
find . -print | cpio -o -H newc > "$OUTPUT_DIR/initrd.cpio"

echo "第三步：检查cpio文件"
if [ ! -f "$OUTPUT_DIR/initrd.cpio" ]; then
    echo "错误：cpio文件创建失败"
    exit 1
fi

echo "第四步：压缩cpio文件为gzip格式"
gzip -f "$OUTPUT_DIR/initrd.cpio"

echo "第五步：将压缩后的文件重命名为initrd.img"
mv "$OUTPUT_DIR/initrd.cpio.gz" "$INITRD_FILE"

# 检查文件是否创建成功
if [ -f "$INITRD_FILE" ]; then
    INITRD_SIZE=$(ls -lh "$INITRD_FILE" | awk '{print $5}')
    echo "Initramfs文件创建成功：$INITRD_FILE"
    echo "文件大小：$INITRD_SIZE"
else
    echo "错误：Initramfs文件创建失败"
    echo "当前目录：$(pwd)"
    echo "尝试列出输出目录内容："
    ls -la "$OUTPUT_DIR"
    exit 1
fi

# 返回到原始目录
cd "$CURRENT_DIR"

echo "============================================"
echo "Initramfs创建完成: $INITRD_FILE"
echo "文件大小: $INITRD_SIZE"
echo ""
echo "可以通过以下命令测试系统:"
echo "  ./build/test.sh"
echo "============================================"
