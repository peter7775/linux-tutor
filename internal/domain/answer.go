package domain

type Answer struct {
	QuestionID string
	Value      string
}
type AnswerResult struct {
	Exact      int
	Partial    int
	Wrong      int
	ScoreDelta int
	Notes      string
}
