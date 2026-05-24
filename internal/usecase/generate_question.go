package usecase

import "linux-tutor/internal/domain"

type GenerateQuestionInput struct{ Topic domain.Topic }
type GenerateQuestionOutput struct{ Question domain.Question }
