package terminal

import (
	"database/sql"
	"linux-tutor/internal/infra/repository"
	tea "github.com/charmbracelet/bubbletea"
)

func Start(db *sql.DB) error {
	repo := repository.ProgressRepo{DB: db}
	p := tea.NewProgram(NewModel(repo))
	_, err := p.Run()
	return err
}
