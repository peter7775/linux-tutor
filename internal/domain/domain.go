package domain

type Progress struct { Correct, Wrong int }

type ShellTask struct { TopicCode, Prompt, Expected, Hint string }

type Topic struct { Code, Title, Area string }
