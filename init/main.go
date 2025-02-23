package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
)

type InitSystem struct {
    state    string
    services map[string]bool
    signals  chan os.Signal
}

func NewInitSystem() *InitSystem {
    return &InitSystem{
        state:    "booting",
        services: make(map[string]bool),
        signals:  make(chan os.Signal, 1),
    }
}

func (i *InitSystem) mountEssentialFS() error {
    log.Println("Mounting essential filesystems...")
    
    // 挂载 proc 文件系统
    if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
        log.Printf("Warning: Failed to mount proc: %v", err)
    }
    
    // 挂载 sysfs
    if err := syscall.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
        log.Printf("Warning: Failed to mount sysfs: %v", err)
    }
    
    // 挂载 devtmpfs
    if err := syscall.Mount("devtmpfs", "/dev", "devtmpfs", syscall.MS_NOSUID, "mode=755"); err != nil {
        log.Printf("Warning: Failed to mount devtmpfs: %v", err)
    }
    
    return nil
}

func (i *InitSystem) initializeDevices() error {
    log.Println("Initializing devices...")
    // 基础设备初始化逻辑将在这里实现
    return nil
}

func (i *InitSystem) handleSignals() {
    for sig := range i.signals {
        switch sig {
        case syscall.SIGTERM:
            log.Println("Received SIGTERM, initiating shutdown...")
            i.shutdown()
        case syscall.SIGINT:
            log.Println("Received SIGINT, initiating shutdown...")
            i.shutdown()
        default:
            log.Printf("Received signal: %v", sig)
        }
    }
}

func (i *InitSystem) shutdown() {
    log.Println("Shutting down...")
    // 实现清理和关闭逻辑
    os.Exit(0)
}

func main() {
    if os.Getpid() != 1 {
        log.Printf("Warning: Not running as PID 1 (current PID: %d)", os.Getpid())
    }
    
    log.Println("LDH-OS Init starting...")
    
    init := NewInitSystem()
    
    // 设置信号处理
    signal.Notify(init.signals,
        syscall.SIGTERM,
        syscall.SIGINT,
        syscall.SIGHUP,
        syscall.SIGQUIT)
    
    if err := init.mountEssentialFS(); err != nil {
        log.Fatal("Failed to mount filesystems:", err)
    }
    
    if err := init.initializeDevices(); err != nil {
        log.Fatal("Failed to initialize devices:", err)
    }
    
    log.Println("Init system ready")
    
    // 处理系统信号
    init.handleSignals()
}
