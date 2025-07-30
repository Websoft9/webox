package agent

import (
	"context"
	"sync"
	"time"

	"websoft9-agent/internal/communication"
	"websoft9-agent/internal/config"
	"websoft9-agent/internal/monitor"
	"websoft9-agent/internal/task"

	"github.com/sirupsen/logrus"
)

// Agent 主要的 Agent 结构
type Agent struct {
	config *config.Config

	// 核心组件
	monitor      *monitor.Monitor
	taskExecutor *task.Executor
	comm         *communication.Manager

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New 创建新的 Agent 实例
func New(cfg *config.Config) (*Agent, error) {
	// 创建监控组件
	mon, err := monitor.New(cfg)
	if err != nil {
		return nil, err
	}

	// 创建任务执行器
	taskExec, err := task.NewExecutor(cfg)
	if err != nil {
		return nil, err
	}

	// 创建通信管理器
	commMgr, err := communication.NewManager(cfg)
	if err != nil {
		return nil, err
	}

	return &Agent{
		config:       cfg,
		monitor:      mon,
		taskExecutor: taskExec,
		comm:         commMgr,
	}, nil
}

// Start 启动 Agent
func (a *Agent) Start(ctx context.Context) error {
	a.ctx, a.cancel = context.WithCancel(ctx)

	logrus.Info("启动 Websoft9 Agent...")

	// 启动通信管理器
	if err := a.comm.Start(a.ctx); err != nil {
		return err
	}

	// 启动监控组件
	if err := a.monitor.Start(a.ctx); err != nil {
		return err
	}

	// 启动任务执行器
	if err := a.taskExecutor.Start(a.ctx); err != nil {
		return err
	}

	// 启动心跳
	a.startHeartbeat()

	logrus.Info("Websoft9 Agent 启动完成")
	return nil
}

// Stop 停止 Agent
func (a *Agent) Stop() {
	if a.cancel != nil {
		a.cancel()
	}

	a.wg.Wait()

	// 停止各组件
	if a.comm != nil {
		a.comm.Stop()
	}
	if a.monitor != nil {
		a.monitor.Stop()
	}
	if a.taskExecutor != nil {
		a.taskExecutor.Stop()
	}
}

// startHeartbeat 启动心跳
func (a *Agent) startHeartbeat() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		ticker := time.NewTicker(time.Duration(a.config.Agent.HeartbeatInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-a.ctx.Done():
				return
			case <-ticker.C:
				a.sendHeartbeat()
			}
		}
	}()
}

// sendHeartbeat 发送心跳
func (a *Agent) sendHeartbeat() {
	logrus.Debug("发送心跳")
	// TODO: 实现心跳发送逻辑
}
