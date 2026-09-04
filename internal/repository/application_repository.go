package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type ApplicationRepository struct {
	db *sql.DB
}

func NewApplicationRepository(db *sql.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(ctx context.Context, app *domain.Application) error {
	app.ID = uuid.New().String()
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	app.Status = "draft"

	applicantDataJSON, _ := json.Marshal(app.ApplicantData)
	healthQuestionsJSON, _ := json.Marshal(app.HealthQuestions)

	query := `
		INSERT INTO applications (id, user_id, product_id, applicant_data, sum_assured, 
			payment_term, premium_amount, health_questions, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		app.ID, app.UserID, app.ProductID, applicantDataJSON, app.SumAssured,
		app.PaymentTerm, app.PremiumAmount, healthQuestionsJSON, app.Status,
		app.CreatedAt, app.UpdatedAt,
	)
	return err
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id string) (*domain.Application, error) {
	app := &domain.Application{}
	var applicantDataJSON, healthQuestionsJSON []byte
	var underwriterID, underwriterNotes, rejectionReason sql.NullString
	var submittedAt, reviewedAt, approvedAt sql.NullTime

	query := `
		SELECT id, user_id, product_id, applicant_data, sum_assured, payment_term, 
			premium_amount, health_questions, status, underwriter_id, underwriter_notes, 
			rejection_reason, submitted_at, reviewed_at, approved_at, created_at, updated_at
		FROM applications
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&app.ID, &app.UserID, &app.ProductID, &applicantDataJSON, &app.SumAssured,
		&app.PaymentTerm, &app.PremiumAmount, &healthQuestionsJSON, &app.Status,
		&underwriterID, &underwriterNotes, &rejectionReason,
		&submittedAt, &reviewedAt, &approvedAt, &app.CreatedAt, &app.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(applicantDataJSON, &app.ApplicantData)
	json.Unmarshal(healthQuestionsJSON, &app.HealthQuestions)

	if underwriterID.Valid {
		app.UnderwriterID = &underwriterID.String
	}
	if underwriterNotes.Valid {
		app.UnderwriterNotes = underwriterNotes.String
	}
	if rejectionReason.Valid {
		app.RejectionReason = rejectionReason.String
	}
	if submittedAt.Valid {
		app.SubmittedAt = &submittedAt.Time
	}
	if reviewedAt.Valid {
		app.ReviewedAt = &reviewedAt.Time
	}
	if approvedAt.Valid {
		app.ApprovedAt = &approvedAt.Time
	}

	return app, nil
}

func (r *ApplicationRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Application, int, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, userID)
		argPos++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	if productID, ok := filters["product_id"].(string); ok && productID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("product_id = $%d", argPos))
		args = append(args, productID)
		argPos++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM applications %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch applications
	query := fmt.Sprintf(`
		SELECT id, user_id, product_id, applicant_data, sum_assured, payment_term, 
			premium_amount, health_questions, status, underwriter_id, underwriter_notes, 
			rejection_reason, submitted_at, reviewed_at, approved_at, created_at, updated_at
		FROM applications
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	applications := []*domain.Application{}
	for rows.Next() {
		app := &domain.Application{}
		var applicantDataJSON, healthQuestionsJSON []byte
		var underwriterID, underwriterNotes, rejectionReason sql.NullString
		var submittedAt, reviewedAt, approvedAt sql.NullTime

		err := rows.Scan(
			&app.ID, &app.UserID, &app.ProductID, &applicantDataJSON, &app.SumAssured,
			&app.PaymentTerm, &app.PremiumAmount, &healthQuestionsJSON, &app.Status,
			&underwriterID, &underwriterNotes, &rejectionReason,
			&submittedAt, &reviewedAt, &approvedAt, &app.CreatedAt, &app.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		json.Unmarshal(applicantDataJSON, &app.ApplicantData)
		json.Unmarshal(healthQuestionsJSON, &app.HealthQuestions)

		if underwriterID.Valid {
			app.UnderwriterID = &underwriterID.String
		}
		if underwriterNotes.Valid {
			app.UnderwriterNotes = underwriterNotes.String
		}
		if rejectionReason.Valid {
			app.RejectionReason = rejectionReason.String
		}
		if submittedAt.Valid {
			app.SubmittedAt = &submittedAt.Time
		}
		if reviewedAt.Valid {
			app.ReviewedAt = &reviewedAt.Time
		}
		if approvedAt.Valid {
			app.ApprovedAt = &approvedAt.Time
		}

		applications = append(applications, app)
	}

	return applications, total, nil
}

func (r *ApplicationRepository) Update(ctx context.Context, app *domain.Application) error {
	app.UpdatedAt = time.Now()

	applicantDataJSON, _ := json.Marshal(app.ApplicantData)
	healthQuestionsJSON, _ := json.Marshal(app.HealthQuestions)

	query := `
		UPDATE applications
		SET user_id = $1, product_id = $2, applicant_data = $3, sum_assured = $4, 
			payment_term = $5, premium_amount = $6, health_questions = $7, status = $8, 
			underwriter_id = $9, underwriter_notes = $10, rejection_reason = $11, 
			submitted_at = $12, reviewed_at = $13, approved_at = $14, updated_at = $15
		WHERE id = $16
	`
	_, err := r.db.ExecContext(ctx, query,
		app.UserID, app.ProductID, applicantDataJSON, app.SumAssured,
		app.PaymentTerm, app.PremiumAmount, healthQuestionsJSON, app.Status,
		app.UnderwriterID, app.UnderwriterNotes, app.RejectionReason,
		app.SubmittedAt, app.ReviewedAt, app.ApprovedAt, app.UpdatedAt, app.ID,
	)
	return err
}

func (r *ApplicationRepository) UpdateStatus(ctx context.Context, id, status string, underwriterID *string, notes, rejectionReason string) error {
	now := time.Now()

	query := `
		UPDATE applications
		SET status = $1, underwriter_id = $2, underwriter_notes = $3, 
			rejection_reason = $4, reviewed_at = $5, 
			approved_at = CASE WHEN $1 = 'approved' THEN $6 ELSE approved_at END,
			updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.ExecContext(ctx, query,
		status, underwriterID, notes, rejectionReason, now, now, now, id,
	)
	return err
}
