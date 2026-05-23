package terminal

import (
	"database/sql"
	"linux-tutor/internal/infra/repository"
	tea "github.com/charmbracelet/bubbletea"
)

func Start(db *sql.DB) error { p := tea.NewProgram(NewModel(repository.ProgressRepo{DB: db})); _, err := p.Run(); return err }
