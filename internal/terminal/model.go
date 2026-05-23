package terminal

import "github.com/charmbracelet/bubbletea"

type Model struct{}

func NewModel(any, any) *Model                           { return &Model{} }
func (m *Model) Init() tea.Cmd                           { return nil }
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *Model) View() string                            { return "Linux Tutor" }
