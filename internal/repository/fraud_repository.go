package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type FraudRepository struct {
	db *sql.DB
}

func NewFraudRepository(db *sql.DB) *FraudRepository {
	return &FraudRepository{db: db}
}

// UpdateRiskScore updates risk score and analysis for an application
func (r *FraudRepository) UpdateRiskScore(ctx context.Context, applicationID string, riskScore int, fraudFlags []string, analysisDetail string) error {
	fraudFlagsJSON, _ := json.Marshal(fraudFlags)
	
	query := `
		UPDATE applications
		SET risk_score = $1, 
		    fraud_flags = $2, 
		    risk_analysis_detail = $3,
		    risk_analyzed_at = $4,
		    updated_at = $5
		WHERE id = $6
	`
	
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, riskScore, fraudFlagsJSON, analysisDetail, now, now, applicationID)
	return err
}

// GetHighRiskApplications returns applications with high risk scores
func (r *FraudRepository) GetHighRiskApplications(ctx context.Context, minRiskScore int, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, user_id, product_id, risk_score, fraud_flags, status, created_at
		FROM applications
		WHERE risk_score >= $1
		ORDER BY risk_score DESC, created_at DESC
		LIMIT $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, minRiskScore, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	results := []map[string]interface{}{}
	for rows.Next() {
		var id, userID, productID, status string
		var riskScore int
		var fraudFlagsJSON []byte
		var createdAt time.Time
		
		if err := rows.Scan(&id, &userID, &productID, &riskScore, &fraudFlagsJSON, &status, &createdAt); err != nil {
			return nil, err
		}
		
		var fraudFlags []string
		json.Unmarshal(fraudFlagsJSON, &fraudFlags)
		
		results = append(results, map[string]interface{}{
			"id":          id,
			"user_id":     userID,
			"product_id":  productID,
			"risk_score":  riskScore,
			"fraud_flags": fraudFlags,
			"status":      status,
			"created_at":  createdAt,
		})
	}
	
	return results, nil
}
