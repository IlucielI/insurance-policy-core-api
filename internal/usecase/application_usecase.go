package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/model"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *domain.Application) error
	GetByID(ctx context.Context, id string) (*domain.Application, error)
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Application, int, error)
	ListWithFilters(ctx context.Context, search, userID, status, productID, priority, dateFrom, dateTo string, limit, offset int) ([]*domain.Application, int, error)
	Update(ctx context.Context, app *domain.Application) error
	UpdateStatus(ctx context.Context, id, status string, underwriterID *string, notes, rejectionReason string) error
}

type PolicyEmailService interface {
	SendPolicyIssuedEmail(to, fullName, policyNumber, productName string, sumAssured int64) error
	PreviewEmail(templateType string) (string, string, error)
}

type UserRepositoryInterface interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type NotificationService interface {
	Create(req *model.NotificationCreateRequest) (*model.Notification, error)
}

type ApplicationUsecase struct {
	appRepo       ApplicationRepository
	productRepo   interface {
		GetByID(ctx context.Context, id string) (*domain.Product, error)
	}
	emailService      PolicyEmailService
	notificationSvc   NotificationService
}

func NewApplicationUsecase(appRepo ApplicationRepository, productRepo interface {
	GetByID(ctx context.Context, id string) (*domain.Product, error)
}, emailService PolicyEmailService) *ApplicationUsecase {
	return &ApplicationUsecase{
		appRepo:      appRepo,
		productRepo:  productRepo,
		emailService: emailService,
	}
}

// SetNotificationService sets the notification service (optional injection)
func (u *ApplicationUsecase) SetNotificationService(svc NotificationService) {
	u.notificationSvc = svc
}

func (u *ApplicationUsecase) CreateApplication(ctx context.Context, app *domain.Application) error {
	product, err := u.productRepo.GetByID(ctx, app.ProductID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}
	if !product.IsActive {
		return errors.New("product is not active")
	}

	if app.SumAssured < product.MinSumAssured || app.SumAssured > product.MaxSumAssured {
		return errors.New("sum assured out of range")
	}

	if app.PaymentTerm < product.MinPaymentTerm || app.PaymentTerm > product.MaxPaymentTerm {
		return errors.New("payment term out of range")
	}

	app.Status = "draft"
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()

	return u.appRepo.Create(ctx, app)
}

func (u *ApplicationUsecase) GetApplicationByID(ctx context.Context, id string) (*domain.Application, error) {
	return u.appRepo.GetByID(ctx, id)
}

func (u *ApplicationUsecase) ListApplications(ctx context.Context, userID, status, productID string, limit, offset int) ([]*domain.Application, int, error) {
	filters := make(map[string]interface{})
	if userID != "" {
		filters["user_id"] = userID
	}
	if status != "" {
		filters["status"] = status
	}
	if productID != "" {
		filters["product_id"] = productID
	}

	return u.appRepo.List(ctx, filters, limit, offset)
}

func (u *ApplicationUsecase) ListApplicationsWithFilters(ctx context.Context, search, userID, status, productID, priority, dateFrom, dateTo string, limit, offset int) ([]*domain.Application, int, error) {
	return u.appRepo.ListWithFilters(ctx, search, userID, status, productID, priority, dateFrom, dateTo, limit, offset)
}

func (u *ApplicationUsecase) SubmitApplication(ctx context.Context, id string) error {
	app, err := u.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}
	if app.Status != "draft" {
		return errors.New("only draft applications can be submitted")
	}

	now := time.Now()
	app.Status = "submitted"
	app.SubmittedAt = &now
	app.UpdatedAt = now

	return u.appRepo.Update(ctx, app)
}

func (u *ApplicationUsecase) ReviewApplication(ctx context.Context, id, underwriterID string) error {
	app, err := u.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}
	if app.Status != "submitted" {
		return errors.New("only submitted applications can be reviewed")
	}

	now := time.Now()
	app.Status = "under_review"
	app.UnderwriterID = &underwriterID
	app.ReviewedAt = &now
	app.UpdatedAt = now

	return u.appRepo.Update(ctx, app)
}

func (u *ApplicationUsecase) ApproveApplication(ctx context.Context, id, underwriterID, notes string) error {
	app, err := u.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}
	if app.Status != "under_review" {
		return errors.New("only applications under review can be approved")
	}

	if err := u.appRepo.UpdateStatus(ctx, id, "approved", &underwriterID, notes, ""); err != nil {
		return err
	}

	// Send real-time notification
	if u.notificationSvc != nil {
		refID := id
		refType := "policy"
		go func() {
			u.notificationSvc.Create(&model.NotificationCreateRequest{
				UserID:        app.UserID,
				Type:          model.NotificationPolicyApproved,
				Title:         "Polis Disetujui",
				Message:       fmt.Sprintf("Pengajuan polis Anda telah disetujui."),
				ReferenceID:   &refID,
				ReferenceType: &refType,
			})
		}()
	}

	// Send email
	if u.emailService != nil {
		go func() {
			product, _ := u.productRepo.GetByID(context.Background(), app.ProductID)
			productName := "Insurance Policy"
			if product != nil {
				productName = product.Name
			}
			policyNumber := "POL-" + id[:8]
			_ = u.emailService.SendPolicyIssuedEmail("user@example.com", "Customer", policyNumber, productName, app.SumAssured)
		}()
	}

	return nil
}

func (u *ApplicationUsecase) RejectApplication(ctx context.Context, id, underwriterID, notes, rejectionReason string) error {
	app, err := u.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}
	if app.Status != "under_review" {
		return errors.New("only applications under review can be rejected")
	}
	if rejectionReason == "" {
		return errors.New("rejection reason is required")
	}

	return u.appRepo.UpdateStatus(ctx, id, "rejected", &underwriterID, notes, rejectionReason)
}

func (u *ApplicationUsecase) UpdateStatus(ctx context.Context, id, status string, underwriterID *string, notes, rejectionReason string) error {
	app, err := u.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}

	validStatuses := []string{"draft", "submitted", "under_review", "approved", "rejected"}
	isValid := false
	for _, s := range validStatuses {
		if s == status {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid status")
	}

	if status == "rejected" && rejectionReason == "" {
		return errors.New("rejection reason is required for rejected status")
	}

	err = u.appRepo.UpdateStatus(ctx, id, status, underwriterID, notes, rejectionReason)
	if err != nil {
		return err
	}

	// Send notification when status changes to approved
	if status == "approved" && u.notificationSvc != nil {
		refID := id
		refType := "policy"
		go func() {
			u.notificationSvc.Create(&model.NotificationCreateRequest{
				UserID:        app.UserID,
				Type:          model.NotificationPolicyApproved,
				Title:         "Polis Disetujui",
				Message:       fmt.Sprintf("Pengajuan polis Anda telah disetujui."),
				ReferenceID:   &refID,
				ReferenceType: &refType,
			})
		}()
	}

	return nil
}

func (u *ApplicationUsecase) SendPolicyIssuedEmailManual(ctx context.Context, email, fullName, policyNumber, productName string, sumAssured int64) error {
	if u.emailService == nil {
		return errors.New("email service not configured")
	}
	return u.emailService.SendPolicyIssuedEmail(email, fullName, policyNumber, productName, sumAssured)
}

func (u *ApplicationUsecase) PreviewEmail(templateType string) (string, string, error) {
	if u.emailService == nil {
		return "", "", errors.New("email service not configured")
	}
	return u.emailService.PreviewEmail(templateType)
}
