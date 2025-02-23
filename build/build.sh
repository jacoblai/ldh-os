#!/bin/bash

# 构建配置
KERNEL_SOURCE="$(pwd)/kernel"
BUILD_DIR="$(pwd)/build/output"
TOOLCHAIN_DIR="$(pwd)/build/toolchain"

# 创建必要目录
mkdir -p "$BUILD_DIR"
mkdir -p "$TOOLCHAIN_DIR"

# 构建内核
build_kernel() {
    cd "$KERNEL_SOURCE"
    make defconfig
    make -j$(nproc)
}

# 构建Init系统
build_init() {
    cd "$(pwd)/init"
    go build -o "$BUILD_DIR/init"
}

# 主构建流程
main() {
    build_kernel
    build_init
}

main "$@"
