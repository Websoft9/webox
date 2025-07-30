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
	if c.config.Server.TLS {
		// TODO: 配置 TLS 证书
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 建立连接
	conn, err := grpc.DialContext(ctx, serverAddr, opts...)
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
		c.conn.Close()
	}
}

// SendHeartbeat 发送心跳
func (c *GRPCClient) SendHeartbeat() error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现心跳发送
	logrus.Debug("发送心跳到服务端")

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// req := &pb.HeartbeatRequest{
	//     AgentId:   c.config.Agent.ID,
	//     Timestamp: time.Now().Unix(),
	//     Status:    "healthy",
	// }

	// _, err := c.agentServiceClient.Heartbeat(ctx, req)
	// return err

	return nil
}

// SendMetrics 发送监控指标
func (c *GRPCClient) SendMetrics(metrics interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现指标发送
	logrus.Debug("发送监控指标到服务端")

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// req := &pb.MetricsRequest{
	//     AgentId:   c.config.Agent.ID,
	//     Timestamp: time.Now().Unix(),
	//     Metrics:   convertMetrics(metrics),
	// }

	// _, err := c.agentServiceClient.SendMetrics(ctx, req)
	// return err

	return nil
}

// SendTaskResult 发送任务结果
func (c *GRPCClient) SendTaskResult(result interface{}) error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现任务结果发送
	logrus.Debug("发送任务结果到服务端")

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// req := &pb.TaskResultRequest{
	//     AgentId: c.config.Agent.ID,
	//     Result:  convertTaskResult(result),
	// }

	// _, err := c.agentServiceClient.SendTaskResult(ctx, req)
	// return err

	return nil
}

// ReceiveTasks 接收任务指令
func (c *GRPCClient) ReceiveTasks(ctx context.Context) error {
	if c.conn == nil {
		return fmt.Errorf("gRPC 连接未建立")
	}

	// TODO: 实现任务接收
	logrus.Info("开始接收任务指令...")

	// req := &pb.TaskStreamRequest{
	//     AgentId: c.config.Agent.ID,
	// }

	// stream, err := c.agentServiceClient.ReceiveTasks(ctx, req)
	// if err != nil {
	//     return err
	// }

	// for {
	//     task, err := stream.Recv()
	//     if err != nil {
	//         return err
	//     }
	//
	//     // 处理接收到的任务
	//     go c.handleTask(task)
	// }

	return nil
}

// handleTask 处理任务
func (c *GRPCClient) handleTask(task interface{}) {
	// TODO: 将任务转发给任务执行器
	logrus.Infof("收到任务: %v", task)
}
