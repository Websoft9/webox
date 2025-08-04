package main

import (
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {
	agent := NewAgent()

	if agent == nil {
		t.Fatal("NewAgent() returned nil")
	}

	if agent.Version != Version {
		t.Errorf("Expected version %s, got %s", Version, agent.Version)
	}

	if agent.ID == "" {
		t.Error("Agent ID should not be empty")
	}
}

func TestProcessTask(t *testing.T) {
	tests := []struct {
		name        string
		taskType    string
		payload     []byte
		expectError bool
	}{
		{"valid task", "deploy", []byte("test payload"), false},
		{"empty task type", "", []byte("test payload"), true},
		{"nil payload", "deploy", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessTask(tt.taskType, tt.payload)

			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessTask(%s, %v) expected error but got none", tt.taskType, tt.payload)
				}
			} else {
				if err != nil {
					t.Errorf("ProcessTask(%s, %v) unexpected error: %v", tt.taskType, tt.payload, err)
				}
			}
		})
	}
}

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()

	if info == nil {
		t.Fatal("GetSystemInfo() returned nil")
	}

	expectedKeys := []string{"cpu_usage", "memory_usage", "disk_usage", "timestamp"}

	for _, key := range expectedKeys {
		if _, exists := info[key]; !exists {
			t.Errorf("Expected key %s not found in system info", key)
		}
	}

	// 检查时间戳是否合理（应该是最近的时间）
	if timestamp, ok := info["timestamp"].(int64); ok {
		now := time.Now().Unix()
		if timestamp < now-10 || timestamp > now+10 {
			t.Errorf("Timestamp %d seems unreasonable (current: %d)", timestamp, now)
		}
	} else {
		t.Error("Timestamp should be int64")
	}
}

// 基准测试
func BenchmarkProcessTask(b *testing.B) {
	payload := []byte("benchmark payload")

	for i := 0; i < b.N; i++ {
		ProcessTask("benchmark", payload)
	}
}

func BenchmarkGetSystemInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetSystemInfo()
	}
}
