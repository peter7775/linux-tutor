package domain

type Progress struct { Correct, Wrong int }

type Topic struct { Code, Title, Area string }

type Task struct { ID string; Topic Topic; Kind string; Prompt string; Expected string; Choices []string; Hint string }
