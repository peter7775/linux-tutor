package terminal

import tea "github.com/charmbracelet/bubbletea"

func Start() error {
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	return err
}
