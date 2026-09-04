package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
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
	return user, err
}
