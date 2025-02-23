package service

import (
	"sync"
)

// StateManager 管理服务状态
type StateManager struct {
	states       map[string]ServiceState
	dependencies map[string][]string
	mu           sync.RWMutex
}

// NewStateManager 创建新的状态管理器
func NewStateManager() *StateManager {
	return &StateManager{
		states:       make(map[string]ServiceState),
		dependencies: make(map[string][]string),
	}
}

// UpdateState 更新服务状态
func (sm *StateManager) UpdateState(service string, state ServiceState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[service] = state
}

// GetState 获取服务状态
func (sm *StateManager) GetState(service string) ServiceState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if state, exists := sm.states[service]; exists {
		return state
	}
	return StateUnknown
}

// SetDependencies 设置服务依赖关系
func (sm *StateManager) SetDependencies(service string, deps []string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.dependencies[service] = deps
}

// GetDependencies 获取服务的依赖
func (sm *StateManager) GetDependencies(service string) []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if deps, exists := sm.dependencies[service]; exists {
		return deps
	}
	return nil
}

// CheckDependencies 检查服务的依赖是否都已启动
func (sm *StateManager) CheckDependencies(service string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	deps := sm.dependencies[service]
	for _, dep := range deps {
		if state, exists := sm.states[dep]; !exists || state != StateRunning {
			return false
		}
	}
	return true
}

// GetAllServices 获取所有服务及其状态
func (sm *StateManager) GetAllServices() map[string]ServiceState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make(map[string]ServiceState)
	for service, state := range sm.states {
		result[service] = state
	}
	return result
}

// RemoveService 从状态管理器中移除服务
func (sm *StateManager) RemoveService(service string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.states, service)
	delete(sm.dependencies, service)
}
