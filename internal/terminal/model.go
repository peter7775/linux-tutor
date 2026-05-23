package terminal

import (
	"fmt"
	"linux-tutor/internal/agent"
	"linux-tutor/internal/domain"
	"linux-tutor/internal/infra/repository"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type screen int
const (
	screenDashboard screen = iota
	screenShell
	screenProgress
)

type Model struct {
	screen screen
	cursor int
	input string
	output []string
	repo repository.ProgressRepo
	agent agent.Agent
	task domain.ShellTask
	correct int
	wrong int
}

func NewModel(repo repository.ProgressRepo, ag agent.Agent) Model {
	c,w,_ := repo.Load()
	return Model{repo: repo, agent: ag, task: ag.Generate("103.5"), correct: c, wrong: w, output: []string{"Mini shell připraven.", "LPIC guidelines aktivní."}}
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
				switch m.cursor { case 0: m.screen = screenShell; case 1: m.screen = screenProgress; case 2: _ = m.repo.Save(m.correct,m.wrong); return m, tea.Quit }
			} else if m.screen == screenShell {
				m.runShell()
			}
		case "backspace":
			if m.screen == screenShell && len(m.input) > 0 { m.input = m.input[:len(m.input)-1] }
		default:
			if m.screen == screenShell && len(msg.String()) == 1 { m.input += msg.String() }
		}
	}
	return m, nil
}

func (m *Model) runShell() {
	cmd := strings.TrimSpace(m.input)
	m.output = append(m.output, "> "+cmd)
	switch cmd {
	case "help":
		m.output = append(m.output, "Příkazy: help, task, answer <cmd>, next, ls, pwd, whoami, clear, exit")
	case "task":
		m.output = append(m.output, "Úloha: "+m.task.Prompt)
		m.output = append(m.output, "Téma: "+m.task.TopicCode+" | Hint: "+m.task.Hint)
	case "next":
		m.task = m.agent.Generate("103.4")
		m.output = append(m.output, "Nová úloha: "+m.task.Prompt)
	case "ls":
		m.output = append(m.output, "cmd  internal  docs  data")
		m.correct++
	case "pwd":
		m.output = append(m.output, "/home/linux-tutor")
		m.correct++
	case "whoami":
		m.output = append(m.output, "student")
		m.correct++
	case "clear":
		m.output = []string{}
	case "exit":
		m.screen = screenDashboard
	default:
		if strings.HasPrefix(cmd, "answer ") {
			ok, msg := m.agent.Evaluate(m.task, strings.TrimSpace(strings.TrimPrefix(cmd, "answer")))
			m.output = append(m.output, msg)
			if ok { m.correct++ } else { m.wrong++ }
			_ = m.repo.Save(m.correct, m.wrong)
		} else if cmd != "" {
			m.output = append(m.output, "Nepodporovaný příkaz v mini shellu.")
			m.wrong++
		}
	}
	_ = m.repo.Save(m.correct, m.wrong)
	m.input = ""
}

func (m Model) View() string {
	switch m.screen {
	case screenDashboard:
		items := []string{"Otevřít mini shell", "Přehled pokroku", "Konec"}
		out := "linux-tutor

"
		for i, item := range items { c := " "; if m.cursor == i { c = ">" }; out += fmt.Sprintf("%s %s
", c, item) }
		return out + "
Pohyb: šipky, Enter, q"
	case screenShell:
		return "Mini shell

" + strings.Join(m.output, "
") + "

> " + m.input + "

Enter spustí příkaz, q vrátí zpět"
	case screenProgress:
		return fmt.Sprintf("Pokrok

Správně: %d
Špatně: %d

q pro návrat", m.correct, m.wrong)
	default:
		return ""
	}
}
