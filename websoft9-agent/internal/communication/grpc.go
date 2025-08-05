package communication

import (
	"context"
	"fmt"

	"websoft9-agent/internal/config"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient gRPC 客户端
type GRPCClient struct {
	config *config.Config
	conn   *grpc.ClientConn

	// TODO: 添加具体的 gRPC 服务客户端
	// agentServiceClient pb.AgentServiceClient
}

// NewGRPCClient 创建 gRPC 客户端
func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	return &GRPCClient{
		config: cfg,
	}, nil
}

// Start 启动 gRPC 客户端
func (c *GRPCClient) Start(ctx context.Context) error {
	serverAddr := fmt.Sprintf("%s:%d", c.config.Server.Host, c.config.Server.Port)

	logrus.Infof("连接到服务端: %s", serverAddr)

	// 创建连接选项
	var opts []grpc.DialOption
	// TODO: 配置 TLS 证书
	// 目前统一使用不安全连接，后续根据配置决定是否启用 TLS
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 建立连接 - 使用 grpc.NewClient 替代废弃的 grpc.DialContext
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return fmt.Errorf("gRPC 连接失败: %v", err)
	}

	c.conn = conn

	// TODO: 初始化服务客户端
	// c.agentServiceClient = pb.NewAgentServiceClient(conn)

	logrus.Info("gRPC 客户端启动成功")
	return nil
}

// Stop 停止 gRPC 客户端
func (c *GRPCClient) Stop() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close gRPC connection")
		}
	}
}

// SendHeartbeat 发送心跳
func (c *GRPCClient) SendHeartbeat() error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现心跳发送
	// 需要实现以下功能：
	// 1. 创建带超时的上下文
	// 2. 构建心跳请求
	// 3. 调用 gRPC 心跳接口
	logrus.Debug("发送心跳到服务端")

	return nil
}

// SendMetrics 发送监控指标
func (c *GRPCClient) SendMetrics(metrics interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现指标发送
	// 需要实现以下功能：
	// 1. 创建带超时的上下文
	// 2. 转换指标数据格式
	// 3. 调用 gRPC 指标发送接口
	logrus.Debug("发送监控指标到服务端")

	return nil
}

// SendTaskResult 发送任务结果
func (c *GRPCClient) SendTaskResult(result interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现任务结果发送
	// 需要实现以下功能：
	// 1. 创建带超时的上下文
	// 2. 转换任务结果格式
	// 3. 调用 gRPC 任务结果发送接口
	logrus.Debug("发送任务结果到服务端")

	return nil
}

// ReceiveTasks 接收任务指令
func (c *GRPCClient) ReceiveTasks(ctx context.Context) error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现任务接收
	// 需要实现以下功能：
	// 1. 创建任务流请求
	// 2. 建立 gRPC 流连接
	// 3. 循环接收任务并处理
	logrus.Info("开始接收任务指令...")

	return nil
}

// TODO: handleTask 处理任务 - 当实现任务接收时需要此函数
// func (c *GRPCClient) handleTask(task interface{}) {
//     // 将任务转发给任务执行器
//     logrus.Infof("收到任务: %v", task)
// }
