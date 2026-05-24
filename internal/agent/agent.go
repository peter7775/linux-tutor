package agent

import "linux-tutor/internal/domain"

type Catalog struct{ Topics []domain.Topic }
type Agent struct{ Catalog Catalog }

func New(path string) Agent {
	return Agent{Catalog: Catalog{Topics: []domain.Topic{{Code: "103.4", Title: "Use streams, pipes and redirects", Area: "GNU and Unix Commands"}, {Code: "103.5", Title: "Create, monitor and kill processes", Area: "GNU and Unix Commands"}}}}
}
func (a Agent) GetCatalog() Catalog { return a.Catalog }
func (a Agent) Generate(code string) domain.Task {
	t := domain.Topic{Code: code, Title: code, Area: "Unknown"}
	for _, x := range a.Catalog.Topics {
		if x.Code == code {
			t = x
			break
		}
	}
	return domain.Task{ID: code + "-1", Topic: t, Kind: "single_command", Prompt: "Zobraz pracovní adresář.", Expected: "pwd", Hint: "basic shell"}
}
func (a Agent) Evaluate(task domain.Task, input string) domain.AnswerResult {
	if input == task.Expected {
		return domain.AnswerResult{Exact: 10, ScoreDelta: 10, Notes: "correct"}
	}
	return domain.AnswerResult{Wrong: 1, ScoreDelta: 0, Notes: "wrong"}
}
