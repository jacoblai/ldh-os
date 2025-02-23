package service

import (
	"time"
)

// ServiceState 表示服务的状态
type ServiceState string

const (
	StateUnknown  ServiceState = "unknown"
	StateStarting ServiceState = "starting"
	StateRunning  ServiceState = "running"
	StateStopping ServiceState = "stopping"
	StateStopped  ServiceState = "stopped"
	StateFailed   ServiceState = "failed"
)

// ServiceType 表示服务的类型
type ServiceType string

const (
	TypeDaemon   ServiceType = "daemon"   // 持续运行的服务
	TypeOneshot  ServiceType = "oneshot"  // 一次性执行的服务
	TypePeriodic ServiceType = "periodic" // 周期性执行的服务
)

// ServiceConfig 定义服务的配置结构
type ServiceConfig struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Type         ServiceType       `yaml:"type"`
	ExecPath     string            `yaml:"exec"`
	Args         []string          `yaml:"args,omitempty"`
	Dependencies []string          `yaml:"dependencies,omitempty"`
	Environment  map[string]string `yaml:"environment,omitempty"`
	Restart      string            `yaml:"restart"`
	MCPConfig    MCPConfig         `yaml:"mcp"`
}

// MCPConfig 定义 MCP 相关配置
type MCPConfig struct {
	Functions   []string `yaml:"functions"`
	Permissions []string `yaml:"permissions"`
}

// ServiceStatus 定义服务的运行时状态
type ServiceStatus struct {
	State        ServiceState
	Pid          int
	StartTime    time.Time
	RestartCount int
	LastError    error
}

// EventType 定义事件类型
type EventType string

const (
	EventStarting EventType = "starting"
	EventStarted  EventType = "started"
	EventStopping EventType = "stopping"
	EventStopped  EventType = "stopped"
	EventFailed   EventType = "failed"
	EventRestart  EventType = "restart"
)
