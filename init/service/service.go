package service

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Service 表示一个服务实例
type Service struct {
	Config   ServiceConfig
	Status   ServiceStatus
	cmd      *exec.Cmd
	eventBus *EventBus
	stopChan chan struct{}
}

// NewService 创建新的服务实例
func NewService(config ServiceConfig, eventBus *EventBus) *Service {
	return &Service{
		Config: config,
		Status: ServiceStatus{
			State:        StateUnknown,
			StartTime:    time.Time{},
			RestartCount: 0,
		},
		eventBus: eventBus,
		stopChan: make(chan struct{}),
	}
}

// Start 启动服务
func (s *Service) Start() error {
	if s.Status.State == StateRunning {
		return fmt.Errorf("service %s is already running", s.Config.Name)
	}

	s.updateState(StateStarting)

	// 准备命令
	s.cmd = exec.Command(s.Config.ExecPath, s.Config.Args...)

	// 设置环境变量
	if len(s.Config.Environment) > 0 {
		env := os.Environ()
		for k, v := range s.Config.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		s.cmd.Env = env
	}

	// 启动进程
	if err := s.cmd.Start(); err != nil {
		s.updateState(StateFailed)
		return fmt.Errorf("failed to start service %s: %v", s.Config.Name, err)
	}

	s.Status.Pid = s.cmd.Process.Pid
	s.Status.StartTime = time.Now()
	s.updateState(StateRunning)

	// 监控进程
	go s.monitor()

	return nil
}

// Stop 停止服务
func (s *Service) Stop() error {
	if s.Status.State != StateRunning {
		return fmt.Errorf("service %s is not running", s.Config.Name)
	}

	s.updateState(StateStopping)
	close(s.stopChan)

	if s.cmd != nil && s.cmd.Process != nil {
		if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil {
			// 如果 SIGTERM 失败，尝试 SIGKILL
			if err := s.cmd.Process.Kill(); err != nil {
				s.updateState(StateFailed)
				return fmt.Errorf("failed to kill service %s: %v", s.Config.Name, err)
			}
		}
	}

	s.updateState(StateStopped)
	return nil
}

// Restart 重启服务
func (s *Service) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}

	s.Status.RestartCount++
	s.stopChan = make(chan struct{}) // 重新创建停止通道

	return s.Start()
}

// monitor 监控服务进程
func (s *Service) monitor() {
	if s.cmd == nil {
		return
	}

	// 等待进程结束
	err := s.cmd.Wait()

	// 检查是否是正常停止
	select {
	case <-s.stopChan:
		// 正常停止，不需要特殊处理
		return
	default:
		// 异常停止
		s.Status.LastError = err
		s.updateState(StateFailed)

		// 根据重启策略处理
		if s.Config.Restart == "always" || (s.Config.Restart == "on-failure" && err != nil) {
			s.Status.RestartCount++
			s.stopChan = make(chan struct{})
			if err := s.Start(); err != nil {
				s.Status.LastError = err
			}
		}
	}
}

// updateState 更新服务状态并发送事件
func (s *Service) updateState(state ServiceState) {
	oldState := s.Status.State
	s.Status.State = state

	// 只有当状态真正发生变化时才发送事件
	if oldState != state && s.eventBus != nil {
		event := ServiceEvent{
			Type:      EventType(state),
			Service:   s.Config.Name,
			Data:      s.Status,
			Timestamp: time.Now(),
		}

		// 同步发送事件，确保状态更新在返回前完成
		s.eventBus.EmitSync(event)
	}
}

// GetStatus 获取服务状态
func (s *Service) GetStatus() ServiceStatus {
	return s.Status
}
