package terminal

import (
	"fmt"
	"linux-tutor/internal/domain"
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
	screen    screen
	cursor    int
	quizzes   []domain.Question
	input     string
	feedback  string
	progress  domain.Progress
}

func NewModel() Model {
	return Model{
		screen: screenDashboard,
		quizzes: []domain.Question{
			{ID: "1", TopicID: "cmd", Prompt: "Jaký příkaz zobrazí běžící procesy?", Answer: "ps"},
			{ID: "2", TopicID: "files", Prompt: "Jaký příkaz vypíše obsah souboru?", Answer: "cat"},
		},
	}
}

func Start() error {
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	return err
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.screen == screenDashboard && m.cursor > 0 { m.cursor-- }
		case "down":
			if m.screen == screenDashboard && m.cursor < 2 { m.cursor++ }
		case "enter":
			if m.screen == screenDashboard {
				switch m.cursor {
				case 0: m.screen = screenQuiz
				case 1: m.screen = screenProgress
				case 2: return m, tea.Quit
				}
			} else if m.screen == screenQuiz {
				q := m.quizzes[0]
				if strings.TrimSpace(strings.ToLower(m.input)) == q.Answer {
					m.feedback = "Správně."
					m.progress.Correct++
				} else {
					m.feedback = fmt.Sprintf("Špatně. Správně je: %s", q.Answer)
					m.progress.Wrong++
				}
				m.input = ""
			}
		case "backspace":
			if m.screen == screenQuiz && len(m.input) > 0 { m.input = m.input[:len(m.input)-1] }
		default:
			if m.screen == screenQuiz && len(msg.String()) == 1 {
				m.input += msg.String()
			}
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
		out += "
Pohyb: šipky, potvrzení: Enter, ukončení: q"
		return out
	case screenQuiz:
		q := m.quizzes[0]
		return fmt.Sprintf("Kvíz

Otázka: %s

Odpověď: %s

%s

Enter pro vyhodnocení, q pro návrat", q.Prompt, m.input, m.feedback)
	case screenProgress:
		return fmt.Sprintf("Pokrok

Správně: %d
Špatně: %d

q pro návrat", m.progress.Correct, m.progress.Wrong)
	default:
		return ""
	}
}
