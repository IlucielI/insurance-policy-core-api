package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *domain.Application) error
	GetByID(ctx context.Context, id string) (*domain.Application, error)
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Application, int, error)
	Update(ctx context.Context, app *domain.Application) error
	UpdateStatus(ctx context.Context, id, status string, underwriterID *string, notes, rejectionReason string) error
}

type ApplicationUsecase struct {
	appRepo     ApplicationRepository
	productRepo interface {
		GetByID(ctx context.Context, id string) (*domain.Product, error)
	}
}

func NewApplicationUsecase(appRepo ApplicationRepository, productRepo interface {
	GetByID(ctx context.Context, id string) (*domain.Product, error)
}) *ApplicationUsecase {
	return &ApplicationUsecase{
		appRepo:     appRepo,
		productRepo: productRepo,
	}
}

func (u *ApplicationUsecase) CreateApplication(ctx context.Context, app *domain.Application) error {
	// Validate product exists
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

	// Validate sum assured range
	if app.SumAssured < product.MinSumAssured || app.SumAssured > product.MaxSumAssured {
		return errors.New("sum assured out of range")
	}

	// Validate payment term
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

	return u.appRepo.UpdateStatus(ctx, id, "approved", &underwriterID, notes, "")
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
