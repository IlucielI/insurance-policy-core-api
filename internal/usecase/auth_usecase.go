package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type AuthUsecase struct {
	userRepo UserRepository
}

func NewAuthUsecase(userRepo UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo}
}

func (u *AuthUsecase) Register(ctx context.Context, email, password, fullName, phone string) (*domain.User, error) {
	// Check existing
	existing, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Phone:        phone,
		Role:         "customer",
		IsVerified:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (u *AuthUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}
