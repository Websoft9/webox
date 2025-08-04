package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"websoft9-agent/internal/constants"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	log.Printf("Websoft9 Agent %s starting...", Version)
	log.Printf("Commit: %s, Build Time: %s", Commit, BuildTime)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动 Agent 服务
	agent := NewAgent()

	// 启动监控
	go agent.StartMonitoring(ctx)

	// 启动任务执行器
	go agent.StartTaskExecutor(ctx)

	// 启动通信管理器
	go agent.StartCommunication(ctx)

	log.Println("Websoft9 Agent started successfully")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down agent...")

	// 优雅关闭
	cancel()

	// 等待所有 goroutine 完成
	time.Sleep(constants.DefaultStartupDelay)

	log.Println("Agent exited")
}

// Agent 结构体
type Agent struct {
	ID      string
	Version string
}

// NewAgent 创建新的 Agent 实例
func NewAgent() *Agent {
	agentID := os.Getenv("AGENT_ID")
	if agentID == "" {
		agentID = "default-agent"
	}

	return &Agent{
		ID:      agentID,
		Version: Version,
	}
}

// StartMonitoring 启动监控服务
func (a *Agent) StartMonitoring(ctx context.Context) {
	log.Println("Starting monitoring service...")

	ticker := time.NewTicker(constants.DefaultHealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Monitoring service stopped")
			return
		case <-ticker.C:
			// 模拟监控数据收集
			log.Printf("Agent %s: Collecting monitoring data...", a.ID)
		}
	}
}

// StartTaskExecutor 启动任务执行器
func (a *Agent) StartTaskExecutor(ctx context.Context) {
	log.Println("Starting task executor...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Task executor stopped")
			return
		default:
			// 模拟任务执行
			time.Sleep(constants.DefaultRetryInterval)
			log.Printf("Agent %s: Checking for tasks...", a.ID)
		}
	}
}

// StartCommunication 启动通信管理器
func (a *Agent) StartCommunication(ctx context.Context) {
	log.Println("Starting communication manager...")

	ticker := time.NewTicker(constants.DefaultMetricsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Communication manager stopped")
			return
		case <-ticker.C:
			// 模拟心跳发送
			log.Printf("Agent %s: Sending heartbeat...", a.ID)
		}
	}
}

// ProcessTask 处理任务 - 修复未使用参数问题
func ProcessTask(taskType string, payload []byte) error {
	if taskType == "" {
		return fmt.Errorf("task type cannot be empty")
	}

	log.Printf("Processing task: %s with payload size: %d bytes", taskType, len(payload))
	return nil
}

func GetSystemInfo() map[string]interface{} {
	return map[string]interface{}{
		"cpu_usage":    constants.TestCPUUsage,
		"memory_usage": constants.TestMemoryUsage,
		"disk_usage":   constants.TestDiskUsage,
		"timestamp":    time.Now().Unix(),
	}
}
