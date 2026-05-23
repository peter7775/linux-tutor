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

type Model struct { screen screen; cursor int; input string; output []string; repo repository.ProgressRepo; agent agent.Agent; task domain.Task; correct int; wrong int; score int; weak map[string]int }

func NewModel(repo repository.ProgressRepo, ag agent.Agent) Model { c, w, _ := repo.Load(); return Model{repo: repo, agent: ag, task: ag.Generate("103.4"), correct: c, wrong: w, output: []string{"Mini shell připraven.", "Použij: task, answer, next, type, topic, help"}, weak: map[string]int{}} }
func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": _ = m.repo.Save(m.correct, m.wrong); return m, tea.Quit
		case "up": if m.screen == screenDashboard && m.cursor > 0 { m.cursor-- }
		case "down": if m.screen == screenDashboard && m.cursor < 2 { m.cursor++ }
		case "enter": if m.screen == screenDashboard { switch m.cursor { case 0: m.screen = screenShell; case 1: m.screen = screenProgress; case 2: _ = m.repo.Save(m.correct, m.wrong); return m, tea.Quit } } else if m.screen == screenShell { m.runShell() }
		case "backspace": if m.screen == screenShell && len(m.input) > 0 { m.input = m.input[:len(m.input)-1] }
		default: if m.screen == screenShell && len(msg.String()) == 1 { m.input += msg.String() }
		}
	}
	return m, nil
}

func (m *Model) addWeak(code string) { m.weak[code]++ }

func (m *Model) nextAdaptive() {
	best := ""
	bestCount := 0
	for code, cnt := range m.weak { if cnt > bestCount { best, bestCount = code, cnt } }
	if best != "" { m.task = m.agent.Generate(best); m.output = append(m.output, "Adaptivně zvoleno slabé téma: "+best); return }
	n, err := m.agent.Next(m.task.Topic.Code)
	if err != nil { m.output = append(m.output, err.Error()); return }
	m.task = m.agent.Generate(n.Code)
	m.output = append(m.output, "Nová úloha: "+m.task.Prompt)
}

func (m *Model) runShell() {
	cmd := strings.TrimSpace(m.input); m.output = append(m.output, "> "+cmd)
	switch {
	case cmd == "help": m.output = append(m.output, "Příkazy: help, task, type, topic, next, answer <...>, ls, pwd, whoami, clear, exit")
	case cmd == "task": m.output = append(m.output, fmt.Sprintf("[%s] %s", m.task.Kind, m.task.Prompt)); if len(m.task.Choices) > 0 { m.output = append(m.output, "Možnosti: "+strings.Join(m.task.Choices, ", ")) }; m.output = append(m.output, "Hint: "+m.task.Hint)
	case cmd == "type": m.output = append(m.output, "Typ úlohy: "+m.task.Kind)
	case cmd == "topic": m.output = append(m.output, fmt.Sprintf("Téma: %s | %s", m.task.Topic.Code, m.task.Topic.Area))
	case cmd == "next": m.nextAdaptive()
	case cmd == "ls": m.output = append(m.output, "cmd  internal  docs  data"); m.correct++; m.score++
	case cmd == "pwd": m.output = append(m.output, "/home/linux-tutor"); m.correct++; m.score++
	case cmd == "whoami": m.output = append(m.output, "student"); m.correct++; m.score++
	case cmd == "clear": m.output = []string{}
	case cmd == "exit": m.screen = screenDashboard
	case strings.HasPrefix(cmd, "answer "):
		ans := strings.TrimSpace(strings.TrimPrefix(cmd, "answer"))
		r := m.agent.Evaluate(m.task, ans)
		m.output = append(m.output, fmt.Sprintf("%s (+%d)", r.Notes, r.ScoreDelta))
		m.score += r.ScoreDelta
		if r.ScoreDelta == 10 { m.correct++ } else if r.ScoreDelta == 5 { m.addWeak(m.task.Topic.Code) } else { m.wrong++; m.addWeak(m.task.Topic.Code) }
		_ = m.repo.Save(m.correct, m.wrong)
	default: if cmd != "" { m.output = append(m.output, "Nepodporovaný příkaz v mini shellu."); m.wrong++; m.score-- }
	}
	_ = m.repo.Save(m.correct, m.wrong); m.input = ""
}

func (m Model) View() string {
	switch m.screen {
	case screenDashboard:
		items := []string{"Otevřít mini shell", "Přehled pokroku", "Konec"}; out := "linux-tutor

"; for i, item := range items { c := " "; if m.cursor == i { c = ">" }; out += fmt.Sprintf("%s %s
", c, item) }; return out + "
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
Skóre: %d
Slabiny: %d témat

q pro návrat", m.correct, m.wrong, m.score, len(m.weak))
	default:
		return ""
	}
}
