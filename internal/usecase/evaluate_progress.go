package usecase

import "linux-tutor/internal/domain"

type EvaluateProgressInput struct {
	Session domain.Session
	Result  domain.AnswerResult
}
type EvaluateProgressOutput struct{ Session domain.Session }
