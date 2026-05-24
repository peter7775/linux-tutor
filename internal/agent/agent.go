package agent

import "linux-tutor/internal/domain"

type Catalog struct{ Topics []domain.Topic }
type Agent struct{ Catalog Catalog }

func New(path string) Agent {
	return Agent{Catalog: Catalog{Topics: []domain.Topic{{Code: "103.4", Title: "Use streams, pipes and redirects", Area: "GNU and Unix Commands"}, {Code: "103.5", Title: "Create, monitor and kill processes", Area: "GNU and Unix Commands"}}}}
}
func (a Agent) GetCatalog() Catalog { return a.Catalog }
func (a Agent) Generate(code string) domain.Question {
	t := domain.Topic{Code: code, Title: code, Area: "Unknown"}
	for _, x := range a.Catalog.Topics {
		if x.Code == code {
			t = x
			break
		}
	}
	return domain.Question{ID: code + "-1", Topic: t, Kind: "single_command", Prompt: "Zobraz pracovní adresář.", Expected: "pwd", Hint: "basic shell"}
}
func (a Agent) Evaluate(q domain.Question, ans string) domain.AnswerResult {
	return ScoreAnswer(q.Expected, ans)
}
func (a Agent) RecommendNext(current domain.Topic, weak map[string]int) domain.Topic { return current }
