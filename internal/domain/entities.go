package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone,omitempty"`
	Role         string    `json:"role"` // customer, admin, underwriter
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Product struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Slug            string                 `json:"slug"`
	Category        string                 `json:"category"` // life, health, vehicle
	Description     string                 `json:"description"`
	CoverageDetails map[string]interface{} `json:"coverage_details"`
	MinSumAssured   int64                  `json:"min_sum_assured"`
	MaxSumAssured   int64                  `json:"max_sum_assured"`
	MinPaymentTerm  int                    `json:"min_payment_term"` // months
	MaxPaymentTerm  int                    `json:"max_payment_term"`
	BasePremiumRate float64                `json:"base_premium_rate"` // percentage
	AgeFactor       map[string]float64     `json:"age_factor,omitempty"`
	IsActive        bool                   `json:"is_active"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type Application struct {
	ID                string                 `json:"id"`
	UserID            string                 `json:"user_id"`
	ProductID         string                 `json:"product_id"`
	ApplicantData     map[string]interface{} `json:"applicant_data"`
	SumAssured        int64                  `json:"sum_assured"`
	PaymentTerm       int                    `json:"payment_term"` // months
	PremiumAmount     int64                  `json:"premium_amount"`
	HealthQuestions   map[string]interface{} `json:"health_questions,omitempty"`
	Status            string                 `json:"status"` // draft, submitted, under_review, approved, rejected
	UnderwriterID     *string                `json:"underwriter_id,omitempty"`
	UnderwriterNotes  string                 `json:"underwriter_notes,omitempty"`
	RejectionReason   string                 `json:"rejection_reason,omitempty"`
	SubmittedAt       *time.Time             `json:"submitted_at,omitempty"`
	ReviewedAt        *time.Time             `json:"reviewed_at,omitempty"`
	ApprovedAt        *time.Time             `json:"approved_at,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

type Payment struct {
	ID                    string     `json:"id"`
	ApplicationID         string     `json:"application_id"`
	OrderID               string     `json:"order_id"`
	MidtransTransactionID string     `json:"midtrans_transaction_id,omitempty"`
	PaymentType           string     `json:"payment_type,omitempty"`
	GrossAmount           int64      `json:"gross_amount"`
	Status                string     `json:"status"` // pending, success, failed, expired
	PaidAt                *time.Time `json:"paid_at,omitempty"`
	ExpiredAt             *time.Time `json:"expired_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type ChatSession struct {
	ID        string    `json:"id"`
	UserID    *string   `json:"user_id,omitempty"`
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatMessage struct {
	ID            string                   `json:"id"`
	ChatSessionID string                   `json:"chat_session_id"`
	Role          string                   `json:"role"` // user, assistant
	Content       string                   `json:"content"`
	ContextDocs   []map[string]interface{} `json:"context_docs,omitempty"`
	CreatedAt     time.Time                `json:"created_at"`
}

type ProductEmbedding struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	ChunkType string    `json:"chunk_type"` // description, benefits, exclusions, faq
	ChunkText string    `json:"chunk_text"`
	Embedding []float32 `json:"-"` // Don't expose raw embedding in JSON
	CreatedAt time.Time `json:"created_at"`
}

type ActivityLog struct {
	ID         string                 `json:"id"`
	UserID     *string                `json:"user_id,omitempty"`
	Action     string                 `json:"action"` // application_status_changed, product_created, etc.
	EntityType string                 `json:"entity_type"`
	EntityID   string                 `json:"entity_id"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}
