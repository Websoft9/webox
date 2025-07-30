package monitor

import (
	"websoft9-agent/internal/config"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/sirupsen/logrus"
)

// SystemMonitor 系统监控器
type SystemMonitor struct {
	config *config.Config
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPU     CPUMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	Disk    []DiskMetrics  `json:"disk"`
	Network NetworkMetrics `json:"network"`
}

// CPUMetrics CPU指标
type CPUMetrics struct {
	Usage   float64 `json:"usage"`
	Cores   int     `json:"cores"`
	LoadAvg float64 `json:"load_avg"`
}

// MemoryMetrics 内存指标
type MemoryMetrics struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Usage     float64 `json:"usage"`
}

// DiskMetrics 磁盘指标
type DiskMetrics struct {
	Device     string  `json:"device"`
	Mountpoint string  `json:"mountpoint"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Free       uint64  `json:"free"`
	Usage      float64 `json:"usage"`
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

// NewSystemMonitor 创建系统监控器
func NewSystemMonitor(cfg *config.Config) (*SystemMonitor, error) {
	return &SystemMonitor{
		config: cfg,
	}, nil
}

// Collect 采集系统指标
func (s *SystemMonitor) Collect() error {
	metrics, err := s.collectMetrics()
	if err != nil {
		return err
	}

	logrus.Debugf("系统指标: CPU使用率=%.2f%%, 内存使用率=%.2f%%",
		metrics.CPU.Usage, metrics.Memory.Usage)

	// TODO: 发送指标到 InfluxDB
	return nil
}

// collectMetrics 采集所有系统指标
func (s *SystemMonitor) collectMetrics() (*SystemMetrics, error) {
	metrics := &SystemMetrics{}

	// CPU 指标
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}
	if len(cpuPercent) > 0 {
		metrics.CPU.Usage = cpuPercent[0]
	}

	cpuCounts, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	metrics.CPU.Cores = cpuCounts

	// 内存指标
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	metrics.Memory = MemoryMetrics{
		Total:     memInfo.Total,
		Used:      memInfo.Used,
		Available: memInfo.Available,
		Usage:     memInfo.UsedPercent,
	}

	// 磁盘指标
	diskPartitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	for _, partition := range diskPartitions {
		diskUsage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		metrics.Disk = append(metrics.Disk, DiskMetrics{
			Device:     partition.Device,
			Mountpoint: partition.Mountpoint,
			Total:      diskUsage.Total,
			Used:       diskUsage.Used,
			Free:       diskUsage.Free,
			Usage:      diskUsage.UsedPercent,
		})
	}

	// 网络指标
	netIO, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
	if len(netIO) > 0 {
		metrics.Network = NetworkMetrics{
			BytesSent:   netIO[0].BytesSent,
			BytesRecv:   netIO[0].BytesRecv,
			PacketsSent: netIO[0].PacketsSent,
			PacketsRecv: netIO[0].PacketsRecv,
		}
	}

	return metrics, nil
}
