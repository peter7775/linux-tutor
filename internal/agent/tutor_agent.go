package agent

import "linux-tutor/internal/domain"

type TutorAgent struct{ Agent Agent }

func NewTutorAgent(a Agent) TutorAgent                        { return TutorAgent{Agent: a} }
func (t TutorAgent) NextQuestion(code string) domain.Question { return t.Agent.Generate(code) }
func (t TutorAgent) Review(q domain.Question, ans string) domain.AnswerResult {
	return t.Agent.Evaluate(q, ans)
}
