package domain

type Answer struct {
	QuestionID string
	Value      string
}
type AnswerResult struct {
	Exact      int
	Partial    bool
	Wrong      int
	ScoreDelta int
	Notes      string
}
