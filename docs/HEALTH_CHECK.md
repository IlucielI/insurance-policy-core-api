# Health Check & Metrics Endpoints

## Endpoints

### 1. `/health` - Liveness Probe
Cek dasar apakah service berjalan. Return 200 OK selama proses masih hidup.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-09-05T10:15:00+07:00",
  "service": "insurance-policy-core-api"
}
```

**HTTP Status:** 200 OK (selalu)

### 2. `/health/ready` - Readiness Probe
Cek lengkap semua dependency: database, Redis, MinIO, LLM API. Return JSON detail per komponen.

**Response (sehat):**
```json
{
  "status": "healthy",
  "timestamp": "2026-09-05T10:15:00+07:00",
  "version": "1.0.0",
  "components": {
    "database": {
      "status": "healthy",
      "latency": "2.3ms",
      "details": {
        "open_connections": 5,
        "in_use": 1,
        "idle": 4
      }
    },
    "redis": {
      "status": "healthy",
      "latency": "0.8ms",
      "details": {
        "hits": 1234,
        "misses": 56,
        "idle_conns": 8,
        "stale_conns": 0,
        "total_conns": 10
      }
    },
    "minio": {
      "status": "healthy",
      "latency": "45ms",
      "details": {
        "bucket": "insurance-documents",
        "exists": true
      }
    },
    "llm": {
      "status": "healthy",
      "latency": "120ms",
      "details": {
        "endpoint": "http://100.103.220.104:20128/v1",
        "status_code": 200
      }
    }
  }
}
```

**Response (degraded - Redis mati):**
```json
{
  "status": "degraded",
  "timestamp": "2026-09-05T10:15:00+07:00",
  "version": "1.0.0",
  "components": {
    "database": {
      "status": "healthy",
      "latency": "2.1ms"
    },
    "redis": {
      "status": "degraded",
      "message": "redis client not configured (cache disabled)"
    },
    ...
  }
}
```

**HTTP Status:** 200 OK (healthy/degraded) atau 503 Service Unavailable (unhealthy)

### 3. `/metrics` - Metrics Endpoint (Opsional Prometheus)

Return metrik dalam format JSON default, atau format Prometheus jika Accept header diisi `application/prometheus` atau query `?format=prometheus`.

**JSON Format:**
```json
{
  "timestamp": "2026-09-05T10:15:00+07:00",
  "uptime": "15h30m45s",
  "components": {
    "database": {
      "open_connections": 5,
      "in_use": 1,
      "idle": 4,
      "wait_count": 0,
      "wait_duration_ms": 0,
      "max_idle_closed": 0,
      "max_lifetime_closed": 0
    },
    "redis": {
      "hits": 1234,
      "misses": 56,
      "timeouts": 0,
      "total_conns": 10,
      "idle_conns": 8,
      "stale_conns": 0
    }
  }
}
```

**Prometheus Format** (`curl -H "Accept: application/prometheus" /metrics`):
```
# HELP insurance_api_uptime_seconds Application uptime in seconds
# TYPE insurance_api_uptime_seconds gauge
insurance_api_uptime_seconds 55845.00

# HELP insurance_api_db_open_connections Number of open database connections
# TYPE insurance_api_db_open_connections gauge
insurance_api_db_open_connections 5

# HELP insurance_api_db_in_use Number of database connections in use
# TYPE insurance_api_db_in_use gauge
insurance_api_db_in_use 1

# HELP insurance_api_db_idle Number of idle database connections
# TYPE insurance_api_db_idle gauge
insurance_api_db_idle 4

# HELP insurance_api_redis_hits Total cache hits
# TYPE insurance_api_redis_hits counter
insurance_api_redis_hits 1234

# HELP insurance_api_redis_misses Total cache misses
# TYPE insurance_api_redis_misses counter
insurance_api_redis_misses 56
```

**HTTP Status:** 200 OK

## Status Kategori

| Status | Arti | HTTP Code |
|--------|------|-----------|
| `healthy` | Semua komponen berfungsi normal | 200 |
| `degraded` | Ada komponen bermasalah tapi masih bisa serve traffic | 200 |
| `unhealthy` | Ada komponen critical yang mati | 503 |

## Penggunaan di Docker

### Dockerfile (Docker Compose)
```dockerfile
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

### docker-compose.yml
```yaml
services:
  backend:
    build: .
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
```

## Penggunaan di Kubernetes

File deployment lengkap di `k8s-deployment.yaml`.

### Liveness Probe
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30
  timeoutSeconds: 5
  failureThreshold: 3
```

### Readiness Probe
```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 15
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

### Startup Probe (untuk startup lambat)
```yaml
startupProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 0
  periodSeconds: 10
  failureThreshold: 30  # Maks 5 menit untuk start
```

## Komponen yang Dicek

| Komponen | Liveness | Readiness | Catatan |
|----------|----------|-----------|---------|
| App Process | ✅ | ✅ | Selalu liveness |
| PostgreSQL | ❌ | ✅ | PingContext, latency check, pool stats |
| Redis | ❌ | ✅ | Ping, pool stats; degraded jika tidak dikonfigurasi |
| MinIO | ❌ | ✅ | BucketExists, latency check |
| LLM API | ❌ | ✅ | HTTP GET ke endpoint /health |
| Uptime | ❌ | ✅ | Via metrics endpoint |

## Setup Prometheus Scraping

Tambah job ini di `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'insurance-api'
    metrics_path: '/metrics'
    params:
      format: ['prometheus']
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: insurance-api
```

Atau scrape langsung:
```yaml
scrape_configs:
  - job_name: 'insurance-api'
    metrics_path: '/metrics'
    params:
      format: ['prometheus']
    static_configs:
      - targets: ['insurance-api-service:8080']
```

## Testing

```bash
# Liveness
curl http://localhost:8080/health

# Readiness  
curl http://localhost:8080/health/ready

# Metrics (JSON)
curl http://localhost:8080/metrics

# Metrics (Prometheus format)
curl -H "Accept: application/prometheus" http://localhost:8080/metrics
# atau
curl "http://localhost:8080/metrics?format=prometheus"
```
