package service

import (
	"sync"
	"time"
)

// ServiceEvent 定义服务事件
type ServiceEvent struct {
	Type      EventType
	Service   string
	Data      interface{}
	Timestamp time.Time
}

// EventHandler 定义事件处理函数类型
type EventHandler func(event ServiceEvent)

// EventBus 实现事件总线
type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus 创建新的事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe 订阅特定类型的事件
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if _, exists := eb.handlers[eventType]; !exists {
		eb.handlers[eventType] = []EventHandler{}
	}
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Unsubscribe 取消订阅特定类型的事件
func (eb *EventBus) Unsubscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if handlers, exists := eb.handlers[eventType]; exists {
		for i, h := range handlers {
			if &h == &handler {
				eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Emit 发送事件
func (eb *EventBus) Emit(event ServiceEvent) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	// 调用所有注册的处理函数
	if handlers, exists := eb.handlers[event.Type]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}

	// 调用通用事件处理函数（如果有的话）
	if handlers, exists := eb.handlers[""]; exists {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

// EmitSync 同步发送事件
func (eb *EventBus) EmitSync(event ServiceEvent) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if handlers, exists := eb.handlers[event.Type]; exists {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
