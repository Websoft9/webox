package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Agent represents the Websoft9 agent
type Agent struct {
	ID     string
	Status string
}

// NewAgent creates a new agent instance
func NewAgent(id string) *Agent {
	return &Agent{
		ID:     id,
		Status: "initialized",
	}
}

// Start starts the agent
func (a *Agent) Start(ctx context.Context) error {
	a.Status = "running"
	log.Printf("Agent %s started", a.ID)

	// Simulate agent work
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.Status = "stopped"
			log.Printf("Agent %s stopped", a.ID)
			return ctx.Err()
		case <-ticker.C:
			log.Printf("Agent %s heartbeat", a.ID)
		}
	}
}

// GetStatus returns the current status
func (a *Agent) GetStatus() string {
	return a.Status
}

// ExecuteTask executes a task
func (a *Agent) ExecuteTask(task string) error {
	if a.Status != "running" {
		return fmt.Errorf("agent is not running")
	}

	log.Printf("Agent %s executing task: %s", a.ID, task)
	// Simulate task execution
	time.Sleep(100 * time.Millisecond)
	log.Printf("Agent %s completed task: %s", a.ID, task)

	return nil
}

func main() {
	agent := NewAgent("agent-001")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := agent.Start(ctx); err != nil {
		log.Printf("Agent stopped: %v", err)
	}
}
