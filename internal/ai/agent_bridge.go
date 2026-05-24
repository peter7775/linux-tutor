package ai

import ("linux-tutor/internal/domain"; "linux-tutor/internal/agent")

func BindCatalog(a *agent.Agent, path string) error { topics, err := LoadCatalog(path); if err != nil { return err }; a.Catalog.Topics = topics; return nil }
func ToQuestion(a agent.Agent, code string) domain.Question { return a.Generate(code) }
