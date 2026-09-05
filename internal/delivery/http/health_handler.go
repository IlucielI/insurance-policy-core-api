package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/cache"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/storage"
	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	db          *sql.DB
	redisClient *cache.RedisClient
	minioClient *storage.MinIOClient
	llmBaseURL  string
}

type ComponentStatus struct {
	Status  string                 `json:"status"` // healthy, degraded, unhealthy
	Message string                 `json:"message,omitempty"`
	Latency string                 `json:"latency,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type HealthResponse struct {
	Status     string                     `json:"status"` // healthy, degraded, unhealthy
	Timestamp  string                     `json:"timestamp"`
	Version    string                     `json:"version,omitempty"`
	Components map[string]ComponentStatus `json:"components"`
}

type MetricsResponse struct {
	Timestamp  string                 `json:"timestamp"`
	Uptime     string                 `json:"uptime"`
	Components map[string]interface{} `json:"components"`
}

var startTime = time.Now()

func NewHealthHandler(db *sql.DB, redisClient *cache.RedisClient, minioClient *storage.MinIOClient, llmBaseURL string) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
		minioClient: minioClient,
		llmBaseURL:  llmBaseURL,
	}
}

// LivenessHealth returns basic liveness check
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "insurance-policy-core-api",
	})
}

// Readiness returns full dependency check
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	components := make(map[string]ComponentStatus)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Check database
	wg.Add(1)
	go func() {
		defer wg.Done()
		status := h.checkDatabase(ctx)
		mu.Lock()
		components["database"] = status
		mu.Unlock()
	}()

	// Check Redis
	wg.Add(1)
	go func() {
		defer wg.Done()
		status := h.checkRedis(ctx)
		mu.Lock()
		components["redis"] = status
		mu.Unlock()
	}()

	// Check MinIO
	wg.Add(1)
	go func() {
		defer wg.Done()
		status := h.checkMinIO(ctx)
		mu.Lock()
		components["minio"] = status
		mu.Unlock()
	}()

	// Check LLM API
	wg.Add(1)
	go func() {
		defer wg.Done()
		status := h.checkLLM(ctx)
		mu.Lock()
		components["llm"] = status
		mu.Unlock()
	}()

	wg.Wait()

	// Determine overall status
	overallStatus := "healthy"
	unhealthyCount := 0
	degradedCount := 0

	for _, comp := range components {
		if comp.Status == "unhealthy" {
			unhealthyCount++
		} else if comp.Status == "degraded" {
			degradedCount++
		}
	}

	if unhealthyCount > 0 {
		overallStatus = "unhealthy"
	} else if degradedCount > 0 {
		overallStatus = "degraded"
	}

	response := HealthResponse{
		Status:     overallStatus,
		Timestamp:  time.Now().Format(time.RFC3339),
		Version:    "1.0.0",
		Components: components,
	}

	statusCode := fiber.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = fiber.StatusServiceUnavailable
	}

	// Log failures for monitoring
	if overallStatus != "healthy" {
		responseJSON, _ := json.Marshal(response)
		fmt.Printf("[HEALTH-ERROR] Health check failed: %s\n", string(responseJSON))
	}

	return c.Status(statusCode).JSON(response)
}

func (h *HealthHandler) checkDatabase(ctx context.Context) ComponentStatus {
	start := time.Now()

	if h.db == nil {
		return ComponentStatus{
			Status:  "unhealthy",
			Message: "database client not initialized",
		}
	}

	err := h.db.PingContext(ctx)
	latency := time.Since(start)

	if err != nil {
		return ComponentStatus{
			Status:  "unhealthy",
			Message: fmt.Sprintf("ping failed: %v", err),
			Latency: latency.String(),
		}
	}

	stats := h.db.Stats()
	details := map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
	}

	if latency > 500*time.Millisecond {
		return ComponentStatus{
			Status:  "degraded",
			Message: "high latency",
			Latency: latency.String(),
			Details: details,
		}
	}

	return ComponentStatus{
		Status:  "healthy",
		Latency: latency.String(),
		Details: details,
	}
}

func (h *HealthHandler) checkRedis(ctx context.Context) ComponentStatus {
	start := time.Now()

	if h.redisClient == nil {
		return ComponentStatus{
			Status:  "degraded",
			Message: "redis not configured (cache disabled)",
		}
	}

	err := h.redisClient.Ping(ctx)
	latency := time.Since(start)

	if err != nil {
		return ComponentStatus{
			Status:  "degraded",
			Message: fmt.Sprintf("ping failed: %v", err),
			Latency: latency.String(),
		}
	}

	poolStats := h.redisClient.PoolStats()
	details := map[string]interface{}{
		"hits":        poolStats.Hits,
		"misses":      poolStats.Misses,
		"total_conns": poolStats.TotalConns,
		"idle_conns":  poolStats.IdleConns,
		"stale_conns": poolStats.StaleConns,
	}

	if latency > 200*time.Millisecond {
		return ComponentStatus{
			Status:  "degraded",
			Message: "high latency",
			Latency: latency.String(),
			Details: details,
		}
	}

	return ComponentStatus{
		Status:  "healthy",
		Latency: latency.String(),
		Details: details,
	}
}

func (h *HealthHandler) checkMinIO(ctx context.Context) ComponentStatus {
	start := time.Now()

	if h.minioClient == nil {
		return ComponentStatus{
			Status:  "unhealthy",
			Message: "minio client not initialized",
		}
	}

	bucketName := h.minioClient.GetBucketName()
	exists, err := h.minioClient.BucketExists(ctx)
	latency := time.Since(start)

	if err != nil {
		return ComponentStatus{
			Status:  "unhealthy",
			Message: fmt.Sprintf("bucket check failed: %v", err),
			Latency: latency.String(),
		}
	}

	if !exists {
		return ComponentStatus{
			Status:  "unhealthy",
			Message: fmt.Sprintf("bucket '%s' does not exist", bucketName),
			Latency: latency.String(),
		}
	}

	details := map[string]interface{}{
		"bucket": bucketName,
		"exists": exists,
	}

	if latency > 1*time.Second {
		return ComponentStatus{
			Status:  "degraded",
			Message: "high latency",
			Latency: latency.String(),
			Details: details,
		}
	}

	return ComponentStatus{
		Status:  "healthy",
		Latency: latency.String(),
		Details: details,
	}
}

func (h *HealthHandler) checkLLM(ctx context.Context) ComponentStatus {
	start := time.Now()

	if h.llmBaseURL == "" {
		return ComponentStatus{
			Status:  "degraded",
			Message: "llm base url not configured",
		}
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", h.llmBaseURL+"/health", nil)
	if err != nil {
		return ComponentStatus{
			Status:  "degraded",
			Message: fmt.Sprintf("failed to create request: %v", err),
		}
	}

	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return ComponentStatus{
			Status:  "degraded",
			Message: fmt.Sprintf("health check failed: %v", err),
			Latency: latency.String(),
		}
	}
	defer resp.Body.Close()

	details := map[string]interface{}{
		"endpoint":    h.llmBaseURL,
		"status_code": resp.StatusCode,
	}

	if resp.StatusCode != http.StatusOK {
		return ComponentStatus{
			Status:  "degraded",
			Message: fmt.Sprintf("unexpected status: %d", resp.StatusCode),
			Latency: latency.String(),
			Details: details,
		}
	}

	if latency > 2*time.Second {
		return ComponentStatus{
			Status:  "degraded",
			Message: "high latency",
			Latency: latency.String(),
			Details: details,
		}
	}

	return ComponentStatus{
		Status:  "healthy",
		Latency: latency.String(),
		Details: details,
	}
}

// Metrics returns Prometheus-compatible metrics
func (h *HealthHandler) Metrics(c *fiber.Ctx) error {
	uptime := time.Since(startTime)

	components := make(map[string]interface{})

	// Database metrics
	if h.db != nil {
		stats := h.db.Stats()
		components["database"] = map[string]interface{}{
			"open_connections":    stats.OpenConnections,
			"in_use":              stats.InUse,
			"idle":                stats.Idle,
			"wait_count":          stats.WaitCount,
			"wait_duration_ms":    stats.WaitDuration.Milliseconds(),
			"max_idle_closed":     stats.MaxIdleClosed,
			"max_lifetime_closed": stats.MaxLifetimeClosed,
		}
	}

	// Redis metrics
	if h.redisClient != nil {
		poolStats := h.redisClient.PoolStats()
		components["redis"] = map[string]interface{}{
			"hits":        poolStats.Hits,
			"misses":      poolStats.Misses,
			"timeouts":    poolStats.Timeouts,
			"total_conns": poolStats.TotalConns,
			"idle_conns":  poolStats.IdleConns,
			"stale_conns": poolStats.StaleConns,
		}
	}

	// Check accept header for format
	accept := c.Get("Accept")
	if accept == "application/prometheus" || c.Query("format") == "prometheus" {
		return h.prometheusFormat(c, uptime, components)
	}

	// Default JSON format
	response := MetricsResponse{
		Timestamp:  time.Now().Format(time.RFC3339),
		Uptime:     uptime.String(),
		Components: components,
	}

	return c.JSON(response)
}

func (h *HealthHandler) prometheusFormat(c *fiber.Ctx, uptime time.Duration, components map[string]interface{}) error {
	var output string

	// Uptime metric
	output += "# HELP insurance_api_uptime_seconds Application uptime\n"
	output += "# TYPE insurance_api_uptime_seconds gauge\n"
	output += fmt.Sprintf("insurance_api_uptime_seconds %.2f\n\n", uptime.Seconds())

	// Database metrics
	if dbMetrics, ok := components["database"].(map[string]interface{}); ok {
		output += "# HELP insurance_api_db_open_connections Open database connections\n"
		output += "# TYPE insurance_api_db_open_connections gauge\n"
		output += fmt.Sprintf("insurance_api_db_open_connections %v\n\n", dbMetrics["open_connections"])

		output += "# HELP insurance_api_db_in_use Database connections in use\n"
		output += "# TYPE insurance_api_db_in_use gauge\n"
		output += fmt.Sprintf("insurance_api_db_in_use %v\n\n", dbMetrics["in_use"])

		output += "# HELP insurance_api_db_idle Idle database connections\n"
		output += "# TYPE insurance_api_db_idle gauge\n"
		output += fmt.Sprintf("insurance_api_db_idle %v\n\n", dbMetrics["idle"])

		output += "# HELP insurance_api_db_wait_count Connection wait count\n"
		output += "# TYPE insurance_api_db_wait_count counter\n"
		output += fmt.Sprintf("insurance_api_db_wait_count %v\n\n", dbMetrics["wait_count"])
	}

	// Redis metrics
	if redisMetrics, ok := components["redis"].(map[string]interface{}); ok {
		output += "# HELP insurance_api_redis_hits Cache hits\n"
		output += "# TYPE insurance_api_redis_hits counter\n"
		output += fmt.Sprintf("insurance_api_redis_hits %v\n\n", redisMetrics["hits"])

		output += "# HELP insurance_api_redis_misses Cache misses\n"
		output += "# TYPE insurance_api_redis_misses counter\n"
		output += fmt.Sprintf("insurance_api_redis_misses %v\n\n", redisMetrics["misses"])

		output += "# HELP insurance_api_redis_total_conns Redis connections\n"
		output += "# TYPE insurance_api_redis_total_conns gauge\n"
		output += fmt.Sprintf("insurance_api_redis_total_conns %v\n\n", redisMetrics["total_conns"])
	}

	c.Set("Content-Type", "text/plain; version=0.0.4")
	return c.SendString(output)
}
