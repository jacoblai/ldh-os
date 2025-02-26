#!/bin/bash

# LDH-OS 构建系统
# 此脚本负责整个LDH-OS的构建流程，包括内核配置、编译和初始化系统

# 设置错误处理
set -e

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 项目根目录（脚本目录的上一级）
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# 目录配置（使用相对路径）
BUILD_DIR="${PROJECT_ROOT}/build/output"
THIRD_PARTY_DIR="${PROJECT_ROOT}/third_party"
KERNEL_DIR="${PROJECT_ROOT}/kernel"
CONFIG_FILE="${SCRIPT_DIR}/kernel_config.conf"

# Linux内核配置
KERNEL_REPO="https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git"

# Banner显示函数
show_banner() {
    echo "========================================"
    echo "  LDH-OS 构建系统"
    echo "========================================"
    echo "构建目录: ${BUILD_DIR}"
    echo "内核目录: ${KERNEL_DIR}"
    echo "配置文件: ${CONFIG_FILE}"
    echo "========================================"
}

# 创建必要目录
create_dirs() {
    mkdir -p "${BUILD_DIR}"
    mkdir -p "${THIRD_PARTY_DIR}"
    echo "创建目录完成"
}

# 下载或更新Linux内核
download_kernel() {
    echo "检查Linux内核..."
    if [ ! -d "${KERNEL_DIR}" ]; then
        echo "下载Linux内核..."
        git clone --depth 1 "${KERNEL_REPO}" "${KERNEL_DIR}"
    else
        echo "更新Linux内核..."
        cd "${KERNEL_DIR}" || exit 1
        git fetch --depth 1 origin
        git reset --hard origin/master
        cd "${PROJECT_ROOT}" || exit 1
    fi
    echo "内核源码准备完成"
}

# 加载内核配置
load_kernel_config() {
    echo "加载内核配置..."
    
    # 检查内核目录是否存在
    if [ ! -d "${KERNEL_DIR}" ]; then
        echo "错误: 内核目录不存在于 ${KERNEL_DIR}"
        echo "请先运行 './build.sh download-kernel' 下载内核源码"
        exit 1
    fi
    
    # 检查配置文件是否存在
    if [ ! -f "${CONFIG_FILE}" ]; then
        echo "错误: 内核配置文件不存在于 ${CONFIG_FILE}"
        exit 1
    fi
    
    echo "将配置文件复制到内核目录..."
    cp "${CONFIG_FILE}" "${KERNEL_DIR}/.config"
    
    # 进入内核目录
    cd "${KERNEL_DIR}" || exit 1
    
    # 确保配置文件格式正确并应用默认值到未指定选项
    echo "应用配置..."
    make olddefconfig
    
    cd "${PROJECT_ROOT}" || exit 1
    
    echo "内核配置已加载"
}

# 构建内核
build_kernel() {
    echo "构建内核..."
    
    # 检查内核目录是否存在
    if [ ! -d "${KERNEL_DIR}" ]; then
        echo "错误: 内核目录不存在"
        exit 1
    fi
    
    # 应用内核配置
    load_kernel_config
    
    # 编译内核
    cd "${KERNEL_DIR}" || exit 1
    echo "开始编译内核，这可能需要一些时间..."
    make -j$(nproc)
    
    # 返回项目根目录
    cd "${PROJECT_ROOT}" || exit 1
    
    echo "内核构建完成"
}

# 构建Init系统
build_init() {
    echo "构建Init系统..."
    
    # 检查init目录是否存在
    if [ ! -d "${PROJECT_ROOT}/init" ]; then
        echo "错误: Init系统目录不存在"
        exit 1
    fi
    
    cd "${PROJECT_ROOT}/init" || exit 1
    go build -o "${BUILD_DIR}/init"
    
    # 返回项目根目录
    cd "${PROJECT_ROOT}" || exit 1
    
    echo "Init系统构建完成"
}

# 清理构建目录
clean() {
    echo "清理构建目录..."
    rm -rf "${BUILD_DIR}"
    mkdir -p "${BUILD_DIR}"
    echo "构建目录已清理"
}

# 清理所有（包括下载的源码）
clean_all() {
    echo "清理所有内容..."
    clean
    echo "删除内核源码..."
    rm -rf "${KERNEL_DIR}"
    echo "所有内容已清理"
}

# 显示帮助信息
show_help() {
    echo "LDH-OS 构建系统使用方法:"
    echo "  ./build.sh [命令]"
    echo ""
    echo "可用命令:"
    echo "  all                 执行完整构建流程（默认）"
    echo "  download-kernel     仅下载/更新内核源码"
    echo "  load-config         仅加载内核配置到内核源码目录"
    echo "  build-kernel        仅构建内核"
    echo "  build-init          仅构建Init系统"
    echo "  clean               清理构建目录"
    echo "  clean-all           清理所有内容（包括源码）"
    echo "  help                显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./build.sh download-kernel  # 下载内核源码"
    echo "  ./build.sh load-config      # 加载内核配置"
    echo "  ./build.sh all              # 执行完整构建"
}

# 主构建流程
main() {
    show_banner
    
    # 检查参数
    case "$1" in
        download-kernel)
            create_dirs
            download_kernel
            ;;
        load-config)
            load_kernel_config
            ;;
        build-kernel)
            build_kernel
            ;;
        build-init)
            build_init
            ;;
        clean)
            clean
            ;;
        clean-all)
            clean_all
            ;;
        help)
            show_help
            ;;
        all|"")
            create_dirs
            download_kernel
            build_kernel
            build_init
            echo "LDH-OS 构建完成"
            ;;
        *)
            echo "错误: 未知命令 '$1'"
            show_help
            exit 1
            ;;
    esac
    
    exit 0
}

# 执行主函数
main "$@" 