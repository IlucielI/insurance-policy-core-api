# MinIO Object Storage Integration

## Overview

MinIO object storage integrated untuk document uploads (policy PDFs, claim photos, ID cards). Supports JPEG, PNG, PDF uploads dengan max size 10MB.

## Architecture

```
Client → API → MinIO Client → MinIO Server
                    ↓
              PostgreSQL (metadata)
```

### Document Types Supported

1. **Policy Documents** → `policies/` folder
   - `policy_certificate`: Policy certificates
   - `endorsement`: Policy endorsements
   - `receipt`: Payment receipts
   - `notice`: Policy notices

2. **Claim Documents** → `claims/` folder
   - `claim_photo`: Accident/damage photos
   - `medical_report`: Medical reports
   - `police_report`: Police reports
   - `invoice`: Medical/repair invoices

3. **Identity Documents** → `identities/` folder
   - `id_card`: Generic ID cards
   - `ktp`: Indonesian KTP
   - `sim`: Indonesian SIM (driver's license)

## Configuration

### Environment Variables

```bash
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=insurance-documents
MINIO_USE_SSL=false
```

### Docker Compose Setup

MinIO container already configured in `docker-compose.yml`:

```bash
# Start MinIO
docker-compose up -d minio

# Check MinIO health
curl http://localhost:9000/minio/health/live

# Access MinIO Console
# http://localhost:9001
# Username: minioadmin
# Password: minioadmin
```

## API Endpoints

### 1. Upload Document

**POST** `/api/v1/documents/upload`

**Headers:**
- `Content-Type: multipart/form-data`
- `Authorization: Bearer <JWT_TOKEN>` (if auth enabled)

**Form Fields:**
- `file` (required): File to upload (JPEG, PNG, PDF, max 10MB)
- `document_type` (required): Type of document (e.g., `policy_certificate`, `claim_photo`, `ktp`)
- `title` (required): Document title
- `description` (optional): Document description
- `policy_id` (optional): Associated policy ID

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/documents/upload \
  -F "file=@policy.pdf" \
  -F "document_type=policy_certificate" \
  -F "title=Life Insurance Policy #12345" \
  -F "description=Policy certificate for John Doe" \
  -F "policy_id=550e8400-e29b-41d4-a716-446655440000"
```

**Response:**
```json
{
  "message": "Document uploaded successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440002",
    "policy_id": "550e8400-e29b-41d4-a716-446655440000",
    "document_type": "policy_certificate",
    "title": "Life Insurance Policy #12345",
    "description": "Policy certificate for John Doe",
    "file_name": "policy.pdf",
    "file_path": "policies/a1b2c3d4-e5f6-7890-abcd-ef1234567890.pdf",
    "file_size": 102400,
    "mime_type": "application/pdf",
    "uploaded_at": "2026-09-04T23:00:00Z"
  }
}
```

### 2. List Documents

**GET** `/api/v1/documents?limit=10&offset=0`

**Headers:**
- `Authorization: Bearer <JWT_TOKEN>` (if auth enabled)

**Response:**
```json
{
  "data": [
    {
      "id": "...",
      "title": "Life Insurance Policy #12345",
      "file_name": "policy.pdf",
      "file_size": 102400,
      "uploaded_at": "2026-09-04T23:00:00Z"
    }
  ],
  "total": 5,
  "limit": 10,
  "offset": 0
}
```

### 3. Download Document

**GET** `/api/v1/documents/:id/download`

**Headers:**
- `Authorization: Bearer <JWT_TOKEN>` (if auth enabled)

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "file_name": "policy.pdf",
  "file_size": 102400,
  "mime_type": "application/pdf",
  "download_url": "http://localhost:9000/insurance-documents/policies/a1b2c3d4.pdf?X-Amz-Algorithm=...",
  "title": "Life Insurance Policy #12345",
  "uploaded_at": "2026-09-04T23:00:00Z"
}
```

**Note:** `download_url` is a presigned URL valid for 1 hour.

### 4. Delete Document

**DELETE** `/api/v1/documents/:id`

**Headers:**
- `Authorization: Bearer <JWT_TOKEN>` (if auth enabled)

**Response:**
```json
{
  "message": "Document deleted successfully"
}
```

## File Validation

### Allowed MIME Types
- `image/jpeg`, `image/jpg`
- `image/png`
- `application/pdf`

### File Size Limit
- Maximum: 10MB (10,485,760 bytes)

### Validation Errors
```json
{
  "error": "file too large: max size is 10MB, got 15728640 bytes"
}
```

```json
{
  "error": "unsupported file type: image/gif (allowed: JPEG, PNG, PDF)"
}
```

## Code Structure

```
internal/
├── infrastructure/
│   └── storage/
│       └── minio_client.go       # MinIO client implementation
├── usecase/
│   └── document_usecase.go       # Business logic + validation
├── delivery/
│   └── http/
│       └── document_handler.go   # HTTP handlers
└── repository/
    └── document_repository.go    # Database operations
```

## Testing

### Run Test Script

```bash
# Start services
docker-compose up -d postgres redis minio

# Run API (in another terminal)
go run cmd/api/main.go

# Run test script
./test-minio.sh
```

### Manual Testing

1. **Access MinIO Console:**
   - URL: http://localhost:9001
   - Username: `minioadmin`
   - Password: `minioadmin`

2. **Upload a file via cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/documents/upload \
  -F "file=@test.pdf" \
  -F "document_type=policy_certificate" \
  -F "title=Test Document"
```

3. **Download a file:**
```bash
# Get document ID from upload response
DOC_ID="550e8400-e29b-41d4-a716-446655440001"

# Get presigned download URL
curl http://localhost:8080/api/v1/documents/$DOC_ID/download

# Download file using presigned URL
curl -o downloaded.pdf "<presigned_url>"
```

4. **Verify in MinIO Console:**
   - Navigate to Buckets → `insurance-documents`
   - Check uploaded files in `policies/`, `claims/`, or `identities/` folders

## Integration with Other Modules

### Claims Module
```go
// Upload claim photos
POST /api/v1/documents/upload
{
  "file": claim-photo.jpg,
  "document_type": "claim_photo",
  "title": "Accident Photo - Front Bumper",
  "policy_id": "claim-uuid"
}
```

### Policy Module
```go
// Upload policy certificate
POST /api/v1/documents/upload
{
  "file": policy-cert.pdf,
  "document_type": "policy_certificate",
  "title": "Policy Certificate",
  "policy_id": "policy-uuid"
}
```

## Production Considerations

### Security

1. **Enable SSL/TLS:**
```bash
MINIO_USE_SSL=true
MINIO_ENDPOINT=minio.yourdomain.com:443
```

2. **Strong Credentials:**
```bash
MINIO_ACCESS_KEY=<strong-access-key>
MINIO_SECRET_KEY=<strong-secret-key>
```

3. **Add Authentication Middleware:**
```go
documents := api.Group("/documents")
documents.Use(middleware.JWTAuth()) // Add JWT middleware
documents.Post("/upload", documentHandler.UploadDocument)
```

### Scalability

1. **Distributed MinIO Setup:**
   - Use MinIO distributed mode for HA
   - Configure multiple MinIO nodes

2. **CDN Integration:**
   - Put CloudFlare/CloudFront in front of MinIO
   - Cache presigned URLs

3. **Monitoring:**
   - Monitor bucket size
   - Set alerts for quota limits
   - Track upload/download metrics

### Backup

```bash
# Backup MinIO data
docker run --rm -v minio_data:/data -v /backup:/backup \
  alpine tar czf /backup/minio-backup-$(date +%Y%m%d).tar.gz /data

# Restore MinIO data
docker run --rm -v minio_data:/data -v /backup:/backup \
  alpine tar xzf /backup/minio-backup-20260904.tar.gz -C /
```

## Troubleshooting

### MinIO Not Starting
```bash
# Check logs
docker-compose logs minio

# Restart container
docker-compose restart minio
```

### Connection Refused
```bash
# Check if MinIO is running
docker ps | grep minio

# Check MinIO health
curl http://localhost:9000/minio/health/live
```

### Upload Failed
```bash
# Check API logs for errors
# Common issues:
# - File too large (>10MB)
# - Invalid file type
# - MinIO credentials incorrect
# - Bucket doesn't exist (should auto-create)
```

### Presigned URL Expired
- Presigned URLs valid for 1 hour
- Generate new URL by calling download endpoint again

## Future Enhancements

1. **Image Compression:**
   - Auto-compress large images before upload
   - Generate thumbnails for photos

2. **Virus Scanning:**
   - Integrate ClamAV for uploaded files
   - Reject infected files

3. **Document OCR:**
   - Extract text from scanned documents
   - Auto-populate form fields from ID cards

4. **Version Control:**
   - Support document versioning
   - Track document history

5. **Batch Operations:**
   - Bulk upload multiple files
   - Zip download for multiple documents

## References

- MinIO Go SDK: https://min.io/docs/minio/linux/developers/go/minio-go.html
- MinIO Docker Setup: https://min.io/docs/minio/container/index.html
- Presigned URLs: https://min.io/docs/minio/linux/developers/go/API.html#PresignedGetObject
