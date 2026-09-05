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
	ID                 string                 `json:"id"`
	UserID             string                 `json:"user_id"`
	ProductID          string                 `json:"product_id"`
	ApplicantData      map[string]interface{} `json:"applicant_data"`
	SumAssured         int64                  `json:"sum_assured"`
	PaymentTerm        int                    `json:"payment_term"` // months
	PremiumAmount      int64                  `json:"premium_amount"`
	HealthQuestions    map[string]interface{} `json:"health_questions,omitempty"`
	Status             string                 `json:"status"` // draft, submitted, under_review, approved, rejected
	UnderwriterID      *string                `json:"underwriter_id,omitempty"`
	UnderwriterNotes   string                 `json:"underwriter_notes,omitempty"`
	RejectionReason    string                 `json:"rejection_reason,omitempty"`
	RiskScore          *int                   `json:"risk_score,omitempty"`          // 0-100 fraud risk score
	RiskLevel          string                 `json:"risk_level,omitempty"`          // low, medium, high
	FraudFlags         []string               `json:"fraud_flags,omitempty"`         // detected suspicious patterns
	RiskAnalysisDetail string                 `json:"risk_analysis_detail,omitempty"` // AI analysis explanation
	RiskAnalyzedAt     *time.Time             `json:"risk_analyzed_at,omitempty"`
	SubmittedAt        *time.Time             `json:"submitted_at,omitempty"`
	ReviewedAt         *time.Time             `json:"reviewed_at,omitempty"`
	ApprovedAt         *time.Time             `json:"approved_at,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
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

type Policy struct {
	ID                   string    `json:"id"`
	PolicyNumber         string    `json:"policy_number"`
	ApplicationID        string    `json:"application_id"`
	UserID               string    `json:"user_id"`
	ProductID            string    `json:"product_id"`
	SumAssured           int64     `json:"sum_assured"`
	PremiumAmount        int64     `json:"premium_amount"`
	PaymentFrequency     string    `json:"payment_frequency"` // monthly, quarterly, annually
	Status               string    `json:"status"`            // active, lapsed, surrendered, expired
	IssueDate            string    `json:"issue_date"`        // DATE format
	ExpiryDate           string    `json:"expiry_date"`
	LastPremiumPaidDate  *string   `json:"last_premium_paid_date,omitempty"`
	NextPremiumDueDate   *string   `json:"next_premium_due_date,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	Product              *Product  `json:"product,omitempty"`
}

type Claim struct {
	ID                  string     `json:"id"`
	ClaimNumber         string     `json:"claim_number"`
	PolicyID            string     `json:"policy_id"`
	UserID              string     `json:"user_id"`
	ClaimType           string     `json:"claim_type"`
	ClaimAmount         int64      `json:"claim_amount"`
	IncidentDate        string     `json:"incident_date"` // DATE format
	IncidentDescription string     `json:"incident_description"`
	Status              string     `json:"status"` // submitted, under_review, approved, rejected, paid
	ReviewerID          *string    `json:"reviewer_id,omitempty"`
	ReviewerNotes       string     `json:"reviewer_notes,omitempty"`
	RejectionReason     string     `json:"rejection_reason,omitempty"`
	ApprovedAmount      *int64     `json:"approved_amount,omitempty"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	ReviewedAt          *time.Time `json:"reviewed_at,omitempty"`
	ApprovedAt          *time.Time `json:"approved_at,omitempty"`
	PaidAt              *time.Time `json:"paid_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	Policy              *Policy    `json:"policy,omitempty"`
}

type ClaimDocument struct {
	ID           string    `json:"id"`
	ClaimID      string    `json:"claim_id"`
	DocumentType string    `json:"document_type"` // medical_report, police_report, invoice, photo, other
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

type ClaimTimeline struct {
	ID          string                 `json:"id"`
	ClaimID     string                 `json:"claim_id"`
	Action      string                 `json:"action"` // submitted, document_uploaded, status_changed, comment_added
	Description string                 `json:"description"`
	ActorID     *string                `json:"actor_id,omitempty"`
	ActorName   string                 `json:"actor_name,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

type Invoice struct {
	ID               string     `json:"id"`
	InvoiceNumber    string     `json:"invoice_number"`
	PolicyID         string     `json:"policy_id"`
	UserID           string     `json:"user_id"`
	InvoiceType      string     `json:"invoice_type"` // premium, admin_fee, penalty, etc.
	Amount           int64      `json:"amount"`
	DueDate          string     `json:"due_date"` // DATE format
	Description      string     `json:"description,omitempty"`
	Status           string     `json:"status"`   // pending, paid, overdue, cancelled
	PaidAmount       int64      `json:"paid_amount"`
	PaidAt           *time.Time `json:"paid_at,omitempty"`
	PaymentMethod    string     `json:"payment_method,omitempty"`
	PaymentReference string     `json:"payment_reference,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Policy           *Policy    `json:"policy,omitempty"`
}

type Document struct {
	ID           string    `json:"id"`
	PolicyID     *string   `json:"policy_id,omitempty"`
	UserID       string    `json:"user_id"`
	DocumentType string    `json:"document_type"` // policy_certificate, endorsement, receipt, notice, other
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

type PolicyEndorsement struct {
	ID               string                 `json:"id"`
	EndorsementNumber string                 `json:"endorsement_number"`
	PolicyID         string                 `json:"policy_id"`
	EndorsementType  string                 `json:"endorsement_type"` // coverage_change, beneficiary_change, premium_adjustment
	Description      string                 `json:"description"`
	EffectiveDate    string                 `json:"effective_date"` // DATE format
	OldValues        map[string]interface{} `json:"old_values,omitempty"`
	NewValues        map[string]interface{} `json:"new_values,omitempty"`
	Status           string                 `json:"status"` // pending, approved, rejected
	ApprovedBy       *string                `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time             `json:"approved_at,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRole struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RoleID     string    `json:"role_id"`
	AssignedBy *string   `json:"assigned_by,omitempty"`
	AssignedAt time.Time `json:"assigned_at"`
}

type NotificationPreferences struct {
	ID                        string    `json:"id"`
	UserID                    string    `json:"user_id"`
	PromotionalEmails         bool      `json:"promotional_emails"`
	PolicyUpdateEmails        bool      `json:"policy_update_emails"`
	ClaimNotificationEmails   bool      `json:"claim_notification_emails"`
	NewsletterEmails          bool      `json:"newsletter_emails"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}
