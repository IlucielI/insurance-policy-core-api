package service

import (
	"context"
	"fmt"
)

// MockLLMService is a temporary mock implementation
type MockLLMService struct{}

func NewMockLLMService() *MockLLMService {
	return &MockLLMService{}
}

func (s *MockLLMService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Return mock embedding (normally would call LLM API)
	return []float32{0.1, 0.2, 0.3}, nil
}

func (s *MockLLMService) GenerateCompletion(ctx context.Context, messages []map[string]string) (string, error) {
	// Return mock completion (normally would call LLM API)
	lastMsg := ""
	for _, msg := range messages {
		if msg["role"] == "user" {
			lastMsg = msg["content"]
		}
	}
	
	return fmt.Sprintf("Terima kasih atas pertanyaan Anda tentang '%s'. Saat ini sistem sedang dalam tahap pengembangan. Untuk informasi lebih lanjut tentang produk asuransi kami, silakan hubungi customer service atau lihat halaman produk.", lastMsg), nil
}
