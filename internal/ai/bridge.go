package ai

import (
	"linux-tutor/internal/agent"
	"linux-tutor/internal/domain"
)

func BindCatalog(a *agent.Agent, path string) error {
	topics, err := LoadCatalog(path)
	if err != nil {
		return err
	}
	a.Catalog.Topics = topics
	return nil
}

func ToQuestion(a agent.Agent, code string) domain.Question { return a.Generate(code) }
