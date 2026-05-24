package ai

import "linux-tutor/internal/domain"

type AIConfig struct{}

type Tutor struct { Config AIConfig; Spec TutorSpec }
func New(cfg AIConfig, spec TutorSpec) Tutor { return Tutor{Config: cfg, Spec: spec} }
func (t Tutor) SuggestNext(current domain.Topic, weak map[string]int) domain.Topic { return current }
func (t Tutor) Explain(q domain.Question, ans string) string { if ans == q.Expected { return "Correct." }; return "Try again." }
