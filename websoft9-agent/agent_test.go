package main

import (
	"context"
	"testing"
	"time"
)

func TestNewAgent(t *testing.T) {
	agent := NewAgent("test-agent")
	if agent == nil {
		t.Fatal("Expected agent to be created")
	}

	if agent.ID != "test-agent" {
		t.Errorf("Expected ID to be 'test-agent', got %s", agent.ID)
	}

	if agent.Status != "initialized" {
		t.Errorf("Expected status to be 'initialized', got %s", agent.Status)
	}
}

func TestGetStatus(t *testing.T) {
	agent := NewAgent("test-agent")

	status := agent.GetStatus()
	if status != "initialized" {
		t.Errorf("Expected status 'initialized', got %s", status)
	}
}

func TestExecuteTask(t *testing.T) {
	agent := NewAgent("test-agent")

	// Test task execution when agent is not running
	err := agent.ExecuteTask("test-task")
	if err == nil {
		t.Error("Expected error when agent is not running")
	}

	// Set agent to running status
	agent.Status = "running"

	// Test successful task execution
	err = agent.ExecuteTask("test-task")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestAgentStart(t *testing.T) {
	agent := NewAgent("test-agent")

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := agent.Start(ctx)
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}

	if agent.Status != "stopped" {
		t.Errorf("Expected status 'stopped', got %s", agent.Status)
	}
}

func TestAgentStartWithTimeout(t *testing.T) {
	agent := NewAgent("test-agent")

	// Test with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := agent.Start(ctx)
	if err == nil {
		t.Error("Expected timeout error")
	}
}
