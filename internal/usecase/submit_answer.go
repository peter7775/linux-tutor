package usecase

import "linux-tutor/internal/domain"

type SubmitAnswerInput struct {
	Question domain.Question
	Answer   domain.Answer
}
type SubmitAnswerOutput struct{ Result domain.AnswerResult }
