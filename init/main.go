package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"ldh-os/init/service"

	"golang.org/x/sys/unix"
)

type InitSystem struct {
	state          string
	serviceManager *service.ServiceManager
	signals        chan os.Signal
}

func NewInitSystem() *InitSystem {
	return &InitSystem{
		state:          "booting",
		serviceManager: service.NewServiceManager(),
		signals:        make(chan os.Signal, 1),
	}
}

func (i *InitSystem) mountEssentialFS() error {
	log.Println("Mounting essential filesystems...")

	// 挂载 proc 文件系统
	if err := unix.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		log.Printf("Warning: Failed to mount proc: %v", err)
	}

	// 挂载 sysfs
	if err := unix.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
		log.Printf("Warning: Failed to mount sysfs: %v", err)
	}

	// 挂载 devtmpfs
	if err := unix.Mount("devtmpfs", "/dev", "devtmpfs", unix.MS_NOSUID, "mode=755"); err != nil {
		log.Printf("Warning: Failed to mount devtmpfs: %v", err)
	}

	return nil
}

func (i *InitSystem) initializeDevices() error {
	log.Println("Initializing devices...")
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
	log.Println("Shutting down all services...")
	if err := i.serviceManager.StopAll(); err != nil {
		log.Printf("Error stopping services: %v", err)
	}

	log.Println("Unmounting filesystems...")
	// 按照相反的顺序卸载文件系统
	unix.Unmount("/dev", 0)
	unix.Unmount("/sys", 0)
	unix.Unmount("/proc", 0)

	os.Exit(0)
}

func (i *InitSystem) loadServices() error {
	// 获取配置文件路径
	configPath := "/etc/ldh-os/services.yaml"
	if os.Getenv("LDH_SERVICES_CONFIG") != "" {
		configPath = os.Getenv("LDH_SERVICES_CONFIG")
	}

	// 确保配置目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := i.createDefaultConfig(configPath); err != nil {
			return err
		}
	}

	// 加载服务配置
	if err := i.serviceManager.LoadServices(configPath); err != nil {
		return err
	}

	return nil
}

func (i *InitSystem) createDefaultConfig(configPath string) error {
	defaultConfig := `
# LDH-OS 默认服务配置
syslog:
  description: "System logging service"
  type: "daemon"
  exec: "/usr/sbin/syslogd"
  restart: "always"
  mcp:
    functions: ["start", "stop", "restart", "status"]
    permissions: ["read", "write"]

cron:
  description: "Cron daemon"
  type: "daemon"
  exec: "/usr/sbin/crond"
  args: ["-n"]
  dependencies: ["syslog"]
  restart: "always"
  mcp:
    functions: ["start", "stop", "restart", "status"]
    permissions: ["read", "write"]
`
	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
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

	// 加载并启动服务
	if err := init.loadServices(); err != nil {
		log.Printf("Warning: Failed to load services: %v", err)
	} else {
		if err := init.serviceManager.StartAll(); err != nil {
			log.Printf("Warning: Failed to start all services: %v", err)
		}
	}

	log.Println("Init system ready")

	// 处理系统信号
	init.handleSignals()
}
