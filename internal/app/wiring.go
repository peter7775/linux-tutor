package app

import (
	"linux-tutor/internal/agent"
	"linux-tutor/internal/infra/repository"
)

func Build() (agent.Agent, repository.ProgressRepo) { return agent.New(""), repository.ProgressRepo{} }
