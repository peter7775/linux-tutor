package terminal

import (
	"database/sql"
	"linux-tutor/internal/agent"
	"linux-tutor/internal/infra/repository"
	tea "github.com/charmbracelet/bubbletea"
)

func Start(db *sql.DB) error { ag := agent.New("internal/catalog/lpic.json"); p := tea.NewProgram(NewModel(repository.ProgressRepo{DB: db}, ag)); _, err := p.Run(); return err }
