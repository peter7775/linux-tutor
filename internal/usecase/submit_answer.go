package usecase

import "linux-tutor/internal/domain"

func SubmitAnswer(q domain.Question, a domain.Answer) domain.AnswerResult {
	_ = a
	return domain.AnswerResult{}
}
