package ai

type TutorSpec struct { Purpose string; TopicScope []string; MaxHints int }
func DefaultTutorSpec() TutorSpec { return TutorSpec{Purpose:"LPIC tutoring", TopicScope:[]string{"103.4","103.5"}, MaxHints:2} }
