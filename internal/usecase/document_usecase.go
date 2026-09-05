package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type DocumentRepositoryInterface interface {
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Document, int, error)
	GetByID(ctx context.Context, id string) (*domain.Document, error)
	Create(ctx context.Context, doc *domain.Document) error
	Delete(ctx context.Context, id string) error
}

type StorageInterface interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	GetFileURL(ctx context.Context, objectName string, expiryDuration time.Duration) (string, error)
	DeleteFile(ctx context.Context, objectName string) error
}

type DocumentUsecase struct {
	documentRepo DocumentRepositoryInterface
	storage      StorageInterface
}

func NewDocumentUsecase(documentRepo DocumentRepositoryInterface, storage StorageInterface) *DocumentUsecase {
	return &DocumentUsecase{
		documentRepo: documentRepo,
		storage:      storage,
	}
}

func (u *DocumentUsecase) GetUserDocuments(ctx context.Context, userID string, limit, offset int) ([]*domain.Document, int, error) {
	return u.documentRepo.GetByUserID(ctx, userID, limit, offset)
}

func (u *DocumentUsecase) GetDocumentByID(ctx context.Context, id string) (*domain.Document, error) {
	doc, err := u.documentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, errors.New("document not found")
	}

	return doc, nil
}

// ValidateFile validates file type and size
func (u *DocumentUsecase) ValidateFile(file *multipart.FileHeader) error {
	// Check file size (max 10MB)
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return fmt.Errorf("file too large: max size is 10MB, got %d bytes", file.Size)
	}

	// Check file type
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
		"application/pdf": true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		// Also check by extension
		ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
		if ext != "jpg" && ext != "jpeg" && ext != "png" && ext != "pdf" {
			return fmt.Errorf("unsupported file type: %s (allowed: JPEG, PNG, PDF)", contentType)
		}
	}

	return nil
}

func (u *DocumentUsecase) UploadDocument(ctx context.Context, userID string, file *multipart.FileHeader, documentType, title, description string, policyID *string) (*domain.Document, error) {
	// Validate file
	if err := u.ValidateFile(file); err != nil {
		return nil, err
	}

	// Determine folder based on document type
	folder := "documents"
	switch documentType {
	case "policy_certificate", "endorsement", "receipt", "notice":
		folder = "policies"
	case "claim_photo", "medical_report", "police_report", "invoice":
		folder = "claims"
	case "id_card", "ktp", "sim":
		folder = "identities"
	}

	// Upload to MinIO
	objectName, err := u.storage.UploadFile(ctx, file, folder)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create document record
	doc := &domain.Document{
		UserID:       userID,
		PolicyID:     policyID,
		DocumentType: documentType,
		Title:        title,
		Description:  description,
		FileName:     file.Filename,
		FilePath:     objectName, // Store MinIO object name
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
	}

	if err := u.documentRepo.Create(ctx, doc); err != nil {
		// Rollback: delete uploaded file
		_ = u.storage.DeleteFile(ctx, objectName)
		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	return doc, nil
}

func (u *DocumentUsecase) GetDocumentDownloadURL(ctx context.Context, documentID string) (string, error) {
	doc, err := u.GetDocumentByID(ctx, documentID)
	if err != nil {
		return "", err
	}

	// Generate presigned URL (valid for 1 hour)
	url, err := u.storage.GetFileURL(ctx, doc.FilePath, time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate download URL: %w", err)
	}

	return url, nil
}

func (u *DocumentUsecase) DeleteDocument(ctx context.Context, id string) error {
	// Validate document exists
	doc, err := u.documentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if doc == nil {
		return errors.New("document not found")
	}

	// Delete from storage
	if err := u.storage.DeleteFile(ctx, doc.FilePath); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from database
	if err := u.documentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete document record: %w", err)
	}

	return nil
}
