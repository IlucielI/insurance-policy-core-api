package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/payment"
	"github.com/IlucielI/insurance-policy-core-api/internal/model"
)

type BillingRepositoryInterface interface {
	GetInvoicesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error)
	GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error)
	GetInvoiceByOrderID(ctx context.Context, orderID string) (*domain.Invoice, error)
	PayInvoice(ctx context.Context, invoice *domain.Invoice) error
	GetPaymentHistory(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error)
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
	UpdatePaymentStatus(ctx context.Context, payment *domain.Payment) error
	GetAllInvoices(ctx context.Context, status string, limit, offset int) ([]*domain.Invoice, int, error)
	CreateInvoice(ctx context.Context, invoice *domain.Invoice) error
	UpdateInvoiceStatus(ctx context.Context, invoiceID, status string) error
}

type BillingUsecase struct {
	billingRepo      BillingRepositoryInterface
	policyRepo       PolicyRepositoryInterface
	userRepo         UserRepositoryInterface
	midtransClient   *payment.MidtransClient
	notificationSvc  NotificationService
}

func NewBillingUsecase(billingRepo BillingRepositoryInterface, policyRepo PolicyRepositoryInterface, userRepo UserRepositoryInterface, midtransClient *payment.MidtransClient) *BillingUsecase {
	return &BillingUsecase{
		billingRepo:    billingRepo,
		policyRepo:     policyRepo,
		userRepo:       userRepo,
		midtransClient: midtransClient,
	}
}

// SetNotificationService sets the notification service (optional injection)
func (u *BillingUsecase) SetNotificationService(svc NotificationService) {
	u.notificationSvc = svc
}

func (u *BillingUsecase) GetInvoices(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	invoices, total, err := u.billingRepo.GetInvoicesByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for _, invoice := range invoices {
		policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
		if err == nil && policy != nil {
			invoice.Policy = policy
		}
	}

	return invoices, total, nil
}

// GetInvoiceByOrderID gets invoice by order_id for payment verification (security: prevent URL manipulation)
func (u *BillingUsecase) GetInvoiceByOrderID(ctx context.Context, orderID string) (*domain.Invoice, error) {
	return u.billingRepo.GetInvoiceByOrderID(ctx, orderID)
}

func (u *BillingUsecase) CreatePayment(ctx context.Context, invoiceID string) (*payment.CreateTransactionResponse, error) {
	invoice, err := u.billingRepo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errors.New("invoice not found")
	}

	if invoice.Status == "paid" {
		return nil, errors.New("invoice already paid")
	}
	if invoice.Status == "cancelled" {
		return nil, errors.New("invoice is cancelled")
	}

	user, err := u.userRepo.GetByID(ctx, invoice.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w", err)
	}

	orderID := fmt.Sprintf("INV-%s-%d", invoice.InvoiceNumber, time.Now().Unix())

	newPayment := &domain.Payment{
		ID:            fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		ApplicationID: policy.ApplicationID,
		OrderID:       orderID,
		GrossAmount:   invoice.Amount,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := u.billingRepo.CreatePayment(ctx, newPayment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	itemName := fmt.Sprintf("Premium Payment - %s", policy.PolicyNumber)
	snapResp, err := u.midtransClient.CreateTransaction(payment.CreateTransactionRequest{
		OrderID:       orderID,
		GrossAmount:   invoice.Amount,
		CustomerName:  user.FullName,
		CustomerEmail: user.Email,
		CustomerPhone: user.Phone,
		ItemName:      itemName,
		ItemPrice:     invoice.Amount,
		ItemQuantity:  1,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	return snapResp, nil
}

func (u *BillingUsecase) PayInvoice(ctx context.Context, invoiceID, paymentMethod, paymentReference string) error {
	invoice, err := u.billingRepo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}

	if invoice.Status == "paid" {
		return errors.New("invoice already paid")
	}
	if invoice.Status == "cancelled" {
		return errors.New("invoice is cancelled")
	}

	invoice.PaidAmount = invoice.Amount
	invoice.PaymentMethod = paymentMethod
	invoice.PaymentReference = paymentReference

	if err := u.billingRepo.PayInvoice(ctx, invoice); err != nil {
		return err
	}

	policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
	if err == nil && policy != nil {
		if policy.NextPremiumDueDate != nil {
			now := time.Now().Format("2006-01-02")
			policy.LastPremiumPaidDate = &now
			u.policyRepo.Update(ctx, policy)
		}
	}

	// Send real-time notification
	if u.notificationSvc != nil {
		refID := invoice.ID
		refType := "invoice"
		go func() {
			u.notificationSvc.Create(&model.NotificationCreateRequest{
				UserID:        invoice.UserID,
				Type:          model.NotificationPaymentConfirmed,
				Title:         "Pembayaran Diterima",
				Message:       fmt.Sprintf("Pembayaran sebesar Rp%d telah kami terima.", invoice.Amount),
				ReferenceID:   &refID,
				ReferenceType: &refType,
			})
		}()
	}

	return nil
}

func (u *BillingUsecase) HandlePaymentWebhook(ctx context.Context, notification *payment.WebhookNotification) error {
	paymentRecord, err := u.billingRepo.GetPaymentByOrderID(ctx, notification.OrderID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if paymentRecord == nil {
		return errors.New("payment not found")
	}

	status, err := u.midtransClient.HandleWebhook(notification)
	if err != nil {
		return fmt.Errorf("failed to handle webhook: %w", err)
	}

	paymentRecord.MidtransTransactionID = status.TransactionID
	paymentRecord.PaymentType = status.PaymentType
	paymentRecord.UpdatedAt = time.Now()

	if payment.IsPaymentSuccess(status.TransactionStatus) {
		paymentRecord.Status = "success"
		now := time.Now()
		paymentRecord.PaidAt = &now

		if err := u.billingRepo.UpdatePaymentStatus(ctx, paymentRecord); err != nil {
			return fmt.Errorf("failed to update payment: %w", err)
		}

		return nil

	} else if payment.IsPaymentFailed(status.TransactionStatus) {
		paymentRecord.Status = "failed"
		if err := u.billingRepo.UpdatePaymentStatus(ctx, paymentRecord); err != nil {
			return fmt.Errorf("failed to update payment: %w", err)
		}
	}

	return nil
}

func (u *BillingUsecase) GetPaymentHistory(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	invoices, total, err := u.billingRepo.GetPaymentHistory(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for _, invoice := range invoices {
		policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
		if err == nil && policy != nil {
			invoice.Policy = policy
		}
	}

	return invoices, total, nil
}

func (u *BillingUsecase) GetInvoiceDetails(ctx context.Context, invoiceID string) (*domain.Invoice, error) {
	invoice, err := u.billingRepo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, fmt.Errorf("invoice not found")
	}

	policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
	if err == nil && policy != nil {
		invoice.Policy = policy
	}

	return invoice, nil
}

func (u *BillingUsecase) GetAllInvoices(ctx context.Context, status string, limit, offset int) ([]*domain.Invoice, int, error) {
	return u.billingRepo.GetAllInvoices(ctx, status, limit, offset)
}

func (u *BillingUsecase) CreateInvoice(ctx context.Context, policyID string, amount int64, dueDate, invoiceType, description string) (*domain.Invoice, error) {
	policy, err := u.policyRepo.GetByID(ctx, policyID)
	if err != nil || policy == nil {
		return nil, fmt.Errorf("policy not found")
	}

	invoice := &domain.Invoice{
		ID:          fmt.Sprintf("inv_%d", time.Now().UnixNano()),
		PolicyID:    policyID,
		UserID:      policy.UserID,
		Amount:      amount,
		DueDate:     dueDate,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.billingRepo.CreateInvoice(ctx, invoice); err != nil {
		return nil, err
	}
	return invoice, nil
}

func (u *BillingUsecase) UpdateInvoiceStatus(ctx context.Context, invoiceID, status string) error {
	return u.billingRepo.UpdateInvoiceStatus(ctx, invoiceID, status)
}
