package agent

import (
	"linux-tutor/internal/domain"
)

type Agent struct{}

func New(path string) Agent { return Agent{} }
func (a Agent) Generate(code string) domain.Task {
	return domain.Task{ID: code, Topic: domain.Topic{Code: code, Title: code, Area: "Unknown"}, Kind: "single_command", Prompt: "demo", Expected: "pwd", Hint: "demo"}
}
func (a Agent) Evaluate(task domain.Task, input string) domain.AnswerResult {
	return domain.AnswerResult{ScoreDelta: 0, Notes: ""}
}
func (a Agent) Next(code string) (domain.Topic, error) { return domain.Topic{Code: code}, nil }
