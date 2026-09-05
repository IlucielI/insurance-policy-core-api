package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type PolicyRepositoryInterface interface {
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Policy, int, error)
	GetByID(ctx context.Context, id string) (*domain.Policy, error)
	CreateEndorsement(ctx context.Context, endorsement *domain.PolicyEndorsement) error
	Update(ctx context.Context, policy *domain.Policy) error
}

type PolicyUsecase struct {
	policyRepo  PolicyRepositoryInterface
	productRepo ProductRepository
}

func NewPolicyUsecase(policyRepo PolicyRepositoryInterface, productRepo ProductRepository) *PolicyUsecase {
	return &PolicyUsecase{
		policyRepo:  policyRepo,
		productRepo: productRepo,
	}
}

func (u *PolicyUsecase) GetUserPolicies(ctx context.Context, userID string, limit, offset int) ([]*domain.Policy, int, error) {
	policies, total, err := u.policyRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Enrich with product details
	for _, policy := range policies {
		product, err := u.productRepo.GetByID(ctx, policy.ProductID)
		if err == nil && product != nil {
			policy.Product = product
		}
	}

	return policies, total, nil
}

func (u *PolicyUsecase) GetPolicyByID(ctx context.Context, id string) (*domain.Policy, error) {
	policy, err := u.policyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, errors.New("policy not found")
	}

	// Enrich with product details
	product, err := u.productRepo.GetByID(ctx, policy.ProductID)
	if err == nil && product != nil {
		policy.Product = product
	}

	return policy, nil
}

func (u *PolicyUsecase) EndorsePolicy(ctx context.Context, policyID string, endorsementType, description, effectiveDate string) error {
	// Validate policy exists
	policy, err := u.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return err
	}
	if policy == nil {
		return errors.New("policy not found")
	}

	// Check policy is active
	if policy.Status != "active" {
		return errors.New("only active policies can be endorsed")
	}

	// Generate endorsement number
	endorsementNumber := fmt.Sprintf("END-%s-%d", policyID[:8], time.Now().Unix())

	endorsement := &domain.PolicyEndorsement{
		EndorsementNumber: endorsementNumber,
		PolicyID:          policyID,
		EndorsementType:   endorsementType,
		Description:       description,
		EffectiveDate:     effectiveDate,
		Status:            "pending",
	}

	return u.policyRepo.CreateEndorsement(ctx, endorsement)
}

func (u *PolicyUsecase) RenewPolicy(ctx context.Context, policyID string) error {
	// Validate policy exists
	policy, err := u.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return err
	}
	if policy == nil {
		return errors.New("policy not found")
	}

	// Check policy status
	if policy.Status == "active" {
		return errors.New("policy is already active")
	}

	// Simple renewal: reactivate and extend expiry
	policy.Status = "active"
	// In real scenario, would extend expiry_date by term period

	return u.policyRepo.Update(ctx, policy)
}
