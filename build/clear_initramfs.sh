#!/bin/bash

# 设置变量
OUTPUT_DIR="build/output"
INITRAMFS_DIR="$OUTPUT_DIR/initramfs"
INIT_BINARY="$OUTPUT_DIR/init"
INITRD_IMG="$OUTPUT_DIR/initrd.img"

# 清理文件和目录
rm -rf "$INITRAMFS_DIR"
rm -f "$INIT_BINARY"
rm -f "$INITRD_IMG"

# 如果output目录为空，也删除它
if [ -d "$OUTPUT_DIR" ] && [ -z "$(ls -A $OUTPUT_DIR)" ]; then
    rmdir "$OUTPUT_DIR"
fi

echo "Cleanup completed."
