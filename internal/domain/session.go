package domain

import "time"

type Session struct {
	ID          string
	UserID      string
	StartedAt   time.Time
	FinishedAt  *time.Time
	Mode        string
	TopicFilter []string
	Score       int
}
