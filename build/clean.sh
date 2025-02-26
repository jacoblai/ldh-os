#!/bin/bash

# 设置错误时退出
set -e

# 清理特定组件
clean_kernel() {
    echo "Cleaning kernel build files..."
    if [ -d "kernel" ]; then
        cd kernel && make clean
        cd ..
    fi
}

clean_third_party() {
    echo "Cleaning third-party libraries..."
    rm -rf third_party/llama.cpp/build
}

clean_initramfs() {
    echo "Cleaning initramfs files..."
    rm -f build/output/initramfs/etc/ldh-os/services.yaml
    rm -f build/output/initrd.img
}

clean_all() {
    echo "Cleaning all build artifacts..."
    clean_kernel
    clean_third_party
    clean_initramfs
    rm -rf build/output/*
}

# 显示帮助信息
show_help() {
    echo "Usage: $0 [component]"
    echo "Components:"
    echo "  kernel      - Clean kernel build files"
    echo "  third-party - Clean third-party libraries"
    echo "  initramfs   - Clean initramfs files"
    echo "  all         - Clean all build artifacts"
    echo "  help        - Show this help message"
}

# 根据参数执行清理
case "$1" in
    "kernel")
        clean_kernel
        ;;
    "third-party")
        clean_third_party
        ;;
    "initramfs")
        clean_initramfs
        ;;
    "all")
        clean_all
        ;;
    "help")
        show_help
        ;;
    *)
        echo "Error: Unknown component '$1'"
        show_help
        exit 1
        ;;
esac

echo "Clean completed." 