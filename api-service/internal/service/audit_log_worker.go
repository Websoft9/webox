package service

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"api-service/pkg/logger"
)

const (
	// Worker pool configuration
	defaultWorkerPoolSize = 10               // Number of concurrent workers
	defaultJobQueueSize   = 1000             // Job queue buffer size
	auditTimeout          = 5 * time.Second  // Timeout for single audit operation
	shutdownTimeout       = 10 * time.Second // Timeout for graceful shutdown
)

// AuditJob represents a single audit log task
type AuditJob struct {
	ginCtx       *gin.Context
	responseBody []byte
	responseTime int
}

// AuditLogWorkerPool manages worker pool for async audit log recording
type AuditLogWorkerPool struct {
	jobQueue chan *AuditJob
	wg       sync.WaitGroup
	logger   logger.Logger
	service  *auditLogService
	once     sync.Once
	stopCh   chan struct{}
}

// NewAuditLogWorkerPool creates a new worker pool
func NewAuditLogWorkerPool(service *auditLogService, logger logger.Logger) *AuditLogWorkerPool {
	pool := &AuditLogWorkerPool{
		jobQueue: make(chan *AuditJob, defaultJobQueueSize),
		logger:   logger,
		service:  service,
		stopCh:   make(chan struct{}),
	}

	// Start workers
	pool.start()

	return pool
}

// start initializes and starts worker goroutines
func (p *AuditLogWorkerPool) start() {
	for i := 0; i < defaultWorkerPoolSize; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	p.logger.Info("Audit log worker pool started",
		logger.Field{Key: "pool_size", Value: defaultWorkerPoolSize},
		logger.Field{Key: "queue_size", Value: defaultJobQueueSize})
}

// worker processes audit jobs from the queue
func (p *AuditLogWorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case job, ok := <-p.jobQueue:
			if !ok {
				// Channel closed, worker should exit
				p.logger.Debug("Worker exiting",
					logger.Field{Key: "worker_id", Value: id})
				return
			}

			// Process the job with timeout
			p.processJob(job)

		case <-p.stopCh:
			// Shutdown signal received
			p.logger.Debug("Worker received shutdown signal",
				logger.Field{Key: "worker_id", Value: id})
			return
		}
	}
}

// processJob processes a single audit job with timeout control
func (p *AuditLogWorkerPool) processJob(job *AuditJob) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), auditTimeout)
	defer cancel()

	// Extract audit data safely (before gin context expires)
	auditReq := p.service.extractAuditInfoFromRequest(
		job.ginCtx,
		job.responseBody,
		job.responseTime,
	)

	// Record audit log
	if err := p.service.RecordLog(ctx, auditReq); err != nil {
		p.logger.Error("Failed to record audit log",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "url", Value: job.ginCtx.Request.RequestURI})
	}
}

// Submit submits a new audit job to the pool
// Returns false if queue is full (non-blocking)
func (p *AuditLogWorkerPool) Submit(job *AuditJob) bool {
	select {
	case p.jobQueue <- job:
		return true
	default:
		// Queue is full, log warning and drop the job
		p.logger.Warn("Audit log queue is full, dropping job",
			logger.Field{Key: "queue_size", Value: defaultJobQueueSize},
			logger.Field{Key: "url", Value: job.ginCtx.Request.RequestURI})
		return false
	}
}

// Shutdown gracefully shuts down the worker pool
func (p *AuditLogWorkerPool) Shutdown() {
	p.once.Do(func() {
		p.logger.Info("Shutting down audit log worker pool")

		// Close stop channel to signal workers
		close(p.stopCh)

		// Close job queue (no more jobs accepted)
		close(p.jobQueue)

		// Wait for workers to finish with timeout
		done := make(chan struct{})
		go func() {
			p.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			p.logger.Info("Audit log worker pool shut down successfully")
		case <-time.After(shutdownTimeout):
			p.logger.Warn("Audit log worker pool shutdown timeout",
				logger.Field{Key: "timeout", Value: shutdownTimeout})
		}
	})
}

// GetQueueLength returns current queue length (for monitoring)
func (p *AuditLogWorkerPool) GetQueueLength() int {
	return len(p.jobQueue)
}
