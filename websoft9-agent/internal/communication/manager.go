package communication

import (
	"context"
	"fmt"
	"sync"

	"websoft9-agent/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// Manager 通信管理器
type Manager struct {
	config *config.Config

	// 通信组件
	grpcClient  *GRPCClient
	redisClient *redis.Client

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewManager 创建通信管理器
func NewManager(cfg *config.Config) (*Manager, error) {
	// 创建 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试 Redis 连接
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %v", err)
	}

	// 创建 gRPC 客户端
	grpcClient, err := NewGRPCClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Manager{
		config:      cfg,
		grpcClient:  grpcClient,
		redisClient: redisClient,
	}, nil
}

// Start 启动通信管理器
func (m *Manager) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)

	logrus.Info("启动通信管理器...")

	// 启动 gRPC 客户端
	if err := m.grpcClient.Start(m.ctx); err != nil {
		return err
	}

	// 启动消息队列监听
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.listenMessages()
	}()

	return nil
}

// Stop 停止通信管理器
func (m *Manager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}

	m.wg.Wait()

	if m.grpcClient != nil {
		m.grpcClient.Stop()
	}

	if m.redisClient != nil {
		if err := m.redisClient.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close Redis client")
		}
	}
}

// listenMessages 监听消息队列
func (m *Manager) listenMessages() {
	logrus.Info("开始监听消息队列...")

	// 订阅 Agent 相关的消息频道
	pubsub := m.redisClient.Subscribe(m.ctx, fmt.Sprintf("agent:%s", m.config.Agent.ID))
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-m.ctx.Done():
			return
		case msg := <-ch:
			m.handleMessage(msg)
		}
	}
}

// handleMessage 处理消息
func (m *Manager) handleMessage(msg *redis.Message) {
	logrus.Debugf("收到消息: %s", msg.Payload)

	// TODO: 解析消息并分发到相应的处理器
	// 消息类型可能包括:
	// - 任务指令
	// - 工作流任务
	// - 配置更新
	// - 控制指令
}

// SendHeartbeat 发送心跳
func (m *Manager) SendHeartbeat() error {
	return m.grpcClient.SendHeartbeat()
}

// SendMetrics 发送监控指标
func (m *Manager) SendMetrics(metrics interface{}) error {
	return m.grpcClient.SendMetrics(metrics)
}

// SendTaskResult 发送任务结果
func (m *Manager) SendTaskResult(result interface{}) error {
	return m.grpcClient.SendTaskResult(result)
}

// SendEvent 发送事件消息
func (m *Manager) SendEvent(event interface{}) error {
	// 发送到 Redis 消息队列
	eventChannel := "events"
	return m.redisClient.Publish(m.ctx, eventChannel, event).Err()
}
