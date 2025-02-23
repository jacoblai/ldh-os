package service

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

// ServiceManager 服务管理器
type ServiceManager struct {
	services     map[string]*Service
	stateManager *StateManager
	eventBus     *EventBus
	mcpHandler   *MCPHandler
	mu           sync.RWMutex
}

// NewServiceManager 创建新的服务管理器
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services:     make(map[string]*Service),
		stateManager: NewStateManager(),
		eventBus:     NewEventBus(),
		mcpHandler:   NewMCPHandler(),
	}
}

// LoadServices 从配置文件加载服务
func (sm *ServiceManager) LoadServices(configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var configs map[string]ServiceConfig
	if err := yaml.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	for name, config := range configs {
		config.Name = name
		if err := sm.RegisterService(config); err != nil {
			return fmt.Errorf("failed to register service %s: %v", name, err)
		}
	}

	return nil
}

// RegisterService 注册新服务
func (sm *ServiceManager) RegisterService(config ServiceConfig) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.services[config.Name]; exists {
		return fmt.Errorf("service %s already exists", config.Name)
	}

	service := NewService(config, sm.eventBus)
	sm.services[config.Name] = service
	sm.stateManager.SetDependencies(config.Name, config.Dependencies)

	// 注册 MCP 功能
	for _, funcName := range config.MCPConfig.Functions {
		sm.mcpHandler.RegisterFunction(config.Name, funcName, sm.createMCPHandler(service, funcName))
	}

	return nil
}

// StartService 启动服务
func (sm *ServiceManager) StartService(name string) error {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	// 检查依赖
	if !sm.stateManager.CheckDependencies(name) {
		// 获取依赖列表以提供更好的错误信息
		deps := sm.stateManager.GetDependencies(name)
		return fmt.Errorf("dependencies not satisfied for service %s: %v", name, deps)
	}

	// 启动服务
	if err := service.Start(); err != nil {
		return fmt.Errorf("failed to start service %s: %v", name, err)
	}

	// 更新状态管理器中的状态
	sm.stateManager.UpdateState(name, service.Status.State)

	return nil
}

// StopService 停止服务
func (sm *ServiceManager) StopService(name string) error {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	return service.Stop()
}

// RestartService 重启服务
func (sm *ServiceManager) RestartService(name string) error {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	return service.Restart()
}

// GetServiceStatus 获取服务状态
func (sm *ServiceManager) GetServiceStatus(name string) (ServiceStatus, error) {
	sm.mu.RLock()
	service, exists := sm.services[name]
	sm.mu.RUnlock()

	if !exists {
		return ServiceStatus{}, fmt.Errorf("service %s not found", name)
	}

	return service.GetStatus(), nil
}

// ListServices 列出所有服务
func (sm *ServiceManager) ListServices() map[string]ServiceStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make(map[string]ServiceStatus)
	for name, service := range sm.services {
		result[name] = service.GetStatus()
	}
	return result
}

// HandleMCPRequest 处理 MCP 请求
func (sm *ServiceManager) HandleMCPRequest(req *MCPRequest) *MCPResponse {
	return sm.mcpHandler.HandleRequest(req)
}

// createMCPHandler 创建 MCP 功能处理函数
func (sm *ServiceManager) createMCPHandler(service *Service, funcName string) MCPFunction {
	return func(params map[string]interface{}) (interface{}, error) {
		var err error
		switch funcName {
		case "start":
			err = service.Start()
			if err == nil {
				sm.stateManager.UpdateState(service.Config.Name, service.Status.State)
			}
			return nil, err
		case "stop":
			err = service.Stop()
			if err == nil {
				sm.stateManager.UpdateState(service.Config.Name, service.Status.State)
			}
			return nil, err
		case "restart":
			err = service.Restart()
			if err == nil {
				sm.stateManager.UpdateState(service.Config.Name, service.Status.State)
			}
			return nil, err
		case "status":
			return service.GetStatus(), nil
		default:
			return nil, fmt.Errorf("unknown function: %s", funcName)
		}
	}
}

// StartAll 启动所有服务
func (sm *ServiceManager) StartAll() error {
	sm.mu.RLock()
	services := make([]*Service, 0, len(sm.services))
	for _, service := range sm.services {
		services = append(services, service)
	}
	sm.mu.RUnlock()

	// 按依赖顺序启动服务
	for _, service := range services {
		if err := sm.StartService(service.Config.Name); err != nil {
			return fmt.Errorf("failed to start service %s: %v", service.Config.Name, err)
		}
	}

	return nil
}

// StopAll 停止所有服务
func (sm *ServiceManager) StopAll() error {
	sm.mu.RLock()
	services := make([]*Service, 0, len(sm.services))
	for _, service := range sm.services {
		services = append(services, service)
	}
	sm.mu.RUnlock()

	// 按依赖顺序的反序停止服务
	for i := len(services) - 1; i >= 0; i-- {
		if err := sm.StopService(services[i].Config.Name); err != nil {
			return fmt.Errorf("failed to stop service %s: %v", services[i].Config.Name, err)
		}
	}

	return nil
}
