package usecase

import "linux-tutor/internal/domain"

type StartLessonInput struct{ Topic domain.Topic }
type StartLessonOutput struct{ Question domain.Question }
