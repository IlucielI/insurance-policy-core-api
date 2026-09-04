package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type ChatRepository interface {
	CreateSession(ctx context.Context, session *domain.ChatSession) error
	GetSessionBySessionID(ctx context.Context, sessionID string) (*domain.ChatSession, error)
	CreateMessage(ctx context.Context, message *domain.ChatMessage) error
	GetMessagesBySessionID(ctx context.Context, chatSessionID string, limit int) ([]*domain.ChatMessage, error)
	SearchProductEmbeddings(ctx context.Context, embedding []float32, limit int) ([]*domain.ProductEmbedding, error)
}

type LLMService interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	GenerateCompletion(ctx context.Context, messages []map[string]string) (string, error)
}

type ChatUsecase struct {
	chatRepo   ChatRepository
	llmService LLMService
}

func NewChatUsecase(chatRepo ChatRepository, llmService LLMService) *ChatUsecase {
	return &ChatUsecase{
		chatRepo:   chatRepo,
		llmService: llmService,
	}
}

func (u *ChatUsecase) GetOrCreateSession(ctx context.Context, sessionID string, userID *string) (*domain.ChatSession, error) {
	// Try to get existing session
	session, err := u.chatRepo.GetSessionBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	// Create new session
	session = &domain.ChatSession{
		SessionID: sessionID,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.chatRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (u *ChatUsecase) SendMessage(ctx context.Context, sessionID, userMessage string, userID *string) (string, []map[string]interface{}, error) {
	// Get or create session
	session, err := u.GetOrCreateSession(ctx, sessionID, userID)
	if err != nil {
		return "", nil, err
	}

	// Save user message
	userMsg := &domain.ChatMessage{
		ChatSessionID: session.ID,
		Role:          "user",
		Content:       userMessage,
		CreatedAt:     time.Now(),
	}
	if err := u.chatRepo.CreateMessage(ctx, userMsg); err != nil {
		return "", nil, err
	}

	// RAG: Generate embedding for user query
	var contextDocs []map[string]interface{}
	var relevantContext string

	// TODO: Implement when LLM endpoint is available
	// For now, use fallback rules-based responses
	assistantReply := u.generateFallbackResponse(userMessage)

	// If LLM service is available, use RAG
	if u.llmService != nil {
		// Generate embedding
		embedding, err := u.llmService.GenerateEmbedding(ctx, userMessage)
		if err == nil {
			// Search similar product embeddings
			embeddings, err := u.chatRepo.SearchProductEmbeddings(ctx, embedding, 3)
			if err == nil && len(embeddings) > 0 {
				// Build context from retrieved documents
				contextParts := []string{}
				for _, emb := range embeddings {
					contextParts = append(contextParts, emb.ChunkText)
					contextDocs = append(contextDocs, map[string]interface{}{
						"product_id": emb.ProductID,
						"chunk_type": emb.ChunkType,
						"text":       emb.ChunkText,
					})
				}
				relevantContext = strings.Join(contextParts, "\n\n")

				// Get recent chat history
				recentMessages, _ := u.chatRepo.GetMessagesBySessionID(ctx, session.ID, 5)
				
				// Build LLM prompt
				llmMessages := []map[string]string{}
				llmMessages = append(llmMessages, map[string]string{
					"role": "system",
					"content": fmt.Sprintf(`Anda adalah asisten AI untuk sistem asuransi. Gunakan konteks berikut untuk menjawab pertanyaan customer:

%s

Jawab dengan ramah, jelas, dan fokus pada informasi produk asuransi. Jika pertanyaan di luar konteks, arahkan ke customer service.`, relevantContext),
				})

				// Add recent history
				for _, msg := range recentMessages {
					llmMessages = append(llmMessages, map[string]string{
						"role":    msg.Role,
						"content": msg.Content,
					})
				}

				// Add current user message
				llmMessages = append(llmMessages, map[string]string{
					"role":    "user",
					"content": userMessage,
				})

				// Generate completion
				reply, err := u.llmService.GenerateCompletion(ctx, llmMessages)
				if err == nil {
					assistantReply = reply
				}
			}
		}
	}

	// Save assistant message
	assistantMsg := &domain.ChatMessage{
		ChatSessionID: session.ID,
		Role:          "assistant",
		Content:       assistantReply,
		ContextDocs:   contextDocs,
		CreatedAt:     time.Now(),
	}
	if err := u.chatRepo.CreateMessage(ctx, assistantMsg); err != nil {
		return "", nil, err
	}

	return assistantReply, contextDocs, nil
}

func (u *ChatUsecase) GetChatHistory(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error) {
	session, err := u.chatRepo.GetSessionBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return []*domain.ChatMessage{}, nil
	}

	return u.chatRepo.GetMessagesBySessionID(ctx, session.ID, limit)
}

// Fallback response when LLM is not available
func (u *ChatUsecase) generateFallbackResponse(message string) string {
	lower := strings.ToLower(message)

	if strings.Contains(lower, "produk") || strings.Contains(lower, "jenis") {
		return "Kami menyediakan 3 jenis produk asuransi:\n\n1. **Asuransi Jiwa** - Perlindungan finansial untuk keluarga Anda\n2. **Asuransi Kesehatan** - Biaya perawatan medis dan rawat inap\n3. **Asuransi Kendaraan** - Perlindungan mobil dan motor dari risiko\n\nSilakan pilih produk di halaman Products untuk informasi lebih detail!"
	}

	if strings.Contains(lower, "premi") || strings.Contains(lower, "harga") || strings.Contains(lower, "biaya") {
		return "Premi asuransi dihitung berdasarkan:\n- Usia Anda\n- Jumlah pertanggungan (sum assured)\n- Jangka waktu pembayaran\n\nGunakan kalkulator premi di halaman produk untuk simulasi langsung!"
	}

	if strings.Contains(lower, "klaim") || strings.Contains(lower, "claim") {
		return "Proses klaim asuransi:\n1. Hubungi customer service kami\n2. Siapkan dokumen pendukung (KTP, polis, bukti kejadian)\n3. Isi formulir klaim\n4. Tim kami akan review dalam 3-7 hari kerja\n\nUntuk bantuan lebih lanjut, hubungi support@insurance.com"
	}

	if strings.Contains(lower, "daftar") || strings.Contains(lower, "apply") || strings.Contains(lower, "beli") {
		return "Cara mengajukan asuransi:\n1. Pilih produk yang sesuai\n2. Hitung premi dengan kalkulator\n3. Klik 'Ajukan Sekarang'\n4. Isi formulir aplikasi\n5. Lakukan pembayaran\n6. Polis terbit dalam 1x24 jam\n\nMulai dari halaman Products!"
	}

	if strings.Contains(lower, "halo") || strings.Contains(lower, "hai") || strings.Contains(lower, "hello") {
		return "Halo! 👋 Selamat datang di Layanan Asuransi kami.\n\nSaya siap membantu Anda dengan:\n- Informasi produk asuransi\n- Perhitungan premi\n- Proses pendaftaran\n- Pertanyaan umum\n\nAda yang bisa saya bantu?"
	}

	return "Maaf, saya tidak sepenuhnya memahami pertanyaan Anda. Saya dapat membantu dengan:\n\n- Informasi produk asuransi (jiwa, kesehatan, kendaraan)\n- Cara menghitung premi\n- Proses klaim\n- Cara mendaftar\n\nSilakan ajukan pertanyaan Anda, atau hubungi customer service kami untuk bantuan lebih lanjut!"
}
