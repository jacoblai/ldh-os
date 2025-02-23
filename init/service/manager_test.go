package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestServiceManager(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "services.yaml")

	config := `
test-service:
  description: "Test Service"
  type: "daemon"
  exec: "/bin/sleep"
  args: ["1000"]
  restart: "never"
  mcp:
    functions: ["start", "stop", "restart", "status"]
    permissions: ["read", "write"]

dependent-service:
  description: "Dependent Service"
  type: "daemon"
  exec: "/bin/sleep"
  args: ["1000"]
  dependencies: ["test-service"]
  restart: "never"
  mcp:
    functions: ["start", "stop", "restart", "status"]
    permissions: ["read", "write"]
`

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	// 创建服务管理器
	sm := NewServiceManager()

	// 测试加载服务
	t.Run("LoadServices", func(t *testing.T) {
		if err := sm.LoadServices(configPath); err != nil {
			t.Errorf("Failed to load services: %v", err)
		}

		services := sm.ListServices()
		if len(services) != 2 {
			t.Errorf("Expected 2 services, got %d", len(services))
		}
	})

	// 测试启动服务
	t.Run("StartService", func(t *testing.T) {
		if err := sm.StartService("test-service"); err != nil {
			t.Errorf("Failed to start test-service: %v", err)
		}

		// 等待服务启动
		time.Sleep(100 * time.Millisecond)

		status, err := sm.GetServiceStatus("test-service")
		if err != nil {
			t.Errorf("Failed to get service status: %v", err)
		}
		if status.State != StateRunning {
			t.Errorf("Expected service state %s, got %s", StateRunning, status.State)
		}
	})

	// 测试依赖检查
	t.Run("DependencyCheck", func(t *testing.T) {
		// 尝试启动依赖服务
		if err := sm.StartService("dependent-service"); err != nil {
			t.Errorf("Failed to start dependent-service: %v", err)
		}

		// 等待服务启动
		time.Sleep(100 * time.Millisecond)

		status, err := sm.GetServiceStatus("dependent-service")
		if err != nil {
			t.Errorf("Failed to get service status: %v", err)
		}
		if status.State != StateRunning {
			t.Errorf("Expected service state %s, got %s", StateRunning, status.State)
		}
	})

	// 测试停止服务
	t.Run("StopService", func(t *testing.T) {
		if err := sm.StopService("dependent-service"); err != nil {
			t.Errorf("Failed to stop dependent-service: %v", err)
		}
		if err := sm.StopService("test-service"); err != nil {
			t.Errorf("Failed to stop test-service: %v", err)
		}

		// 等待服务停止
		time.Sleep(100 * time.Millisecond)

		status, err := sm.GetServiceStatus("test-service")
		if err != nil {
			t.Errorf("Failed to get service status: %v", err)
		}
		if status.State != StateStopped {
			t.Errorf("Expected service state %s, got %s", StateStopped, status.State)
		}
	})

	// 测试 MCP 功能
	t.Run("MCPFunctions", func(t *testing.T) {
		req := &MCPRequest{
			Service:  "test-service",
			Function: "status",
			Params:   map[string]interface{}{},
		}

		resp := sm.HandleMCPRequest(req)
		if !resp.Success {
			t.Errorf("MCP request failed: %s", resp.Error)
		}
	})
}

func TestServiceManagerErrors(t *testing.T) {
	sm := NewServiceManager()

	// 测试启动不存在的服务
	t.Run("StartNonExistentService", func(t *testing.T) {
		err := sm.StartService("non-existent")
		if err == nil {
			t.Error("Expected error when starting non-existent service")
		}
	})

	// 测试注册重复服务
	t.Run("RegisterDuplicateService", func(t *testing.T) {
		config := ServiceConfig{
			Name:      "test",
			Type:      TypeDaemon,
			ExecPath:  "/bin/sleep",
			Args:      []string{"1000"},
			Restart:   "never",
			MCPConfig: MCPConfig{},
		}

		if err := sm.RegisterService(config); err != nil {
			t.Errorf("Failed to register first service: %v", err)
		}

		err := sm.RegisterService(config)
		if err == nil {
			t.Error("Expected error when registering duplicate service")
		}
	})
}
