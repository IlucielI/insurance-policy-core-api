package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateSession(ctx context.Context, session *domain.ChatSession) error {
	session.ID = uuid.New().String()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	query := `
		INSERT INTO chat_sessions (id, user_id, session_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.SessionID, session.CreatedAt, session.UpdatedAt,
	)
	return err
}

func (r *ChatRepository) GetSessionBySessionID(ctx context.Context, sessionID string) (*domain.ChatSession, error) {
	session := &domain.ChatSession{}
	var userID sql.NullString

	query := `
		SELECT id, user_id, session_id, created_at, updated_at
		FROM chat_sessions
		WHERE session_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &userID, &session.SessionID, &session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if userID.Valid {
		session.UserID = &userID.String
	}

	return session, nil
}

func (r *ChatRepository) CreateMessage(ctx context.Context, message *domain.ChatMessage) error {
	message.ID = uuid.New().String()
	message.CreatedAt = time.Now()

	contextDocsJSON, _ := json.Marshal(message.ContextDocs)

	query := `
		INSERT INTO chat_messages (id, chat_session_id, role, content, context_docs, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		message.ID, message.ChatSessionID, message.Role, message.Content, contextDocsJSON, message.CreatedAt,
	)
	return err
}

func (r *ChatRepository) GetMessagesBySessionID(ctx context.Context, chatSessionID string, limit int) ([]*domain.ChatMessage, error) {
	query := `
		SELECT id, chat_session_id, role, content, context_docs, created_at
		FROM chat_messages
		WHERE chat_session_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, chatSessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*domain.ChatMessage{}
	for rows.Next() {
		message := &domain.ChatMessage{}
		var contextDocsJSON []byte

		err := rows.Scan(
			&message.ID, &message.ChatSessionID, &message.Role, &message.Content,
			&contextDocsJSON, &message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(contextDocsJSON, &message.ContextDocs)
		messages = append(messages, message)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *ChatRepository) SearchProductEmbeddings(ctx context.Context, embedding []float32, limit int) ([]*domain.ProductEmbedding, error) {
	// Vector similarity search using pgvector
	// This is a placeholder - actual implementation needs pgvector extension
	query := `
		SELECT id, product_id, chunk_type, chunk_text, created_at
		FROM product_embeddings
		ORDER BY embedding <=> $1
		LIMIT $2
	`
	
	// Convert []float32 to pgvector format (string representation)
	// For now, we'll skip the actual vector query and return empty
	// TODO: Implement proper pgvector query with embedding parameter
	
	rows, err := r.db.QueryContext(ctx, query, "[0,0,0]", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	embeddings := []*domain.ProductEmbedding{}
	for rows.Next() {
		emb := &domain.ProductEmbedding{}
		err := rows.Scan(&emb.ID, &emb.ProductID, &emb.ChunkType, &emb.ChunkText, &emb.CreatedAt)
		if err != nil {
			return nil, err
		}
		embeddings = append(embeddings, emb)
	}

	return embeddings, nil
}
