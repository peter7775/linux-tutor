package gui

import (
	"database/sql"
	"fmt"
	"linux-tutor/internal/agent"
	"linux-tutor/internal/domain"
	"linux-tutor/internal/infra/repository"
	"math"
	"sort"
	"strings"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type state struct {
	repo          repository.ProgressRepo
	ag            agent.Agent
	topicIdx      int
	task          domain.Task
	correct       int
	wrong         int
	score         int
	weak          map[string]int
	area          map[string]domain.AreaStat
	topic         map[string]domain.TopicStat
	attempts      []domain.Attempt
	testRemaining int
	testMode      bool
}

func Start(db *sql.DB) {
	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())
	w := a.NewWindow("linux-tutor")
	w.Resize(fyne.NewSize(1280, 860))

	st := newState(db)
	content, timerLbl, questionLbl, answerEntry, feedbackLbl, progressBar, statsEntry := buildUI(st)

	w.SetContent(content)
	st.startCountdown(timerLbl, progressBar, statsEntry, feedbackLbl, questionLbl, answerEntry)
	w.ShowAndRun()
}

func newState(db *sql.DB) *state {
	repo := repository.ProgressRepo{DB: db}
	c, w, _ := repo.Load()
	st := &state{
		repo:          repo,
		ag:            agent.New("internal/catalog/lpic.json"),
		weak:          map[string]int{},
		area:          map[string]domain.AreaStat{},
		topic:         map[string]domain.TopicStat{},
		attempts:      []domain.Attempt{},
		correct:       c,
		wrong:         w,
		testRemaining: 60,
	}
	if len(st.ag.Catalog.Topics) > 0 {
		st.task = st.ag.Generate(st.ag.Catalog.Topics[0].Code).Task
	}
	return st
}

func (s *state) add(delta int, ans string) {
	switch delta {
	case 10:
		s.correct++
	case 5:
		s.weak[s.task.Topic.Code]++
	default:
		s.wrong++
		s.weak[s.task.Topic.Code]++
	}

	s.score += delta

	a := s.area[s.task.Topic.Area]
	a.Area = s.task.Topic.Area
	if delta > 0 {
		a.Correct++
	} else {
		a.Wrong++
	}
	s.area[s.task.Topic.Area] = a

	t := s.topic[s.task.Topic.Code]
	t.Code = s.task.Topic.Code
	t.LastSeen = time.Now()
	if delta > 0 {
		t.Correct++
	} else {
		t.Wrong++
	}
	s.topic[s.task.Topic.Code] = t

	s.attempts = append(s.attempts, domain.Attempt{
		TopicCode:  s.task.Topic.Code,
		Prompt:     s.task.Prompt,
		Answer:     ans,
		Notes:      fmt.Sprintf("%d", delta),
		ScoreDelta: delta,
		CreatedAt:  time.Now(),
	})

	_ = s.repo.Save(s.correct, s.wrong)
	_ = s.repo.SaveAttempt(domain.Attempt{
		TopicCode: s.task.Topic.Code,
		Prompt:    s.task.Prompt,
		Answer:    ans,
		Notes:     fmt.Sprintf("%d", delta),
	})
}

func (s *state) nextAdaptive() {
	best := ""
	bestScore := math.MaxInt

	for code, st := range s.topic {
		if st.Correct+st.Wrong > 0 {
			sc := st.Wrong*2 - st.Correct
			if sc < bestScore {
				bestScore = sc
				best = code
			}
		}
	}

	if best == "" && len(s.ag.Catalog.Topics) > 0 {
		best = s.ag.Catalog.Topics[s.topicIdx%len(s.ag.Catalog.Topics)].Code
		s.topicIdx++
	}

	if best != "" {
		s.task = s.ag.Generate(best).Task
	}
}

func buildUI(s *state) (fyne.CanvasObject, *widget.Label, *widget.Label, *widget.Entry, *widget.Label, *widget.ProgressBar, *widget.Entry) {
	title := widget.NewLabelWithStyle("Linux Tutor", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("LPIC-1 learning flow with adaptive practice and timed tests")
	timer := widget.NewLabel("60s")
	progress := widget.NewProgressBar()
	progress.SetValue(0.5)

	answer := widget.NewEntry()
	answer.SetPlaceHolder("Type your answer here")

	feedback := widget.NewLabel("")
	question := widget.NewLabel("")
	question.Wrapping = fyne.TextWrapWord

	lesson := widget.NewLabel("")
	lesson.Wrapping = fyne.TextWrapWord

	stats := widget.NewMultiLineEntry()
	stats.Disable()

	topicList := widget.NewList(
		func() int { return len(s.ag.Catalog.Topics) },
		func() fyne.CanvasObject { return widget.NewLabel("topic") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < len(s.ag.Catalog.Topics) {
				t := s.ag.Catalog.Topics[i]
				o.(*widget.Label).SetText(fmt.Sprintf("%s  %s", t.Code, t.Title))
			}
		},
	)

	refresh := func(code string) {
		s.task = s.ag.Generate(code).Task
		lesson.SetText(lessonText(s.task.Topic))
		question.SetText(renderQuestion(s.task))
		stats.SetText(renderStats(s))
		feedback.SetText("")
		answer.SetText("")
	}

	if len(s.ag.Catalog.Topics) > 0 {
		refresh(s.ag.Catalog.Topics[0].Code)
	}

	topicList.OnSelected = func(id widget.ListItemID) {
		if id < len(s.ag.Catalog.Topics) {
			refresh(s.ag.Catalog.Topics[id].Code)
		}
	}

	submit := func() {
		ans := strings.TrimSpace(answer.Text)
		r := s.ag.Evaluate(*s.task.Question, ans)
		s.add(r.ScoreDelta, ans)
		feedback.SetText(renderFeedback(r, s.task, ans))
		progress.SetValue(float64(s.correct) / math.Max(1, float64(s.correct+s.wrong+1)))
		stats.SetText(renderStats(s))
		s.nextAdaptive()
		question.SetText(renderQuestion(s.task))
		answer.SetText("")
	}

	testStart := widget.NewButton("Start 60s test", func() {
		s.testMode = true
		s.testRemaining = 60
		timer.SetText("60s")
		feedback.SetText("Test started. Focus on accuracy and speed.")
	})

	checkButton := widget.NewButton("Check answer", submit)
	nextButton := widget.NewButton("Next question", func() {
		s.nextAdaptive()
		question.SetText(renderQuestion(s.task))
		feedback.SetText("")
		answer.SetText("")
	})

	controls := container.NewVBox(checkButton, nextButton, testStart)
	mainCard := container.NewVBox(title, subtitle, widget.NewSeparator(), lesson, widget.NewSeparator(), question, answer, feedback, controls)
	progressCard := container.NewVBox(progress, stats)

	tabs := container.NewAppTabs(
		container.NewTabItem("Practice", mainCard),
		container.NewTabItem("Progress", progressCard),
	)

	left := container.NewVBox(widget.NewLabelWithStyle("Topics", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), topicList)
	right := container.NewVBox(timer, tabs)
	root := container.NewBorder(nil, nil, left, nil, right)

	return root, timer, question, answer, feedback, progress, stats
}

func (s *state) startCountdown(timer *widget.Label, progress *widget.ProgressBar, stats *widget.Entry, feedback *widget.Label, question *widget.Label, answer *widget.Entry) {
	_ = timer
	_ = progress
	_ = stats
	_ = feedback
	_ = question
	_ = answer
}

func lessonText(t domain.Topic) string {
	lessons := map[string]string{
		"103.4": "Learn how to redirect stdout and stderr, chain commands, and inspect command output efficiently.",
		"103.5": "Practice process discovery, process control, and safe termination workflows.",
		"104.5": "Understand chmod, chown, umask, and how permissions affect access.",
		"105.2": "Write, execute, and debug simple shell scripts with predictable structure.",
		"107.1": "Work with users, groups, and account-related files on the system.",
		"107.2": "Schedule recurring jobs and one-off tasks using cron and at.",
		"109.3": "Diagnose connectivity, routing, DNS, and basic network issues.",
		"110.2": "Check services, service states, and host security basics.",
	}

	body := lessons[t.Code]
	if body == "" {
		body = "This topic is available in the catalog, but no lesson text is defined yet."
	}

	return fmt.Sprintf("%s — %s\n\n%s\n\nArea: %s", t.Code, t.Title, body, t.Area)
}

func renderQuestion(task domain.Task) string {
	if task.ID == "" {
		return "No question loaded yet."
	}

	lines := []string{
		task.Prompt,
		fmt.Sprintf("Topic: %s", task.Topic.Title),
		fmt.Sprintf("Area: %s", task.Topic.Area),
		fmt.Sprintf("Kind: %s", task.Kind),
	}

	if len(task.Choices) > 0 {
		lines = append(lines, "Choices:")
		for i, c := range task.Choices {
			lines = append(lines, fmt.Sprintf("  %d) %s", i+1, c))
		}
	}

	return strings.Join(lines, "\n")
}

func renderFeedback(r domain.AnswerResult, task domain.Task, ans string) string {
	status := "Wrong"

	switch r.ScoreDelta {
	case 10:
		status = "Correct"
	case 5:
		status = "Partially correct"
	}

	return fmt.Sprintf("%s. Score +%d. Your answer: %q. Expected hint: %q", status, r.ScoreDelta, ans, task.Expected)
}

func renderStats(s *state) string {
	areas := make([]string, 0, len(s.area))
	for _, a := range s.area {
		areas = append(areas, fmt.Sprintf("%s: %d correct / %d wrong", a.Area, a.Correct, a.Wrong))
	}
	sort.Strings(areas)
	if len(areas) == 0 {
		areas = []string{"No area data yet."}
	}

	topics := make([]string, 0, len(s.topic))
	for _, t := range s.topic {
		topics = append(topics, fmt.Sprintf("%s: %d correct / %d wrong", t.Code, t.Correct, t.Wrong))
	}
	sort.Strings(topics)
	if len(topics) == 0 {
		topics = []string{"No topic data yet."}
	}

	weak := make([]string, 0, len(s.weak))
	for code, cnt := range s.weak {
		weak = append(weak, fmt.Sprintf("%s: %d", code, cnt))
	}
	sort.Strings(weak)
	if len(weak) == 0 {
		weak = []string{"No weak topics yet."}
	}

	return fmt.Sprintf("Correct: %d\nWrong: %d\nScore: %d\nAttempts: %d\n\nWeak topics:\n%s\n\nAreas:\n%s\n\nTopic stats:\n%s",
		s.correct, s.wrong, s.score, len(s.attempts), strings.Join(weak, "\n"), strings.Join(areas, "\n"), strings.Join(topics, "\n"))
}
