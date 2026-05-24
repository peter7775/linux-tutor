package usecase

import "linux-tutor/internal/domain"

func EvaluateProgress(s domain.Session, r domain.AnswerResult) domain.Session {
	s.Score += r.ScoreDelta
	return s
}
