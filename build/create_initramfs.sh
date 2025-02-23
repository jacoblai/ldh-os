#!/bin/bash

# 设置变量
OUTPUT_DIR="build/output"
INITRAMFS_DIR="$OUTPUT_DIR/initramfs"
INIT_BINARY="$OUTPUT_DIR/init"

# 创建initramfs目录结构
mkdir -p "$INITRAMFS_DIR"/{bin,sbin,proc,sys,dev,etc,lib,lib64,mcp,llm}

# 复制init二进制
cp "$INIT_BINARY" "$INITRAMFS_DIR/init"

# 复制必要的共享库
ldd "$INIT_BINARY" | while read -r line; do
    if [[ $line =~ [^[:space:]]+$ ]]; then
        lib="${BASH_REMATCH[0]}"
        if [ -f "$lib" ]; then
            mkdir -p "$INITRAMFS_DIR$(dirname "$lib")"
            cp "$lib" "$INITRAMFS_DIR$lib"
        fi
    fi
done

# 创建initramfs
cd "$INITRAMFS_DIR"
find . | cpio -H newc -o | gzip > "../initrd.img"

echo "Initramfs created at $OUTPUT_DIR/initrd.img"
