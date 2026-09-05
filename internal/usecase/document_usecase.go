package usecase

import (
	"context"
	"errors"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type DocumentRepositoryInterface interface {
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Document, int, error)
	GetByID(ctx context.Context, id string) (*domain.Document, error)
	Create(ctx context.Context, doc *domain.Document) error
	Delete(ctx context.Context, id string) error
}

type DocumentUsecase struct {
	documentRepo DocumentRepositoryInterface
}

func NewDocumentUsecase(documentRepo DocumentRepositoryInterface) *DocumentUsecase {
	return &DocumentUsecase{
		documentRepo: documentRepo,
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

func (u *DocumentUsecase) UploadDocument(ctx context.Context, doc *domain.Document) error {
	// Basic validation
	if doc.UserID == "" {
		return errors.New("user_id is required")
	}
	if doc.Title == "" {
		return errors.New("title is required")
	}
	if doc.FileName == "" || doc.FilePath == "" {
		return errors.New("file name and path are required")
	}

	return u.documentRepo.Create(ctx, doc)
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

	return u.documentRepo.Delete(ctx, id)
}
