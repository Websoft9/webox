package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"websoft9-agent/internal/agent"
	"websoft9-agent/internal/config"

	"github.com/sirupsen/logrus"
)

var (
	configFile = flag.String("config", "/etc/websoft9/agent.yaml", "配置文件路径")
	version    = flag.Bool("version", false, "显示版本信息")
)

func main() {
	flag.Parse()

	if *version {
		logrus.Info("Websoft9 Agent v1.0.0")
		return
	}

	// 加载配置
	cfg, err := config.Load(*configFile)
	if err != nil {
		logrus.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	setupLogger(cfg.Log.Level)

	// 创建 Agent 实例
	agentInstance, err := agent.New(cfg)
	if err != nil {
		logrus.Fatalf("创建 Agent 失败: %v", err)
	}

	// 启动 Agent
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := agentInstance.Start(ctx); err != nil {
		logrus.Fatalf("启动 Agent 失败: %v", err)
	}

	// 等待信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logrus.Info("正在关闭 Agent...")
	agentInstance.Stop()
	logrus.Info("Agent 已关闭")
}

func setupLogger(level string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
