package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/lib/pq"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// GetUserRoles retrieves all roles for a user
func (r *RoleRepository) GetUserRoles(ctx context.Context, userID string) ([]*domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.display_name, r.description, r.permissions
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		var permissionsJSON []byte

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.DisplayName,
			&role.Description,
			&permissionsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}

	return roles, rows.Err()
}

// GetRoleNames returns just the role names for a user (for JWT claims)
func (r *RoleRepository) GetRoleNames(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT r.name
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roleNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roleNames = append(roleNames, name)
	}

	return roleNames, rows.Err()
}

// AssignRole assigns a role to a user
func (r *RoleRepository) AssignRole(ctx context.Context, userID, roleID, assignedBy string) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, roleID, assignedBy)
	return err
}

// RemoveRole removes a role from a user
func (r *RoleRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	return err
}

// GetAllRoles returns all available roles in the system
func (r *RoleRepository) GetAllRoles(ctx context.Context) ([]*domain.Role, error) {
	query := `
		SELECT id, name, display_name, description, permissions
		FROM roles
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		var permissionsJSON []byte

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.DisplayName,
			&role.Description,
			&permissionsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}

	return roles, rows.Err()
}

// GetRoleByName retrieves a role by its name
func (r *RoleRepository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `
		SELECT id, name, display_name, description, permissions
		FROM roles
		WHERE name = $1
	`

	var role domain.Role
	var permissionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.DisplayName,
		&role.Description,
		&permissionsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
		return nil, err
	}

	return &role, nil
}

// GetRoleByID retrieves a role by its ID
func (r *RoleRepository) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	query := `
		SELECT id, name, display_name, description, permissions
		FROM roles
		WHERE id = $1
	`

	var role domain.Role
	var permissionsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.DisplayName,
		&role.Description,
		&permissionsJSON,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
		return nil, err
	}

	return &role, nil
}

// UserHasPermission checks if a user has a specific permission
func (r *RoleRepository) UserHasPermission(ctx context.Context, userID, permission string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1
			AND r.permissions @> $2::jsonb
		)
	`

	permJSON, _ := json.Marshal([]string{permission})
	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, permJSON).Scan(&exists)
	return exists, err
}

// UserHasRole checks if a user has a specific role
func (r *RoleRepository) UserHasRole(ctx context.Context, userID, roleName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM roles r
			INNER JOIN user_roles ur ON r.id = ur.role_id
			WHERE ur.user_id = $1 AND r.name = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, roleName).Scan(&exists)
	return exists, err
}

// GetUsersWithRole returns all users with a specific role
func (r *RoleRepository) GetUsersWithRole(ctx context.Context, roleName string) ([]string, error) {
	query := `
		SELECT ur.user_id
		FROM user_roles ur
		INNER JOIN roles r ON ur.role_id = r.id
		WHERE r.name = $1
	`

	rows, err := r.db.QueryContext(ctx, query, roleName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, rows.Err()
}

// BulkAssignRoles assigns multiple roles to a user
func (r *RoleRepository) BulkAssignRoles(ctx context.Context, userID string, roleNames []string, assignedBy string) error {
	// First, get role IDs
	query := `SELECT id FROM roles WHERE name = ANY($1)`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(roleNames))
	if err != nil {
		return err
	}
	defer rows.Close()

	var roleIDs []string
	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return err
		}
		roleIDs = append(roleIDs, roleID)
	}

	// Insert all role assignments
	insertQuery := `
		INSERT INTO user_roles (user_id, role_id, assigned_by)
		SELECT $1, unnest($2::uuid[]), $3
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err = r.db.ExecContext(ctx, insertQuery, userID, pq.Array(roleIDs), assignedBy)
	return err
}
