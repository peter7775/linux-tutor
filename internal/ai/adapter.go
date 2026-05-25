package ai

import "linux-tutor/internal/domain"

type Adapter interface {
	SuggestNext(domain.Topic, map[string]int) domain.Topic
	Explain(domain.Question, string) string
}
