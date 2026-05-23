package terminal

import (
	"fmt"
	"linux-tutor/internal/agent"
	"linux-tutor/internal/domain"
	"linux-tutor/internal/infra/repository"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
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
	task domain.Task
	correct int
	wrong int
	score int
	weak map[string]int
	area map[string]domain.AreaStat
	topic map[string]domain.TopicStat
	bar progress.Model
	ready bool
}

func NewModel(repo repository.ProgressRepo, ag agent.Agent) Model {
	c, w, _ := repo.Load()
	b := progress.New(progress.WithDefaultGradient())
	return Model{repo: repo, agent: ag, task: ag.Generate("103.4"), correct: c, wrong: w, output: []string{"Mini shell ready.", "Use: task, answer, next, type, topic, help"}, weak: map[string]int{}, area: map[string]domain.AreaStat{}, topic: map[string]domain.TopicStat{}, bar: b}
}

func (m Model) Init() tea.Cmd { return nil }

func (m *Model) updateProgress() { m.bar.SetPercent(float64(m.correct) / math.Max(1, float64(m.correct+m.wrong+1))) }

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

func (m *Model) bump(area string, code string, delta int) {
	sa := m.area[area]
	sa.Area = area
	st := m.topic[code]
	st.Code = code
	st.LastSeen = time.Now()
	if delta > 0 { sa.Correct++; st.Correct++ } else { sa.Wrong++; st.Wrong++ }
	m.area[area] = sa
	m.topic[code] = st
}

func (m *Model) pickAdaptive() domain.Topic {
	bestCode, bestScore := "", math.MaxInt
	for code, st := range m.topic {
		sum := st.Correct + st.Wrong
		if sum == 0 { continue }
		s := st.Wrong*2 - st.Correct
		if s < bestScore { bestScore, bestCode = s, code }
	}
	if bestCode != "" { for _, t := range m.agent.Catalog.Topics { if t.Code == bestCode { return t } } }
	bestWeak, maxWeak := "", -1
	for code, cnt := range m.weak { if cnt > maxWeak { bestWeak, maxWeak = code, cnt } }
	if bestWeak != "" { for _, t := range m.agent.Catalog.Topics { if t.Code == bestWeak { return t } } }
	if n, err := m.agent.Next(m.task.Topic.Code); err == nil { return n }
	return m.task.Topic
}

func (m *Model) nextAdaptive() { t := m.pickAdaptive(); m.task = m.agent.Generate(t.Code); m.output = append(m.output, "Adaptive topic: "+t.Code+" - "+t.Area) }

func (m *Model) runShell() {
	cmd := strings.TrimSpace(m.input)
	m.output = append(m.output, "> "+cmd)
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
		if r.ScoreDelta == 10 { m.correct++ } else if r.ScoreDelta == 5 { m.weak[m.task.Topic.Code]++ } else { m.wrong++; m.weak[m.task.Topic.Code]++ }
		m.score += r.ScoreDelta
		m.bump(m.task.Topic.Area, m.task.Topic.Code, r.ScoreDelta)
		m.output = append(m.output, fmt.Sprintf("%s (+%d)", r.Notes, r.ScoreDelta))
		_ = m.repo.Save(m.correct, m.wrong)
	default:
		if cmd != "" { m.output = append(m.output, "Nepodporovaný příkaz v mini shellu."); m.wrong++; m.score-- }
	}
	_ = m.repo.Save(m.correct, m.wrong)
	m.updateProgress()
	m.input = ""
}

func (m Model) renderStats() string {
	p := lipgloss.NewStyle().Bold(true).Render("LPIC stats") + "
"
	for _, a := range m.area { p += fmt.Sprintf("%s: %d correct, %d wrong
", a.Area, a.Correct, a.Wrong) }
	if len(m.area) == 0 { p += "No area data yet
" }
	p += "
Weak topics:
"
	for code, cnt := range m.weak { p += fmt.Sprintf("%s: %d
", code, cnt) }
	if len(m.weak) == 0 { p += "No weak topics yet
" }
	return p
}

func (m Model) View() string {
	switch m.screen {
	case screenDashboard:
		items := []string{"Open learning room", "Open progress", "Quit"}
		out := "linux-tutor

"
		for i, item := range items { c := " "; if m.cursor == i { c = ">" }; out += fmt.Sprintf("%s %s
", c, item) }
		return out + "
Progress: " + m.bar.ViewAs(float64(m.correct)/math.Max(1, float64(m.correct+m.wrong+1))) + "
Use arrows + Enter"
	case screenShell:
		return "Learning room

" + strings.Join(m.output, "
") + "

> " + m.input + "

Enter runs the command, q returns"
	case screenProgress:
		return fmt.Sprintf("Progress

Correct: %d
Wrong: %d
Score: %d

%s
q to return", m.correct, m.wrong, m.score, m.renderStats())
	default:
		return ""
	}
}
