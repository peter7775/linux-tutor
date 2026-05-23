package domain

import "time"

type Progress struct {
	UserID       string
	TopicID      string
	CorrectCount int
	WrongCount   int
	MasteryScore float64
	LastSeenAt   time.Time
	Streak       int
	WeakSignals  []string
}
