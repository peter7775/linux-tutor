package terminal

import (
	"database/sql"
	"fmt"
	"linux-tutor/internal/agent"
	"linux-tutor/internal/domain"
	"linux-tutor/internal/infra/repository"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int
const ( screenDashboard screen = iota; screenLesson; screenPractice; screenTest; screenProgress )

type Model struct { screen screen; cursor int; input string; output []string; repo repository.ProgressRepo; ag agent.Agent; topicIdx int; task domain.Task; correct, wrong, score int; weak map[string]int; area map[string]domain.AreaStat; topic map[string]domain.TopicStat; attempts []domain.Attempt; bar progress.Model; testRemaining int; testMode bool }

func NewModel(repo repository.ProgressRepo, ag agent.Agent) Model { c,w,_ := repo.Load(); b := progress.New(progress.WithDefaultGradient()); m := Model{repo:repo, ag:ag, output:[]string{"Linux tutor ready."}, weak: map[string]int{}, area: map[string]domain.AreaStat{}, topic: map[string]domain.TopicStat{}, attempts: []domain.Attempt{}, bar:b, correct:c, wrong:w}; m.loadNext(); return m }
func (m Model) Init() tea.Cmd { return m.tickCmd() }
func (m Model) tickCmd() tea.Cmd { return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) }) }
type tickMsg time.Time
func (m *Model) syncProgress() { total := m.correct + m.wrong; m.bar.SetPercent(float64(m.correct) / max(1, float64(total+1))) }
func max(a, b float64) float64 { if a > b { return a }; return b }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { switch msg := msg.(type) { case tickMsg: if m.testMode && m.testRemaining > 0 { m.testRemaining--; if m.testRemaining == 0 { m.testMode = false; m.output = append(m.output, "Test finished.") } }; return m, m.tickCmd(); case tea.KeyMsg: switch msg.String() { case "ctrl+c", "q": _ = m.repo.Save(m.correct, m.wrong); return m, tea.Quit; case "up": if m.screen == screenDashboard && m.cursor > 0 { m.cursor-- }; case "down": if m.screen == screenDashboard && m.cursor < 4 { m.cursor++ }; case "enter": if m.screen == screenDashboard { switch m.cursor { case 0: m.screen = screenLesson; case 1: m.screen = screenPractice; case 2: m.screen = screenTest; case 3: m.screen = screenProgress; case 4: _ = m.repo.Save(m.correct, m.wrong); return m, tea.Quit } } else if m.screen == screenLesson { m.loadNext() } else if m.screen == screenPractice || m.screen == screenTest { m.submit() }; case "backspace": if (m.screen == screenPractice || m.screen == screenTest) && len(m.input) > 0 { m.input = m.input[:len(m.input)-1] }; default: if (m.screen == screenPractice || m.screen == screenTest) && len(msg.String()) == 1 { m.input += msg.String() } } }; return m, nil }
func (m *Model) loadNext() { if len(m.ag.Catalog.Topics) == 0 { return }; t := m.ag.Catalog.Topics[m.topicIdx%len(m.ag.Catalog.Topics)]; m.task = m.ag.Generate(t.Code); m.topicIdx++ }
func (m *Model) nextAdaptive() { best := ""; bestScore := 1<<30; for code, st := range m.topic { if st.Correct+st.Wrong > 0 { sc := st.Wrong*2 - st.Correct; if sc < bestScore { bestScore, best = sc, code } } }; if best == "" { m.loadNext(); return }; m.task = m.ag.Generate(best) }
func (m *Model) addStats(delta int, ans string) { a := m.area[m.task.Topic.Area]; a.Area = m.task.Topic.Area; if delta > 0 { a.Correct++ } else { a.Wrong++ }; m.area[m.task.Topic.Area] = a; t := m.topic[m.task.Topic.Code]; t.Code = m.task.Topic.Code; t.LastSeen = time.Now(); if delta > 0 { t.Correct++ } else { t.Wrong++ }; m.topic[m.task.Topic.Code] = t; m.attempts = append(m.attempts, domain.Attempt{TopicCode:m.task.Topic.Code, Prompt:m.task.Prompt, Answer:ans, Notes:fmt.Sprintf("%d", delta), ScoreDelta:delta, CreatedAt:time.Now()}) }
func (m *Model) submit() { ans := strings.TrimSpace(m.input); r := m.ag.Evaluate(m.task, ans); if r.ScoreDelta == 10 { m.correct++ } else if r.ScoreDelta == 5 { m.weak[m.task.Topic.Code]++ } else { m.wrong++; m.weak[m.task.Topic.Code]++ }; m.score += r.ScoreDelta; m.addStats(r.ScoreDelta, ans); _ = m.repo.Save(m.correct, m.wrong); _ = m.repo.SaveAttempt(m.task.Topic.Code, m.task.Prompt, ans, r.Notes, r.ScoreDelta); m.output = append(m.output, fmt.Sprintf("%s (+%d)", r.Notes, r.ScoreDelta)); m.input = ""; m.syncProgress(); if m.screen == screenTest && m.testRemaining > 0 { m.nextAdaptive() } else { m.nextAdaptive() } }
func (m Model) View() string { switch m.screen { case screenDashboard: return "linux-tutor

> Lesson
  Practice
  Test
  Progress
  Quit

Progress: " + m.bar.View(); case screenLesson: return "Lesson

" + m.lessonView() + "

Enter next, q back"; case screenPractice: return "Practice

" + m.practiceView() + "

> " + m.input + "

Enter submit, q back"; case screenTest: return fmt.Sprintf("Test (%ds left)

%s

> %s

Enter submit, q back", m.testRemaining, m.practiceView(), m.input); case screenProgress: return "Progress

" + m.progressView() + "

q back"; default: return "" } }
func (m Model) lessonView() string { if len(m.ag.Catalog.Topics) == 0 { return "No topics" }; t := m.ag.Catalog.Topics[(m.topicIdx-1+len(m.ag.Catalog.Topics))%len(m.ag.Catalog.Topics)]; lessons := map[string]string{"103.4":"Practice redirecting stdout and stderr.","103.5":"Find, inspect, and kill processes.","104.5":"Work with chmod, chown, and umask.","105.2":"Write simple shell scripts.","107.1":"Learn /etc/passwd and groups.","107.2":"Schedule tasks with cron.","109.3":"Diagnose DNS, routing, and connectivity.","110.2":"Check services and host security."}; return fmt.Sprintf("%s
%s

%s

Task preview: %s", t.Code, t.Title, lessons[t.Code], m.ag.Generate(t.Code).Prompt) }
func (m Model) practiceView() string { if m.task.ID == "" { return "No task loaded" }; return fmt.Sprintf("[%s] %s
Hint: %s", m.task.Kind, m.task.Prompt, m.task.Hint) }
func (m Model) progressView() string { out := fmt.Sprintf("Correct: %d
Wrong: %d
Score: %d
Attempts: %d

Weak topics:
", m.correct, m.wrong, m.score, len(m.attempts)); for code, c := range m.weak { out += fmt.Sprintf("%s: %d
", code, c) }; out += "
Areas:
"; for _, a := range m.area { out += fmt.Sprintf("%s: %d correct / %d wrong
", a.Area, a.Correct, a.Wrong) }; return out }
