package terminal

import (
	"fmt"
	"linux-tutor/internal/infra/repository"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type screen int
const (
	screenDashboard screen = iota
	screenQuiz
	screenProgress
)

type Model struct {
	screen   screen
	cursor   int
	input    string
	feedback string
	repo     repository.ProgressRepo
	correct  int
	wrong    int
}

func NewModel(repo repository.ProgressRepo) Model {
	c, w, _ := repo.Load()
	return Model{repo: repo, screen: screenDashboard, cursor: 0, correct: c, wrong: w}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			_ = m.repo.Save(m.correct, m.wrong)
			return m, tea.Quit
		case "up":
			if m.screen == screenDashboard && m.cursor > 0 { m.cursor-- }
		case "down":
			if m.screen == screenDashboard && m.cursor < 2 { m.cursor++ }
		case "enter":
			if m.screen == screenDashboard {
				switch m.cursor { case 0: m.screen = screenQuiz; case 1: m.screen = screenProgress; case 2: _ = m.repo.Save(m.correct,m.wrong); return m, tea.Quit }
			} else if m.screen == screenQuiz {
				if strings.TrimSpace(strings.ToLower(m.input)) == "ps" { m.feedback = "Správně."; m.correct++ } else { m.feedback = "Špatně. Správně je: ps"; m.wrong++ }
				_ = m.repo.Save(m.correct, m.wrong)
				m.input = ""
			}
		case "backspace":
			if m.screen == screenQuiz && len(m.input) > 0 { m.input = m.input[:len(m.input)-1] }
		default:
			if m.screen == screenQuiz && len(msg.String()) == 1 { m.input += msg.String() }
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case screenDashboard:
		items := []string{"Spustit kvíz", "Přehled pokroku", "Konec"}
		out := "linux-tutor

"
		for i, item := range items {
			cursor := " "
			if m.cursor == i { cursor = ">" }
			out += fmt.Sprintf("%s %s
", cursor, item)
		}
		return out + "
Pohyb: šipky, Enter, q"
	case screenQuiz:
		return fmt.Sprintf("Otázka: Jaký příkaz zobrazí běžící procesy?

Odpověď: %s

%s", m.input, m.feedback)
	case screenProgress:
		return fmt.Sprintf("Pokrok

Správně: %d
Špatně: %d

q pro návrat", m.correct, m.wrong)
	default:
		return ""
	}
}
