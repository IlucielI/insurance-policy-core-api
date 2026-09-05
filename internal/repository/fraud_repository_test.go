package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFraudRepository_UpdateRiskScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewFraudRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		fraudFlags := []string{"high_value", "rapid_submission"}
		fraudFlagsJSON, _ := json.Marshal(fraudFlags)

		mock.ExpectExec("UPDATE applications SET risk_score").
			WithArgs(85, fraudFlagsJSON, "High risk detected", sqlmock.AnyArg(), sqlmock.AnyArg(), "app-123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateRiskScore(ctx, "app-123", 85, fraudFlags, "High risk detected")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		fraudFlags := []string{"test"}
		mock.ExpectExec("UPDATE applications SET risk_score").
			WillReturnError(assert.AnError)

		err := repo.UpdateRiskScore(ctx, "app-123", 50, fraudFlags, "Test")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestFraudRepository_GetHighRiskApplications(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewFraudRepository(db)
	ctx := context.Background()

	t.Run("success with results", func(t *testing.T) {
		now := time.Now()
		fraudFlags1, _ := json.Marshal([]string{"high_value", "rapid_submission"})
		fraudFlags2, _ := json.Marshal([]string{"suspicious_pattern"})

		rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "risk_score", "fraud_flags", "status", "created_at"}).
			AddRow("app-1", "user-1", "prod-1", 90, fraudFlags1, "pending_review", now).
			AddRow("app-2", "user-2", "prod-2", 75, fraudFlags2, "pending_review", now)

		mock.ExpectQuery("SELECT (.+) FROM applications WHERE risk_score").
			WithArgs(70, 10).
			WillReturnRows(rows)

		results, err := repo.GetHighRiskApplications(ctx, 70, 10)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "app-1", results[0]["id"])
		assert.Equal(t, 90, results[0]["risk_score"])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no results", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM applications WHERE risk_score").
			WithArgs(95, 10).
			WillReturnRows(sqlmock.NewRows([]string{}))

		results, err := repo.GetHighRiskApplications(ctx, 95, 10)
		assert.NoError(t, err)
		assert.Empty(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM applications WHERE risk_score").
			WithArgs(70, 10).
			WillReturnError(assert.AnError)

		results, err := repo.GetHighRiskApplications(ctx, 70, 10)
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
