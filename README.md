# LDH-OS: LLM-Driven Host Operating System

## 项目概述
LDH-OS是一个基于Linux内核的实验性操作系统，旨在构建一个以LLM（大语言模型）为核心的操作系统。该系统使用自定义的init系统，并计划集成LLM作为系统的核心决策组件。

## 当前状态
- [x] 基础内核配置与编译
- [x] 自定义Init系统（基础功能）
- [x] 基本文件系统挂载
- [x] 信号处理机制
- [ ] LLM集成
- [ ] MCP协议实现
- [ ] 系统服务管理

## 系统架构
```
Linux Kernel
|
Custom Init System
|
LLM Agent System (计划中)
|
Essential System Services
```

## 目录结构
```
ldh-os/
├── kernel/          # Linux内核源码
├── init/            # Init系统实现
├── llm/             # LLM相关实现（计划中）
│   ├── models/     # 模型文件目录
│   ├── data/       # 数据目录
│   └── config/     # 配置文件
├── mcp/             # MCP协议实现（计划中）
│   ├── functions/  # 功能模块
│   └── plugins/    # 插件
├── build/           # 构建脚本和工具
│   ├── build.sh
│   ├── test.sh
│   └── create_initramfs.sh
└── docs/            # 文档
```

## 环境要求
### 系统依赖
```bash
# 基础开发工具
- build-essential
- gcc/g++
- golang
- make
- cmake
- ninja-build
- qemu-system-x86

# 内核编译依赖
- bison
- flex
- libssl-dev
- libelf-dev
- bc
```

## 构建说明

### 1. 编译内核
```bash
cd kernel
make defconfig
make -j$(nproc)
```

### 2. 编译Init系统
```bash
cd init
go build -o ../build/output/init
```

### 3. 创建initramfs
```bash
./build/create_initramfs.sh
```

### 4. 运行测试
```bash
./build/test.sh
```

## 开发指南

### 设置开发环境
```bash
# 安装依赖
sudo apt update
sudo apt install -y build-essential gcc g++ gdb make cmake ninja-build \
    git golang python3 python3-pip qemu-system-x86 bison flex \
    libssl-dev libelf-dev bc
```

### 克隆项目
```bash
git clone https://github.com/jacoblai/ldh-os.git ldh-os
cd ldh-os
```

## linux内核下载
```bash
# 克隆Linux内核
git clone --depth 1 https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git kernel
```

## 执行构建和测试
```bash
# 构建init系统
cd init && go build -o ../build/output/init && cd ..

# 构建系统
./build/build.sh

# 创建initramfs
./build/create_initramfs.sh

# 运行测试
./build/test.sh
```

### 目前已实现的功能
1. 基础Init系统
    - 系统启动流程管理
    - 基本文件系统挂载
    - 信号处理机制
    - 简单的设备初始化框架

2. 构建系统
    - 自动化内核编译
    - Init系统构建
    - initramfs生成
    - QEMU测试环境

## 开发路线图

### Phase 1 - 当前阶段
- [x] 配置精简Linux内核
- [x] 实现基础Init系统
- [x] 构建最小文件系统
- [ ] 完善服务管理功能

### Phase 2 - 计划中
- [ ] 集成llama.cpp
- [ ] 实现基础MCP协议
- [ ] 构建系统代理框架

### Phase 3 - 未来计划
- [ ] 实现语音交互服务
- [ ] 集成图数据库
- [ ] 实现扩展系统功能

## 调试指南

### QEMU调试
系统通过QEMU启动，使用以下参数：
- 内存: 2GB
- CPU核心: 2
- 控制台: ttyS0

### 日志查看
Init系统的日志直接输出到控制台，可以通过QEMU串口查看。

## 贡献指南
1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证
[MIT License]

## 联系方式
[黎东海 - 邮箱：229292620@qq.com]
