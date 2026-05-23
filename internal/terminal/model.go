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
}

func NewModel(repo repository.ProgressRepo, ag agent.Agent) Model {
	c, w, _ := repo.Load()
	b := progress.New(progress.WithDefaultGradient())
	return Model{repo: repo, agent: ag, task: ag.Generate("103.4"), correct: c, wrong: w, output: []string{"Mini shell ready.", "Use: task, answer, next, type, topic, help"}, weak: map[string]int{}, area: map[string]domain.AreaStat{}, topic: map[string]domain.TopicStat{}, bar: b}
}

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

func (m *Model) bumpArea(area string, ok bool) {
	s := m.area[area]
	s.Area = area
	if ok { s.Correct++ } else { s.Wrong++ }
	m.area[area] = s
}

func (m *Model) bumpTopic(code string, ok bool) {
	s := m.topic[code]
	s.Code = code
	s.LastSeen = time.Now()
	if ok { s.Correct++ } else { s.Wrong++ }
	m.topic[code] = s
}

func (m *Model) pickAdaptive() domain.Topic {
	bestCode := ""
	bestScore := math.MaxInt
	for code, s := range m.topic {
		sum := s.Correct + s.Wrong
		if sum == 0 { continue }
		score := s.Wrong*2 - s.Correct
		if score < bestScore { bestScore, bestCode = score, code }
	}
	if bestCode != "" {
		for _, t := range m.agent.Catalog.Topics { if t.Code == bestCode { return t } }
	}
	bestWeak := ""
	maxWeak := -1
	for code, cnt := range m.weak { if cnt > maxWeak { bestWeak, maxWeak = code, cnt } }
	if bestWeak != "" { for _, t := range m.agent.Catalog.Topics { if t.Code == bestWeak { return t } } }
	if n, err := m.agent.Next(m.task.Topic.Code); err == nil { return n }
	return m.task.Topic
}

func (m *Model) nextAdaptive() {
	t := m.pickAdaptive()
	m.task = m.agent.Generate(t.Code)
	m.output = append(m.output, "Adaptive topic: "+t.Code+" - "+t.Area)
}

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
		m.bumpArea(m.task.Topic.Area, r.ScoreDelta > 0)
		m.bumpTopic(m.task.Topic.Code, r.ScoreDelta > 0)
		m.output = append(m.output, fmt.Sprintf("%s (+%d)", r.Notes, r.ScoreDelta))
		_ = m.repo.Save(m.correct, m.wrong)
	default:
		if cmd != "" { m.output = append(m.output, "Nepodporovaný příkaz v mini shellu."); m.wrong++; m.score-- }
	}
	_ = m.repo.Save(m.correct, m.wrong)
	m.input = ""
}

func (m Model) areaLines() []string {
	areas := []string{}
	for _, a := range m.area { areas = append(areas, fmt.Sprintf("- %s: %d correct, %d wrong", a.Area, a.Correct, a.Wrong)) }
	if len(areas) == 0 { areas = []string{"- no data yet"} }
	return areas
}

func (m Model) topicLines() []string {
	lines := []string{}
	for code, t := range m.topic { lines = append(lines, fmt.Sprintf("- %s: %d correct, %d wrong", code, t.Correct, t.Wrong)) }
	if len(lines) == 0 { lines = []string{"- no data yet"} }
	return lines
}

func (m Model) View() string {
	switch m.screen {
	case screenDashboard:
		items := []string{"Otevřít mini shell", "Přehled pokroku", "Konec"}
		out := "linux-tutor\n\n"
		for i, item := range items {
			c := " "
			if m.cursor == i {
				c = ">"
			}
			out += fmt.Sprintf("%s %s\n", c, item)
		}
		return out + fmt.Sprintf("\nProgress: %s\n\nPohyb: šipky, Enter, q", m.bar.ViewAs(float64(m.correct)/math.Max(1, float64(m.correct+m.wrong+1))))
	case screenShell:
		return "Mini shell\n\n" + strings.Join(m.output, "\n") + "\n\n> " + m.input + "\n\nEnter spustí příkaz, q vrátí zpět"
	case screenProgress:
		return "Pokrok\n\n" +
			"Správně: " + fmt.Sprint(m.correct) + "\n" +
			"Špatně: " + fmt.Sprint(m.wrong) + "\n" +
			"Skóre: " + fmt.Sprint(m.score) + "\n\n" +
			"LPIC areas:\n" +
			strings.Join(m.areaLines(), "\n") + "\n\n" +
			"Weak topics:\n" +
			strings.Join(m.topicLines(), "\n") + "\n\n" +
			"q pro návrat"
	default:
		return ""
	}
}
