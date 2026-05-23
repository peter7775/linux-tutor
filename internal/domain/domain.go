package domain

import "time"

type Topic struct{ Code, Title, Area string }
type Task struct {
	ID                           string
	Topic                        Topic
	Kind, Prompt, Expected, Hint string
	Choices                      []string
}
type AnswerResult struct {
	Exact, Partial, Wrong, ScoreDelta int
	Notes                             string
}
type AreaStat struct {
	Area           string
	Correct, Wrong int
}
type TopicStat struct {
	Code           string
	Correct, Wrong int
	LastSeen       time.Time
}
type Attempt struct {
	TopicCode, Prompt, Answer, Notes string
	ScoreDelta                       int
	CreatedAt                        time.Time
}
