package terminal

import (
	"database/sql"
	"linux-tutor/internal/agent"
	"linux-tutor/internal/infra/repository"

	tea "github.com/charmbracelet/bubbletea"
)

func Start(db *sql.DB) error {
	repo := repository.ProgressRepo{DB: db}
	ag := agent.New("internal/catalog/lpic.json")
	p := tea.NewProgram(NewModel(repo, ag))
	_, err := p.Run()
	return err
}
