package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/cache"
	"github.com/google/uuid"
)

type UserRepository struct {
	db    *sql.DB
	cache *cache.RedisClient
}

func NewUserRepository(db *sql.DB, cache *cache.RedisClient) *UserRepository {
	return &UserRepository{db: db, cache: cache}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, password_hash, full_name, phone, role, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FullName, user.Phone,
		user.Role, user.IsVerified, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, password_hash, full_name, phone, role, is_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone,
		&user.Role, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	// Build cache key for user profile
	cacheKey := fmt.Sprintf("user:profile:%s", id)
	
	// Try cache first (15 min TTL for user profile/session data)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var user domain.User
			if json.Unmarshal([]byte(cached), &user) == nil {
				return &user, nil
			}
		}
	}
	
	user := &domain.User{}
	query := `
		SELECT id, email, password_hash, full_name, phone, role, is_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone,
		&user.Role, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Cache user profile for 15 minutes
	if r.cache != nil {
		if cached, err := json.Marshal(user); err == nil {
			_ = r.cache.Set(ctx, cacheKey, string(cached), 15*time.Minute)
		}
	}
	
	return user, err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, newPasswordHash, time.Now(), userID)
	
	// Invalidate user cache
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:profile:%s", userID)
		_ = r.cache.Delete(ctx, cacheKey)
	}
	
	return err
}

func (r *UserRepository) CreatePasswordResetToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, false, $4)
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt, time.Now())
	return err
}

func (r *UserRepository) GetPasswordResetToken(ctx context.Context, token string) (userID string, valid bool, err error) {
	query := `
		SELECT user_id, expires_at, used
		FROM password_reset_tokens
		WHERE token = $1
	`
	
	var expiresAt time.Time
	var used bool
	
	err = r.db.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAt, &used)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	
	// Check if token is valid (not expired and not used)
	if used || time.Now().After(expiresAt) {
		return userID, false, nil
	}
	
	return userID, true, nil
}

func (r *UserRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	query := `
		UPDATE password_reset_tokens
		SET used = true
		WHERE token = $1
	`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

// ListAllUsers returns all users without pagination for export
func (r *UserRepository) ListAllUsers(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, email, COALESCE(password_hash, '') as password_hash, full_name, COALESCE(phone, '') as phone, 
			role, is_verified, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		u := &domain.User{}
		err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone,
			&u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// ListUsers returns paginated users (for admin panel)
func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, email, COALESCE(password_hash, '') as password_hash, full_name, COALESCE(phone, '') as phone,
			role, is_verified, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		u := &domain.User{}
		err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone,
			&u.Role, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}
