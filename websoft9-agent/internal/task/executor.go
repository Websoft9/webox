package task

import (
	"context"
	"sync"
	"time"

	"websoft9-agent/internal/config"

	"github.com/sirupsen/logrus"
)

// Executor 任务执行器
type Executor struct {
	config *config.Config

	// 任务处理器
	handlers map[string]TaskHandler

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// TaskHandler 任务处理器接口
type TaskHandler interface {
	Execute(ctx context.Context, task *Task) (*TaskResult, error)
}

// Task 任务定义
type Task struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Params   map[string]interface{} `json:"params"`
	Timeout  int                    `json:"timeout"`
	Priority int                    `json:"priority"`
}

// TaskResult 任务执行结果
type TaskResult struct {
	TaskID   string                 `json:"task_id"`
	Status   string                 `json:"status"` // success, failed, timeout
	Message  string                 `json:"message"`
	Data     map[string]interface{} `json:"data"`
	Duration int64                  `json:"duration"` // 执行时间(毫秒)
}

// NewExecutor 创建任务执行器
func NewExecutor(cfg *config.Config) (*Executor, error) {
	executor := &Executor{
		config:   cfg,
		handlers: make(map[string]TaskHandler),
	}

	// 注册任务处理器
	executor.registerHandlers()

	return executor, nil
}

// Start 启动任务执行器
func (e *Executor) Start(ctx context.Context) error {
	e.ctx, e.cancel = context.WithCancel(ctx)

	logrus.Info("启动任务执行器...")

	// 启动任务监听
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.listenTasks()
	}()

	return nil
}

// Stop 停止任务执行器
func (e *Executor) Stop() {
	if e.cancel != nil {
		e.cancel()
	}
	e.wg.Wait()
}

// registerHandlers 注册任务处理器
func (e *Executor) registerHandlers() {
	e.handlers["deploy_app"] = &AppDeployHandler{}
	e.handlers["manage_app"] = &AppManageHandler{}
	e.handlers["system_command"] = &SystemCommandHandler{}
	e.handlers["file_transfer"] = &FileTransferHandler{}
	e.handlers["service_manage"] = &ServiceManageHandler{}
}

// listenTasks 监听任务
func (e *Executor) listenTasks() {
	logrus.Info("开始监听任务...")

	for {
		select {
		case <-e.ctx.Done():
			return
		default:
			// TODO: 从 Redis 队列获取任务
			// task := getTaskFromQueue()
			// e.executeTask(task)
		}
	}
}

// executeTask 执行任务
func (e *Executor) executeTask(task *Task) {
	logrus.Infof("执行任务: %s (类型: %s)", task.ID, task.Type)

	handler, exists := e.handlers[task.Type]
	if !exists {
		logrus.Errorf("未知的任务类型: %s", task.Type)
		return
	}

	// 创建任务上下文
	taskCtx := e.ctx
	if task.Timeout > 0 {
		var cancel context.CancelFunc
		taskCtx, cancel = context.WithTimeout(e.ctx, time.Duration(task.Timeout)*time.Second)
		defer cancel()
	}

	// 执行任务
	result, err := handler.Execute(taskCtx, task)
	if err != nil {
		logrus.Errorf("任务执行失败: %v", err)
		result = &TaskResult{
			TaskID:  task.ID,
			Status:  "failed",
			Message: err.Error(),
		}
	}

	// TODO: 发送结果到服务端
	logrus.Infof("任务 %s 执行完成: %s", task.ID, result.Status)
}
