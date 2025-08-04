package monitor

import (
	"context"
	"sync"
	"time"

	"websoft9-agent/internal/config"
	"websoft9-agent/internal/constants"

	"github.com/sirupsen/logrus"
)

// Monitor 监控组件
type Monitor struct {
	config *config.Config

	// 监控器
	systemMonitor    *SystemMonitor
	containerMonitor *ContainerMonitor
	healthChecker    *HealthChecker

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New 创建监控组件
func New(cfg *config.Config) (*Monitor, error) {
	sysMon, err := NewSystemMonitor(cfg)
	if err != nil {
		return nil, err
	}

	containerMon, err := NewContainerMonitor(cfg)
	if err != nil {
		return nil, err
	}

	healthChecker, err := NewHealthChecker(cfg)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		config:           cfg,
		systemMonitor:    sysMon,
		containerMonitor: containerMon,
		healthChecker:    healthChecker,
	}, nil
}

// Start 启动监控
func (m *Monitor) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)

	logrus.Info("启动监控组件...")

	// 启动系统监控
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.runSystemMonitor()
	}()

	// 启动容器监控
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.runContainerMonitor()
	}()

	// 启动健康检查
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.runHealthChecker()
	}()

	return nil
}

// Stop 停止监控
func (m *Monitor) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
	m.wg.Wait()
}

// runSystemMonitor 运行系统监控
func (m *Monitor) runSystemMonitor() {
	ticker := time.NewTicker(time.Duration(m.config.Agent.MonitorInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.systemMonitor.Collect(); err != nil {
				logrus.Errorf("系统监控采集失败: %v", err)
			}
		}
	}
}

// runContainerMonitor 运行容器监控
func (m *Monitor) runContainerMonitor() {
	ticker := time.NewTicker(time.Duration(m.config.Agent.MonitorInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.containerMonitor.Collect(); err != nil {
				logrus.Errorf("容器监控采集失败: %v", err)
			}
		}
	}
}

// runHealthChecker 运行健康检查
func (m *Monitor) runHealthChecker() {
	ticker := time.NewTicker(constants.DefaultHealthCheckInterval) // 健康检查间隔
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.healthChecker.Check(); err != nil {
				logrus.Errorf("健康检查失败: %v", err)
			}
		}
	}
}
