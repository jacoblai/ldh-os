package service

import (
	"encoding/json"
	"fmt"
)

// MCPFunction 定义 MCP 功能处理函数类型
type MCPFunction func(params map[string]interface{}) (interface{}, error)

// MCPRequest 定义 MCP 请求结构
type MCPRequest struct {
	Service  string                 `json:"service"`
	Function string                 `json:"function"`
	Params   map[string]interface{} `json:"params"`
}

// MCPResponse 定义 MCP 响应结构
type MCPResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// MCPHandler MCP 协议处理器
type MCPHandler struct {
	functions map[string]map[string]MCPFunction // service -> function -> handler
}

// NewMCPHandler 创建新的 MCP 处理器
func NewMCPHandler() *MCPHandler {
	return &MCPHandler{
		functions: make(map[string]map[string]MCPFunction),
	}
}

// RegisterFunction 注册 MCP 功能
func (h *MCPHandler) RegisterFunction(service, name string, fn MCPFunction) error {
	if _, exists := h.functions[service]; !exists {
		h.functions[service] = make(map[string]MCPFunction)
	}

	if _, exists := h.functions[service][name]; exists {
		return fmt.Errorf("function %s already registered for service %s", name, service)
	}

	h.functions[service][name] = fn
	return nil
}

// HandleRequest 处理 MCP 请求
func (h *MCPHandler) HandleRequest(req *MCPRequest) *MCPResponse {
	// 检查服务是否存在
	serviceFuncs, exists := h.functions[req.Service]
	if !exists {
		return &MCPResponse{
			Success: false,
			Error:   fmt.Sprintf("service %s not found", req.Service),
		}
	}

	// 检查功能是否存在
	fn, exists := serviceFuncs[req.Function]
	if !exists {
		return &MCPResponse{
			Success: false,
			Error:   fmt.Sprintf("function %s not found in service %s", req.Function, req.Service),
		}
	}

	// 执行功能
	result, err := fn(req.Params)
	if err != nil {
		return &MCPResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return &MCPResponse{
		Success: true,
		Data:    result,
	}
}

// NotifyLLM 向 LLM 发送通知
func (h *MCPHandler) NotifyLLM(event ServiceEvent) error {
	// 将事件转换为 JSON 格式
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	// TODO: 实现实际的 LLM 通知机制
	// 这里可以通过 HTTP、WebSocket 或其他方式将事件发送给 LLM
	fmt.Printf("LLM Notification: %s\n", string(eventJSON))
	return nil
}

// GetRegisteredFunctions 获取已注册的功能列表
func (h *MCPHandler) GetRegisteredFunctions(service string) []string {
	if serviceFuncs, exists := h.functions[service]; exists {
		functions := make([]string, 0, len(serviceFuncs))
		for name := range serviceFuncs {
			functions = append(functions, name)
		}
		return functions
	}
	return nil
}
