package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Document, int, error) {
	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM documents WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch documents
	query := `
		SELECT id, policy_id, user_id, document_type, title, description,
			file_name, file_path, file_size, mime_type, uploaded_at
		FROM documents
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	documents := []*domain.Document{}
	for rows.Next() {
		doc := &domain.Document{}
		err := rows.Scan(
			&doc.ID, &doc.PolicyID, &doc.UserID, &doc.DocumentType, &doc.Title,
			&doc.Description, &doc.FileName, &doc.FilePath, &doc.FileSize,
			&doc.MimeType, &doc.UploadedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		documents = append(documents, doc)
	}

	return documents, total, nil
}

func (r *DocumentRepository) GetByID(ctx context.Context, id string) (*domain.Document, error) {
	doc := &domain.Document{}
	query := `
		SELECT id, policy_id, user_id, document_type, title, description,
			file_name, file_path, file_size, mime_type, uploaded_at
		FROM documents
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&doc.ID, &doc.PolicyID, &doc.UserID, &doc.DocumentType, &doc.Title,
		&doc.Description, &doc.FileName, &doc.FilePath, &doc.FileSize,
		&doc.MimeType, &doc.UploadedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (r *DocumentRepository) Create(ctx context.Context, doc *domain.Document) error {
	doc.ID = uuid.New().String()
	doc.UploadedAt = time.Now()

	query := `
		INSERT INTO documents (id, policy_id, user_id, document_type, title, description,
			file_name, file_path, file_size, mime_type, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		doc.ID, doc.PolicyID, doc.UserID, doc.DocumentType, doc.Title,
		doc.Description, doc.FileName, doc.FilePath, doc.FileSize,
		doc.MimeType, doc.UploadedAt,
	)
	return err
}

func (r *DocumentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
