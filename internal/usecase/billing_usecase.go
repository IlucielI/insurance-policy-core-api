package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type BillingRepositoryInterface interface {
	GetInvoicesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error)
	GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error)
	PayInvoice(ctx context.Context, invoice *domain.Invoice) error
	GetPaymentHistory(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error)
}

type BillingUsecase struct {
	billingRepo BillingRepositoryInterface
	policyRepo  PolicyRepositoryInterface
}

func NewBillingUsecase(billingRepo BillingRepositoryInterface, policyRepo PolicyRepositoryInterface) *BillingUsecase {
	return &BillingUsecase{
		billingRepo: billingRepo,
		policyRepo:  policyRepo,
	}
}

func (u *BillingUsecase) GetInvoices(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	invoices, total, err := u.billingRepo.GetInvoicesByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Enrich with policy details
	for _, invoice := range invoices {
		policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
		if err == nil && policy != nil {
			invoice.Policy = policy
		}
	}

	return invoices, total, nil
}

func (u *BillingUsecase) PayInvoice(ctx context.Context, invoiceID, paymentMethod, paymentReference string) error {
	// Get invoice
	invoice, err := u.billingRepo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return err
	}
	if invoice == nil {
		return errors.New("invoice not found")
	}

	// Validate invoice status
	if invoice.Status == "paid" {
		return errors.New("invoice already paid")
	}
	if invoice.Status == "cancelled" {
		return errors.New("invoice is cancelled")
	}

	// Update invoice
	invoice.PaidAmount = invoice.Amount
	invoice.PaymentMethod = paymentMethod
	invoice.PaymentReference = paymentReference

	// Process payment
	if err := u.billingRepo.PayInvoice(ctx, invoice); err != nil {
		return err
	}

	// Update policy next premium due date
	policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
	if err == nil && policy != nil {
		// Simple logic: add 1 month to next due date
		if policy.NextPremiumDueDate != nil {
			// Parse and add 1 month
			// In production, use proper date handling
			now := time.Now().Format("2006-01-02")
			policy.LastPremiumPaidDate = &now
			u.policyRepo.Update(ctx, policy)
		}
	}

	return nil
}

func (u *BillingUsecase) GetPaymentHistory(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	invoices, total, err := u.billingRepo.GetPaymentHistory(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Enrich with policy details
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

	// Enrich with policy
	policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
	if err == nil && policy != nil {
		invoice.Policy = policy
	}

	return invoice, nil
}
