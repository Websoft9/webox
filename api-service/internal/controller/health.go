package controller

import (
	"api-service/internal/config"
	response "api-service/internal/dto/common"
	"api-service/pkg/database"
	"api-service/pkg/errors"
	"api-service/pkg/redis"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthController handles health check endpoints
type HealthController struct {
	config *config.Config
}

// NewHealthController creates a new health controller
func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{
		config: cfg,
	}
}

// DatabaseHealth checks database connection health
// @Summary Check database health
// @Description Checks the database connection status and returns connection information
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "Database is healthy"
// @Failure 500 {object} response.Response "Failed to get database info"
// @Failure 503 {object} response.Response "Database connection failed"
// @Router /health/database [get]
func (hc *HealthController) DatabaseHealth(c *gin.Context) {
	start := time.Now()

	// Test database connection
	err := database.TestDatabaseConnection(hc.config)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		response.ServiceUnavailable(c, err)
		return
	}

	// Get database info
	dbInfo, err := database.GetDatabaseInfo(hc.config)
	if err != nil {
		response.WithErrorCode(c, errors.CodeDatabaseConnectionFailed)
		return
	}

	dbInfo["status"] = "healthy"
	dbInfo["response_time"] = duration
	dbInfo["timestamp"] = start.Unix()

	response.SuccessWithData(c, dbInfo)
}

// SystemHealth performs comprehensive system health check
// @Summary Comprehensive system health check
// @Description Performs health checks on all system components including database, Redis, and InfluxDB
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "All systems healthy"
// @Failure 503 {object} map[string]interface{} "Some systems are unhealthy"
// @Router /health/system [get]
func (hc *HealthController) SystemHealth(c *gin.Context) {
	start := time.Now()
	healthData := map[string]interface{}{
		"timestamp": start.Unix(),
		"status":    "healthy",
		"checks":    map[string]interface{}{},
	}

	overallHealthy := true
	checks := healthData["checks"].(map[string]interface{})

	// Run all health checks
	overallHealthy = hc.checkDatabase(checks) && overallHealthy
	overallHealthy = hc.checkRedis(c, checks) && overallHealthy
	overallHealthy = hc.checkInfluxDB(c, checks) && overallHealthy

	// Set overall status and response
	hc.finalizeHealthResponse(c, healthData, overallHealthy, start)
}

func (hc *HealthController) checkDatabase(checks map[string]interface{}) bool {
	dbStart := time.Now()
	dbErr := database.TestDatabaseConnection(hc.config)
	dbDuration := time.Since(dbStart).Milliseconds()

	if dbErr != nil {
		checks["database"] = map[string]interface{}{
			"status":        "unhealthy",
			"error":         dbErr.Error(),
			"response_time": dbDuration,
		}
		return false
	}

	dbInfo, _ := database.GetDatabaseInfo(hc.config)
	checks["database"] = map[string]interface{}{
		"status":        "healthy",
		"response_time": dbDuration,
		"info":          dbInfo,
	}
	return true
}

func (hc *HealthController) checkRedis(c *gin.Context, checks map[string]interface{}) bool {
	if hc.config.Redis.Host == "" {
		return true // Redis not configured, consider healthy
	}

	begin := time.Now()

	ctx := c.Request.Context()
	if _, err := redis.Ping(ctx); err != nil {
		redisDuration := time.Since(begin).Milliseconds()
		checks["redis"] = map[string]interface{}{
			"status":        "unhealthy",
			"error":         err.Error(),
			"response_time": redisDuration,
		}
		return false
	}
	redisDuration := time.Since(begin).Milliseconds()
	checks["redis"] = map[string]interface{}{
		"status":        "healthy",
		"response_time": redisDuration,
		"host":          hc.config.Redis.Host,
		"port":          hc.config.Redis.Port,
		"db":            hc.config.Redis.DB,
	}
	return true
}

func (hc *HealthController) checkInfluxDB(c *gin.Context, checks map[string]interface{}) bool {
	if hc.config.InfluxDB.URL == "" {
		return true // InfluxDB not configured, consider healthy
	}

	influxStart := time.Now()
	influxClient, err := database.InitInfluxDB(hc.config)
	influxDuration := time.Since(influxStart).Milliseconds()

	if err != nil {
		checks["influxdb"] = map[string]interface{}{
			"status":        "unhealthy",
			"error":         err.Error(),
			"response_time": influxDuration,
		}
		return false
	}

	ctx := c.Request.Context()
	health, healthErr := influxClient.Health(ctx)
	if healthErr != nil {
		checks["influxdb"] = map[string]interface{}{
			"status":        "unhealthy",
			"error":         healthErr.Error(),
			"response_time": influxDuration,
		}
		return false
	}

	checks["influxdb"] = map[string]interface{}{
		"status":        "healthy",
		"response_time": influxDuration,
		"url":           hc.config.InfluxDB.URL,
		"org":           hc.config.InfluxDB.Org,
		"bucket":        hc.config.InfluxDB.Bucket,
		"health_status": health.Status,
	}
	return true
}

func (hc *HealthController) finalizeHealthResponse(c *gin.Context, healthData map[string]interface{}, overallHealthy bool, start time.Time) {
	if !overallHealthy {
		healthData["status"] = "unhealthy"
	}

	totalDuration := time.Since(start).Milliseconds()
	healthData["total_response_time"] = totalDuration

	statusCode := http.StatusOK
	if !overallHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"success": overallHealthy,
		"message": func() string {
			if overallHealthy {
				return "All systems healthy"
			}
			return "Some systems are unhealthy"
		}(),
		"data": healthData,
	})
}

// DatabaseStats returns detailed database statistics
// @Summary Get database statistics
// @Description Returns detailed statistics about database connections, performance metrics, and configuration
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "Database statistics retrieved successfully"
// @Failure 500 {object} response.Response "Failed to get database statistics"
// @Router /health/database/stats [get]
func (hc *HealthController) DatabaseStats(c *gin.Context) {
	dbInfo, err := database.GetDatabaseInfo(hc.config)
	if err != nil {
		response.WithError(c, err)
		return
	}
	response.SuccessWithData(c, dbInfo)
}

// Ping is a simple ping endpoint
// @Summary Ping endpoint
// @Description Simple ping endpoint that returns pong with timestamp
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Ping successful"
// @Router /ping [get]
func (hc *HealthController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now().Unix(),
	})
}

// Readiness check for Kubernetes readiness probe
// @Summary Kubernetes readiness probe
// @Description Readiness probe endpoint for Kubernetes deployments. Checks if the application is ready to serve traffic
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Application is ready"
// @Failure 503 {object} map[string]interface{} "Database not ready"
// @Router /readiness [get]
func (hc *HealthController) Readiness(c *gin.Context) {
	// Check if database is ready
	if err := database.TestDatabaseConnection(hc.config); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":   false,
			"message": "Database not ready",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ready":   true,
		"message": "Application is ready",
	})
}

// Liveness check for Kubernetes liveness probe
// @Summary Kubernetes liveness probe
// @Description Liveness probe endpoint for Kubernetes deployments. Indicates if the application is alive and running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Application is alive"
// @Router /liveness [get]
func (hc *HealthController) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"alive":   true,
		"message": "Application is alive",
	})
}
