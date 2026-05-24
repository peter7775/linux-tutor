package usecase

import "linux-tutor/internal/domain"

type RecommendNextInput struct {
	Current domain.Topic
	Weak    map[string]int
}
type RecommendNextOutput struct{ Next domain.Topic }
