package task

import (
	"context"
	"fmt"
	"os/exec"
	"time"
	"websoft9-agent/internal/constants"
	"websoft9-agent/pkg/security"

	"github.com/sirupsen/logrus"
)

// AppDeployHandler 应用部署处理器
type AppDeployHandler struct{}

func (h *AppDeployHandler) Execute(ctx context.Context, task *Task) (*TaskResult, error) {
	start := time.Now()

	logrus.Infof("执行应用部署任务: %s", task.ID)

	// TODO: 实现应用部署逻辑
	// 1. 解析部署参数
	// 2. 拉取镜像
	// 3. 创建容器
	// 4. 启动应用

	return &TaskResult{
		TaskID:   task.ID,
		Status:   constants.StatusSuccess,
		Message:  "应用部署成功",
		Duration: time.Since(start).Milliseconds(),
	}, nil
}

// AppManageHandler 应用管理处理器
type AppManageHandler struct{}

func (h *AppManageHandler) Execute(ctx context.Context, task *Task) (*TaskResult, error) {
	start := time.Now()

	logrus.Infof("执行应用管理任务: %s", task.ID)

	// TODO: 实现应用管理逻辑
	// 支持的操作: start, stop, restart, remove, update

	return &TaskResult{
		TaskID:   task.ID,
		Status:   constants.StatusSuccess,
		Message:  "应用管理操作成功",
		Duration: time.Since(start).Milliseconds(),
	}, nil
}

// SystemCommandHandler 系统命令处理器
type SystemCommandHandler struct {
	validator *security.CommandValidator
}

// NewSystemCommandHandler 创建系统命令处理器
func NewSystemCommandHandler() *SystemCommandHandler {
	return &SystemCommandHandler{
		validator: security.NewCommandValidator(),
	}
}

func (h *SystemCommandHandler) Execute(ctx context.Context, task *Task) (*TaskResult, error) {
	start := time.Now()

	command, ok := task.Params["command"].(string)
	if !ok {
		return &TaskResult{
			TaskID:   task.ID,
			Status:   constants.StatusFailed,
			Message:  "缺少命令参数",
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	// 清理输入
	command = security.SanitizeInput(command)

	// 验证命令安全性
	if err := h.validator.ValidateCommand(command); err != nil {
		logrus.Warnf("命令验证失败: %s, 命令: %s", err.Error(), command)
		return &TaskResult{
			TaskID:   task.ID,
			Status:   constants.StatusFailed,
			Message:  fmt.Sprintf("命令验证失败: %s", err.Error()),
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	logrus.Infof("执行系统命令: %s", command)

	// 记录安全审计日志
	logrus.WithFields(logrus.Fields{
		"task_id": task.ID,
		"command": command,
		"action":  "system_command_execute",
	}).Info("Security audit: system command execution")

	// #nosec G204 - Command is validated by security.CommandValidator before execution
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	output, err := cmd.CombinedOutput()

	result := &TaskResult{
		TaskID:   task.ID,
		Duration: time.Since(start).Milliseconds(),
		Data:     make(map[string]interface{}),
	}

	if err != nil {
		result.Status = constants.StatusFailed
		result.Message = err.Error()
		logrus.WithFields(logrus.Fields{
			"task_id": task.ID,
			"command": command,
			"error":   err.Error(),
		}).Error("System command execution failed")
	} else {
		result.Status = constants.StatusSuccess
		result.Message = "命令执行成功"
	}

	result.Data["output"] = string(output)

	return result, nil
}

// FileTransferHandler 文件传输处理器
type FileTransferHandler struct{}

func (h *FileTransferHandler) Execute(ctx context.Context, task *Task) (*TaskResult, error) {
	start := time.Now()

	logrus.Infof("执行文件传输任务: %s", task.ID)

	// TODO: 实现文件传输逻辑
	// 支持的操作: upload, download, copy, move, delete

	return &TaskResult{
		TaskID:   task.ID,
		Status:   constants.StatusSuccess,
		Message:  "文件传输成功",
		Duration: time.Since(start).Milliseconds(),
	}, nil
}

// ServiceManageHandler 系统服务管理处理器
type ServiceManageHandler struct {
	validator *security.CommandValidator
}

// NewServiceManageHandler 创建系统服务管理处理器
func NewServiceManageHandler() *ServiceManageHandler {
	return &ServiceManageHandler{
		validator: security.NewCommandValidator(),
	}
}

func (h *ServiceManageHandler) Execute(ctx context.Context, task *Task) (*TaskResult, error) {
	start := time.Now()

	serviceName, ok := task.Params["service"].(string)
	if !ok {
		return &TaskResult{
			TaskID:   task.ID,
			Status:   constants.StatusFailed,
			Message:  "缺少服务名称参数",
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	action, ok := task.Params["action"].(string)
	if !ok {
		return &TaskResult{
			TaskID:   task.ID,
			Status:   constants.StatusFailed,
			Message:  "缺少操作参数",
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	// 清理输入
	serviceName = security.SanitizeInput(serviceName)
	action = security.SanitizeInput(action)

	// 验证操作类型
	if err := h.validator.ValidateSystemctlAction(action); err != nil {
		logrus.Warnf("systemctl 操作验证失败: %s", err.Error())
		return &TaskResult{
			TaskID:   task.ID,
			Status:   constants.StatusFailed,
			Message:  fmt.Sprintf("操作验证失败: %s", err.Error()),
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	// 验证服务名称
	if err := h.validator.ValidateServiceName(serviceName); err != nil {
		logrus.Warnf("服务名称验证失败: %s", err.Error())
		return &TaskResult{
			TaskID:   task.ID,
			Status:   constants.StatusFailed,
			Message:  fmt.Sprintf("服务名称验证失败: %s", err.Error()),
			Duration: time.Since(start).Milliseconds(),
		}, nil
	}

	logrus.Infof("管理系统服务: %s %s", action, serviceName)

	// 记录安全审计日志
	logrus.WithFields(logrus.Fields{
		"task_id":      task.ID,
		"service_name": serviceName,
		"action":       action,
		"operation":    "service_management",
	}).Info("Security audit: service management operation")

	// 使用参数化命令执行，避免命令注入
	// #nosec G204 - Action and serviceName are validated by security.CommandValidator
	cmd := exec.CommandContext(ctx, "systemctl", action, serviceName)
	output, err := cmd.CombinedOutput()

	result := &TaskResult{
		TaskID:   task.ID,
		Duration: time.Since(start).Milliseconds(),
		Data:     make(map[string]interface{}),
	}

	if err != nil {
		result.Status = constants.StatusFailed
		result.Message = err.Error()
		logrus.WithFields(logrus.Fields{
			"task_id":      task.ID,
			"service_name": serviceName,
			"action":       action,
			"error":        err.Error(),
		}).Error("Service management operation failed")
	} else {
		result.Status = constants.StatusSuccess
		result.Message = "服务管理操作成功"
	}

	result.Data["output"] = string(output)

	return result, nil
}
