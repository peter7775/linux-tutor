package domain

type Question struct { ID, TopicID, Prompt, Answer string }

type Progress struct { Correct, Wrong int }
