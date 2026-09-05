package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type ClaimRepositoryInterface interface {
	Create(ctx context.Context, claim *domain.Claim) error
	GetByID(ctx context.Context, id string) (*domain.Claim, error)
	AddDocument(ctx context.Context, doc *domain.ClaimDocument) error
	AddTimelineEntry(ctx context.Context, entry *domain.ClaimTimeline) error
	GetTimeline(ctx context.Context, claimID string) ([]*domain.ClaimTimeline, error)
	Update(ctx context.Context, claim *domain.Claim) error
}

type ClaimUsecase struct {
	claimRepo  ClaimRepositoryInterface
	policyRepo PolicyRepositoryInterface
}

func NewClaimUsecase(claimRepo ClaimRepositoryInterface, policyRepo PolicyRepositoryInterface) *ClaimUsecase {
	return &ClaimUsecase{
		claimRepo:  claimRepo,
		policyRepo: policyRepo,
	}
}

func (u *ClaimUsecase) CreateClaim(ctx context.Context, claim *domain.Claim) error {
	// Validate policy exists and is active
	policy, err := u.policyRepo.GetByID(ctx, claim.PolicyID)
	if err != nil {
		return err
	}
	if policy == nil {
		return errors.New("policy not found")
	}
	if policy.Status != "active" {
		return errors.New("claims can only be filed for active policies")
	}

	// Validate claim amount
	if claim.ClaimAmount <= 0 {
		return errors.New("claim amount must be greater than 0")
	}
	if claim.ClaimAmount > policy.SumAssured {
		return errors.New("claim amount cannot exceed sum assured")
	}

	// Generate claim number
	claim.ClaimNumber = fmt.Sprintf("CLM-%s-%d", policy.PolicyNumber[:8], time.Now().Unix())
	claim.Status = "submitted"
	claim.UserID = policy.UserID

	// Create claim
	if err := u.claimRepo.Create(ctx, claim); err != nil {
		return err
	}

	// Add initial timeline entry
	timeline := &domain.ClaimTimeline{
		ClaimID:     claim.ID,
		Action:      "submitted",
		Description: fmt.Sprintf("Claim submitted for %s", claim.ClaimType),
		ActorName:   "Customer",
	}
	return u.claimRepo.AddTimelineEntry(ctx, timeline)
}

func (u *ClaimUsecase) GetClaimByID(ctx context.Context, id string) (*domain.Claim, error) {
	claim, err := u.claimRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, errors.New("claim not found")
	}

	// Enrich with policy details
	policy, err := u.policyRepo.GetByID(ctx, claim.PolicyID)
	if err == nil && policy != nil {
		claim.Policy = policy
	}

	return claim, nil
}

func (u *ClaimUsecase) UploadDocument(ctx context.Context, claimID, documentType, fileName, filePath, mimeType string, fileSize int64) error {
	// Validate claim exists
	claim, err := u.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return err
	}
	if claim == nil {
		return errors.New("claim not found")
	}

	// Create document record
	doc := &domain.ClaimDocument{
		ClaimID:      claimID,
		DocumentType: documentType,
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     fileSize,
		MimeType:     mimeType,
	}

	if err := u.claimRepo.AddDocument(ctx, doc); err != nil {
		return err
	}

	// Add timeline entry
	timeline := &domain.ClaimTimeline{
		ClaimID:     claimID,
		Action:      "document_uploaded",
		Description: fmt.Sprintf("Document uploaded: %s (%s)", fileName, documentType),
		ActorName:   "Customer",
	}
	return u.claimRepo.AddTimelineEntry(ctx, timeline)
}

func (u *ClaimUsecase) GetClaimTimeline(ctx context.Context, claimID string) ([]*domain.ClaimTimeline, error) {
	// Validate claim exists
	claim, err := u.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, errors.New("claim not found")
	}

	return u.claimRepo.GetTimeline(ctx, claimID)
}
