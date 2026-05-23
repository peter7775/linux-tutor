package domain

import "time"

type Progress struct { Correct, Wrong int }
type Topic struct { Code, Title, Area string }
type Task struct { ID string; Topic Topic; Kind string; Prompt string; Expected string; Choices []string; Hint string }
type AnswerResult struct { Exact int; Partial int; Wrong int; Notes string; ScoreDelta int }
type TopicStat struct { Code string; Correct int; Wrong int; LastSeen time.Time }
