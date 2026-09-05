package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/model"
)

type ClaimRepositoryInterface interface {
	Create(ctx context.Context, claim *domain.Claim) error
	GetByID(ctx context.Context, id string) (*domain.Claim, error)
	AddDocument(ctx context.Context, doc *domain.ClaimDocument) error
	AddTimelineEntry(ctx context.Context, entry *domain.ClaimTimeline) error
	GetTimeline(ctx context.Context, claimID string) ([]*domain.ClaimTimeline, error)
	Update(ctx context.Context, claim *domain.Claim) error
	ListClaimsWithFilters(ctx context.Context, search, status, claimType, dateFrom, dateTo, amountMin, amountMax string, limit, offset int) ([]*domain.Claim, int, error)
	ListAllClaimsForExport(ctx context.Context, status, claimType string) ([]*domain.Claim, error)
}

type ClaimEmailService interface {
	SendClaimStatusUpdateEmail(to, fullName, claimNumber, status, notes string) error
}

type ClaimUsecase struct {
	claimRepo        ClaimRepositoryInterface
	policyRepo       PolicyRepositoryInterface
	emailService     ClaimEmailService
	notificationSvc  NotificationService
}

func NewClaimUsecase(claimRepo ClaimRepositoryInterface, policyRepo PolicyRepositoryInterface, emailService ClaimEmailService) *ClaimUsecase {
	return &ClaimUsecase{
		claimRepo:    claimRepo,
		policyRepo:   policyRepo,
		emailService: emailService,
	}
}

// SetNotificationService sets the notification service (optional injection)
func (u *ClaimUsecase) SetNotificationService(svc NotificationService) {
	u.notificationSvc = svc
}

func (u *ClaimUsecase) CreateClaim(ctx context.Context, claim *domain.Claim) error {
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

	if claim.ClaimAmount <= 0 {
		return errors.New("claim amount must be greater than 0")
	}
	if claim.ClaimAmount > policy.SumAssured {
		return errors.New("claim amount cannot exceed sum assured")
	}

	claim.ClaimNumber = fmt.Sprintf("CLM-%s-%d", policy.PolicyNumber[:8], time.Now().Unix())
	claim.Status = "submitted"
	claim.UserID = policy.UserID

	if err := u.claimRepo.Create(ctx, claim); err != nil {
		return err
	}

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

	policy, err := u.policyRepo.GetByID(ctx, claim.PolicyID)
	if err == nil && policy != nil {
		claim.Policy = policy
	}

	return claim, nil
}

func (u *ClaimUsecase) UploadDocument(ctx context.Context, claimID, documentType, fileName, filePath, mimeType string, fileSize int64) error {
	claim, err := u.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return err
	}
	if claim == nil {
		return errors.New("claim not found")
	}

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

	timeline := &domain.ClaimTimeline{
		ClaimID:     claimID,
		Action:      "document_uploaded",
		Description: fmt.Sprintf("Document uploaded: %s (%s)", fileName, documentType),
		ActorName:   "Customer",
	}
	return u.claimRepo.AddTimelineEntry(ctx, timeline)
}

func (u *ClaimUsecase) GetClaimTimeline(ctx context.Context, claimID string) ([]*domain.ClaimTimeline, error) {
	claim, err := u.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return nil, err
	}
	if claim == nil {
		return nil, errors.New("claim not found")
	}

	return u.claimRepo.GetTimeline(ctx, claimID)
}

func (u *ClaimUsecase) UpdateClaimStatus(ctx context.Context, claimID, newStatus, notes string) error {
	claim, err := u.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return err
	}
	if claim == nil {
		return errors.New("claim not found")
	}

	claim.Status = newStatus
	claim.UpdatedAt = time.Now()
	if err := u.claimRepo.Update(ctx, claim); err != nil {
		return err
	}

	timeline := &domain.ClaimTimeline{
		ClaimID:     claimID,
		Action:      "status_updated",
		Description: fmt.Sprintf("Claim status updated to: %s", newStatus),
		ActorName:   "System",
	}
	if err := u.claimRepo.AddTimelineEntry(ctx, timeline); err != nil {
		return err
	}

	// Send real-time notification
	if u.notificationSvc != nil {
		refID := claimID
		refType := "claim"
		go func() {
			u.notificationSvc.Create(&model.NotificationCreateRequest{
				UserID:        claim.UserID,
				Type:          model.NotificationClaimUpdated,
				Title:         "Klaim Diperbarui",
				Message:       fmt.Sprintf("Status klaim #%s berubah menjadi: %s", claim.ClaimNumber, newStatus),
				ReferenceID:   &refID,
				ReferenceType: &refType,
			})
		}()
	}

	// Send email
	if u.emailService != nil {
		go func() {
			_ = u.emailService.SendClaimStatusUpdateEmail("user@example.com", "Customer", claim.ClaimNumber, newStatus, notes)
		}()
	}

	return nil
}

func (u *ClaimUsecase) SendClaimStatusEmailManual(ctx context.Context, email, fullName, claimNumber, status, notes string) error {
	if u.emailService == nil {
		return errors.New("email service not configured")
	}
	return u.emailService.SendClaimStatusUpdateEmail(email, fullName, claimNumber, status, notes)
}

// ListClaimsWithFilters returns claims with search and filters (admin).
// The priority parameter from handler is accepted but not used in current DB schema.
func (u *ClaimUsecase) ListClaimsWithFilters(ctx context.Context, search, status, claimType, priority, dateFrom, dateTo, amountMin, amountMax string, limit, offset int) ([]*domain.Claim, int, error) {
	return u.claimRepo.ListClaimsWithFilters(ctx, search, status, claimType, dateFrom, dateTo, amountMin, amountMax, limit, offset)
}

// ListAllClaimsForExport returns all claims for report export
func (u *ClaimUsecase) ListAllClaimsForExport(ctx context.Context, status, claimType string) ([]*domain.Claim, error) {
	return u.claimRepo.ListAllClaimsForExport(ctx, status, claimType)
}
