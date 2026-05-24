package agent

import "linux-tutor/internal/domain"

func ScoreAnswer(expected, actual string) domain.AnswerResult {
	if expected == actual {
		return domain.AnswerResult{Exact: 10, ScoreDelta: 10, Notes: "correct"}
	}
	if actual != "" {
		return domain.AnswerResult{Partial: 5, ScoreDelta: 5, Notes: "partial"}
	}
	return domain.AnswerResult{Wrong: 1, ScoreDelta: 0, Notes: "wrong"}
}
