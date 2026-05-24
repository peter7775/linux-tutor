package domain

import "time"

type Progress struct{ Correct, Wrong int }
type Topic struct{ Code, Title, Area string }
type Task struct {
	ID       string
	Topic    Topic
	Kind     string
	Prompt   string
	Expected string
	Choices  []string
	Question *Question
}
type AreaStat struct {
	Area           string
	Correct, Wrong int
}
type TopicStat struct {
	Code     string
	Correct  int
	Wrong    int
	LastSeen time.Time
}
type Attempt struct {
	TopicCode, Prompt, Answer, Notes string
	ScoreDelta                       int
	CreatedAt                        time.Time
}
