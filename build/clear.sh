#!/bin/bash

# 设置项目根目录
PROJECT_ROOT="$(pwd)"
KERNEL_SOURCE="$PROJECT_ROOT/kernel"
BUILD_DIR="$PROJECT_ROOT/build/output"

# 清理内核编译
echo "Cleaning kernel build files..."
if [ -d "$KERNEL_SOURCE" ]; then
    cd "$KERNEL_SOURCE" || exit 1
    if [ -f "Makefile" ]; then
        make clean    # 清理大部分编译文件
        make mrproper # 完全清理，包括.config文件
    else
        echo "Warning: Kernel source not found in $KERNEL_SOURCE"
    fi
    cd "$PROJECT_ROOT" || exit 1
else
    echo "Warning: Kernel directory not found at $KERNEL_SOURCE"
fi

# 清理构建输出目录
echo "Cleaning build output directory..."
if [ -d "$BUILD_DIR" ]; then
    rm -rf "${BUILD_DIR:?}"/*
    echo "Build output directory cleaned"
else
    echo "Build output directory not found at $BUILD_DIR"
fi