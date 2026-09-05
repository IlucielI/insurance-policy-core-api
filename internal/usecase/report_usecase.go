package usecase

import (
	"context"
	"fmt"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/report"
)

// ReportRepositoryInterface defines data access for report generation
type ReportRepositoryInterface interface {
	GetAnalyticsSummary(ctx context.Context) (*report.AnalyticsSummary, error)
}

// ReportUserRepoInterface is subset of user repository needed for reports
type ReportUserRepoInterface interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	ListAllUsers(ctx context.Context) ([]*domain.User, error)
}

// ReportBillingRepoInterface is subset of billing repository needed for reports
type ReportBillingRepoInterface interface {
	GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error)
}

// ReportPolicyRepoInterface is subset of policy repository needed for reports
type ReportPolicyRepoInterface interface {
	GetByID(ctx context.Context, id string) (*domain.Policy, error)
}

// ReportClaimRepoInterface is subset of claim repository needed for reports
type ReportClaimRepoInterface interface {
	ListAllClaimsForExport(ctx context.Context, status, claimType string) ([]*domain.Claim, error)
}

type ReportUsecase struct {
	reportRepo  ReportRepositoryInterface
	userRepo    ReportUserRepoInterface
	billingRepo ReportBillingRepoInterface
	policyRepo  ReportPolicyRepoInterface
	claimRepo   ReportClaimRepoInterface
	pdfGen      *report.PDFGenerator
	excelGen    *report.ExcelGenerator
}

func NewReportUsecase(
	reportRepo ReportRepositoryInterface,
	userRepo ReportUserRepoInterface,
	billingRepo ReportBillingRepoInterface,
	policyRepo ReportPolicyRepoInterface,
	claimRepo ReportClaimRepoInterface,
) *ReportUsecase {
	return &ReportUsecase{
		reportRepo:  reportRepo,
		userRepo:    userRepo,
		billingRepo: billingRepo,
		policyRepo:  policyRepo,
		claimRepo:   claimRepo,
		pdfGen:      report.NewPDFGenerator(),
		excelGen:    report.NewExcelGenerator(),
	}
}

// GenerateBillingStatementPDF returns PDF bytes for a billing statement
func (u *ReportUsecase) GenerateBillingStatementPDF(ctx context.Context, invoiceID string) ([]byte, error) {
	invoice, err := u.billingRepo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("get invoice: %w", err)
	}
	if invoice == nil {
		return nil, fmt.Errorf("invoice not found")
	}

	user, err := u.userRepo.GetByID(ctx, invoice.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	policy, err := u.policyRepo.GetByID(ctx, invoice.PolicyID)
	if err != nil {
		return nil, fmt.Errorf("get policy: %w", err)
	}

	return u.pdfGen.BillingStatementPDF(invoice, user, policy)
}

// GenerateClaimsReportExcel returns Excel bytes for claims report
func (u *ReportUsecase) GenerateClaimsReportExcel(ctx context.Context, status, claimType string) ([]byte, error) {
	claims, err := u.claimRepo.ListAllClaimsForExport(ctx, status, claimType)
	if err != nil {
		return nil, fmt.Errorf("list claims: %w", err)
	}

	rows := make([]report.ClaimsReportRow, len(claims))
	for i, c := range claims {
		policyNum := ""
		if c.Policy != nil {
			policyNum = c.Policy.PolicyNumber
		}
		rows[i] = report.ClaimsReportRow{
			No:           i + 1,
			ClaimNumber:  c.ClaimNumber,
			PolicyNum:    policyNum,
			ClaimType:    c.ClaimType,
			ClaimAmount:  c.ClaimAmount,
			Status:       c.Status,
			IncidentDate: c.IncidentDate,
			SubmittedAt:  c.SubmittedAt.Format("02 Jan 2006"),
		}
	}

	filters := map[string]string{}
	if status != "" {
		filters["status"] = status
	}
	if claimType != "" {
		filters["claim_type"] = claimType
	}

	return u.excelGen.ClaimsReportExcel(rows, filters)
}

// GenerateCustomerListExcel returns Excel bytes for customer list
func (u *ReportUsecase) GenerateCustomerListExcel(ctx context.Context) ([]byte, error) {
	users, err := u.userRepo.ListAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	rows := make([]report.CustomerRow, len(users))
	for i, user := range users {
		rows[i] = report.CustomerRow{
			No:        i + 1,
			FullName:  user.FullName,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("02 Jan 2006"),
		}
	}

	return u.excelGen.CustomerListExcel(rows)
}

// GenerateAnalyticsSummaryExcel returns Excel bytes for analytics summary
func (u *ReportUsecase) GenerateAnalyticsSummaryExcel(ctx context.Context) ([]byte, error) {
	summary, err := u.reportRepo.GetAnalyticsSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("get analytics: %w", err)
	}

	return u.excelGen.AnalyticsSummaryExcel(summary)
}
