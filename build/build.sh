#!/bin/bash

# 构建配置
PROJECT_ROOT="$(pwd)"
KERNEL_SOURCE="$PROJECT_ROOT/kernel"
BUILD_DIR="$PROJECT_ROOT/build/output"
TOOLCHAIN_DIR="$PROJECT_ROOT/build/toolchain"

# 创建必要目录
mkdir -p "$BUILD_DIR"
mkdir -p "$TOOLCHAIN_DIR"

# 构建内核
build_kernel() {
    echo "Building kernel..."
    cd "$KERNEL_SOURCE" || exit 1
    make defconfig
    make -j$(nproc)
    cd "$PROJECT_ROOT" || exit 1
}

# 构建Init系统
build_init() {
    echo "Building init system..."
    cd "$PROJECT_ROOT/init" || exit 1
    go build -o "$BUILD_DIR/init"
    cd "$PROJECT_ROOT" || exit 1
}

# 主构建流程
main() {
    build_kernel
    build_init
}

main "$@"
