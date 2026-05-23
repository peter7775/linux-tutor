package gui

import (
	"database/sql"
	"fmt"
	"linux-tutor/internal/agent"
	"linux-tutor/internal/domain"
	"linux-tutor/internal/infra/repository"
	"math"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type state struct { repo repository.ProgressRepo; ag agent.Agent; topicIdx int; task domain.Task; correct, wrong, score int; weak map[string]int; area map[string]domain.AreaStat; topic map[string]domain.TopicStat; attempts []domain.Attempt; testRemaining int; testMode bool; tickStop chan struct{} }

func Start(db *sql.DB) { a := app.New(); w := a.NewWindow("linux-tutor"); w.Resize(fyne.NewSize(1200, 800)); st := newState(db); content, timerLbl := buildUI(st, w); w.SetContent(content); st.startCountdown(timerLbl, w); w.ShowAndRun() }
func newState(db *sql.DB) *state { repo := repository.ProgressRepo{DB: db}; c,w,_ := repo.Load(); st := &state{repo:repo, ag:agent.New("internal/catalog/lpic.json"), weak:map[string]int{}, area:map[string]domain.AreaStat{}, topic:map[string]domain.TopicStat{}, attempts:[]domain.Attempt{}, correct:c, wrong:w, testRemaining:60, tickStop: make(chan struct{})}; if len(st.ag.Catalog.Topics)>0 { st.task = st.ag.Generate(st.ag.Catalog.Topics[0].Code) }; return st }
func (s *state) add(delta int, ans string) { if delta == 10 { s.correct++ } else if delta == 5 { s.weak[s.task.Topic.Code]++ } else { s.wrong++; s.weak[s.task.Topic.Code]++ }; s.score += delta; a := s.area[s.task.Topic.Area]; a.Area = s.task.Topic.Area; if delta > 0 { a.Correct++ } else { a.Wrong++ }; s.area[s.task.Topic.Area] = a; t := s.topic[s.task.Topic.Code]; t.Code = s.task.Topic.Code; t.LastSeen = time.Now(); if delta > 0 { t.Correct++ } else { t.Wrong++ }; s.topic[s.task.Topic.Code] = t; s.attempts = append(s.attempts, domain.Attempt{TopicCode:s.task.Topic.Code, Prompt:s.task.Prompt, Answer:ans, Notes:fmt.Sprintf("%d", delta), ScoreDelta:delta, CreatedAt:time.Now()}); _ = s.repo.Save(s.correct, s.wrong); _ = s.repo.SaveAttempt(s.task.Topic.Code, s.task.Prompt, ans, fmt.Sprintf("%d", delta), delta) }
func (s *state) nextAdaptive() { best := ""; bestScore := 1<<30; for code, st := range s.topic { if st.Correct+st.Wrong > 0 { sc := st.Wrong*2 - st.Correct; if sc < bestScore { bestScore, best = sc, code } } }; if best == "" && len(s.ag.Catalog.Topics)>0 { best = s.ag.Catalog.Topics[s.topicIdx%len(s.ag.Catalog.Topics)].Code; s.topicIdx++ }; if best != "" { s.task = s.ag.Generate(best) } }
func buildUI(s *state, win fyne.Window) (fyne.CanvasObject, *widget.Label) { title := widget.NewLabelWithStyle("linux-tutor", fyne.TextAlignCenter, fyne.TextStyle{Bold:true}); progress := widget.NewProgressBar(); progress.SetValue(0.5); answer := widget.NewEntry(); feedback := widget.NewLabel(""); timer := widget.NewLabel("60s"); topicList := widget.NewList(func() int { return len(s.ag.Catalog.Topics) }, func() fyne.CanvasObject { return widget.NewLabel("topic") }, func(i widget.ListItemID, o fyne.CanvasObject) { if i < len(s.ag.Catalog.Topics) { t := s.ag.Catalog.Topics[i]; o.(*widget.Label).SetText(fmt.Sprintf("%s  %s", t.Code, t.Title)) } }); lesson := widget.NewLabel(""); practice := widget.NewLabel(""); stats := widget.NewMultiLineEntry(); stats.Disable(); refresh := func(code string) { s.task = s.ag.Generate(code); ttitle := code; for _, t := range s.ag.Catalog.Topics { if t.Code == code { ttitle = t.Title; break } }; lesson.SetText(fmt.Sprintf("Lesson %s

%s", code, lessonText(code, ttitle))); practice.SetText(fmt.Sprintf("[%s] %s", s.task.Kind, s.task.Prompt)); stats.SetText(renderStats(s)); feedback.SetText(""); answer.SetText("") }; if len(s.ag.Catalog.Topics)>0 { refresh(s.ag.Catalog.Topics[0].Code) }
	topicList.OnSelected = func(id widget.ListItemID) { if id < len(s.ag.Catalog.Topics) { refresh(s.ag.Catalog.Topics[id].Code) } }
	submit := func() { ans := answer.Text; r := s.ag.Evaluate(s.task, ans); s.add(r.ScoreDelta, ans); feedback.SetText(fmt.Sprintf("%s (+%d)", r.Notes, r.ScoreDelta)); progress.SetValue(float64(s.correct)/math.Max(1, float64(s.correct+s.wrong+1))); stats.SetText(renderStats(s)); s.nextAdaptive(); practice.SetText(fmt.Sprintf("[%s] %s", s.task.Kind, s.task.Prompt)); answer.SetText("") }
	testStart := widget.NewButton("Start 60s test", func() { s.testMode = true; s.testRemaining = 60; s.nextAdaptive(); timer.SetText("60s") })
	lessonTab := container.NewVBox(lesson)
	practiceTab := container.NewVBox(practice, answer, widget.NewButton("Submit", submit), feedback)
	testTab := container.NewVBox(timer, testStart, widget.NewButton("Submit test answer", submit), widget.NewLabel("Test mode uses the same learning engine with a countdown."))
	progressTab := container.NewVBox(progress, stats)
	tabs := container.NewAppTabs(container.NewTabItem("Lesson", lessonTab), container.NewTabItem("Practice", practiceTab), container.NewTabItem("Test", testTab), container.NewTabItem("Progress", progressTab))
	return container.NewBorder(container.NewVBox(title), nil, topicList, nil, tabs), timer
}
func lessonText(code, title string) string { return map[string]string{"103.4":"Practice redirecting stdout and stderr.","103.5":"Find, inspect, and kill processes.","104.5":"Work with chmod, chown, and umask.","105.2":"Write simple shell scripts.","107.1":"Learn /etc/passwd and groups.","107.2":"Schedule tasks with cron.","109.3":"Diagnose DNS, routing, and connectivity.","110.2":"Check services and host security."}[code] }
func renderStats(s *state) string { out := fmt.Sprintf("Correct: %d
Wrong: %d
Score: %d
Attempts: %d

Weak topics:
", s.correct, s.wrong, s.score, len(s.attempts)); for code, c := range s.weak { out += fmt.Sprintf("%s: %d
", code, c) }; out += "
Areas:
"; for _, a := range s.area { out += fmt.Sprintf("%s: %d correct / %d wrong
", a.Area, a.Correct, a.Wrong) }; return strings.TrimSpace(out) }


func (s *state) startCountdown(timer *widget.Label, win fyne.Window) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if s.testRemaining <= 0 { continue }
			s.testRemaining--
			val := s.testRemaining
			fyne.Do(func() {
				timer.SetText(fmt.Sprintf("%ds", val))
				if val == 0 {
					feedback := widget.NewLabel("Time is up")
					_ = feedback
				}
			})
		}
	}()
}
