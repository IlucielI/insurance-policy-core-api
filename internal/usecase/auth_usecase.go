package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
	CreatePasswordResetToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	GetPasswordResetToken(ctx context.Context, token string) (userID string, valid bool, err error)
	MarkTokenAsUsed(ctx context.Context, token string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int, error)
}

type RoleRepository interface {
	GetRoleNames(ctx context.Context, userID string) ([]string, error)
	AssignRole(ctx context.Context, userID, roleID, assignedBy string) error
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	RemoveRole(ctx context.Context, userID, roleID string) error
	GetAllRoles(ctx context.Context) ([]*domain.Role, error)
	GetRoleByID(ctx context.Context, id string) (*domain.Role, error)
}

type EmailService interface {
	SendWelcomeEmail(to, fullName string) error
	SendPasswordResetEmail(to, fullName, resetToken string) error
}

type AuthUsecase struct {
	userRepo     UserRepository
	roleRepo     RoleRepository
	emailService EmailService
}

func NewAuthUsecase(userRepo UserRepository, roleRepo RoleRepository, emailService EmailService) *AuthUsecase {
	return &AuthUsecase{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		emailService: emailService,
	}
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

	// Assign default customer role
	if u.roleRepo != nil {
		customerRole, err := u.roleRepo.GetRoleByName(ctx, "customer")
		if err == nil && customerRole != nil {
			_ = u.roleRepo.AssignRole(ctx, user.ID, customerRole.ID, user.ID)
		}
	}

	// Send welcome email (non-blocking - don't fail registration if email fails)
	if u.emailService != nil {
		go func() {
			_ = u.emailService.SendWelcomeEmail(email, fullName)
		}()
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*domain.User, []string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	// Get user roles
	var roles []string
	if u.roleRepo != nil {
		roles, err = u.roleRepo.GetRoleNames(ctx, user.ID)
		if err != nil {
			// Fallback to legacy role column if RBAC not set up
			roles = []string{user.Role}
		}
	} else {
		roles = []string{user.Role}
	}

	return user, roles, nil
}

func (u *AuthUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

// RequestPasswordReset generates reset token and sends email
func (u *AuthUsecase) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		// Don't reveal if email exists - security best practice
		return nil
	}

	// Generate secure token (32 random bytes = 64 hex chars)
	token := generateSecureToken()
	expiresAt := time.Now().Add(1 * time.Hour)

	// Save token to database
	if err := u.userRepo.CreatePasswordResetToken(ctx, user.ID, token, expiresAt); err != nil {
		return err
	}

	// Send reset email (non-blocking)
	if u.emailService != nil {
		go func() {
			_ = u.emailService.SendPasswordResetEmail(email, user.FullName, token)
		}()
	}

	return nil
}

// ResetPassword validates token and updates password
func (u *AuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate token
	userID, valid, err := u.userRepo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	if err := u.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return err
	}

	// Mark token as used
	if err := u.userRepo.MarkTokenAsUsed(ctx, token); err != nil {
		return err
	}

	return nil
}

// generateSecureToken creates cryptographically secure random token
func generateSecureToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to timestamp-based token (less secure but functional)
		return fmt.Sprintf("%d-%s", time.Now().UnixNano(), generateRandomString(32))
	}
	return fmt.Sprintf("%x", b)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

// Admin methods

// ListUsers returns all users with pagination
func (u *AuthUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	return u.userRepo.ListUsers(ctx, limit, offset)
}

// AssignRole assigns a role to a user
func (u *AuthUsecase) AssignRole(ctx context.Context, userID, roleID, assignedBy string) error {
	if u.roleRepo == nil {
		return errors.New("role repository not available")
	}
	return u.roleRepo.AssignRole(ctx, userID, roleID, assignedBy)
}

// RemoveRole removes a role from a user
func (u *AuthUsecase) RemoveRole(ctx context.Context, userID, roleID string) error {
	if u.roleRepo == nil {
		return errors.New("role repository not available")
	}
	return u.roleRepo.RemoveRole(ctx, userID, roleID)
}

// GetRoles returns all roles
func (u *AuthUsecase) GetRoles(ctx context.Context) ([]*domain.Role, error) {
	if u.roleRepo == nil {
		return nil, errors.New("role repository not available")
	}
	return u.roleRepo.GetAllRoles(ctx)
}
