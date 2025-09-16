package middle

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime,omitempty"`
}

var startTime = time.Now()

// HealthCheck 健康检查中间件
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Version:   "1.0.0",
			Uptime:    time.Since(startTime).String(),
		})
	}
}

// ReadyCheck 就绪检查
func ReadyCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加数据库、Redis等依赖检查
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now(),
		})
	}
}

// Metrics 指标端点
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"timestamp":     time.Now(),
			"uptime":        time.Since(startTime).String(),
			"goroutines":    runtime.NumGoroutine(),
			"memory_alloc": fmt.Sprintf("%d MB", runtime.MemStats{}.Alloc/1024/1024),
		})
	}
}