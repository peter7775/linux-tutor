package domain

type Progress struct { Correct, Wrong int }

type ShellTask struct {
	TopicCode string
	Prompt    string
	Expected  string
	Hint      string
}
