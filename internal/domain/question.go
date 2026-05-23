package domain

type Question struct {
	ID          string
	TopicID     string
	Type        string
	Prompt      string
	Choices     []string
	Correct     []string
	Explanation string
	Commands    []string
	Difficulty  int
	Hint        string
	Meta        map[string]string
}
