package gui

import (
	"context"
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
	_ "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type state struct {
	repo repository.ProgressRepo
	ag   agent.Agent

	topicIdx int
	task     domain.Question

	correct int
	wrong   int
	score   int

	weak     map[string]int
	area     map[string]domain.AreaStat
	topic    map[string]domain.TopicStat
	attempts []domain.Attempt

	testRemaining int
	testMode      bool
}

func Start(db *sql.DB) {
	a := app.New()
	a.Settings().SetTheme(&lemonTheme{})
	w := a.NewWindow("linux-tutor")
	w.Resize(fyne.NewSize(1440, 900))

	st := newState(db)
	content, timerLbl, questionLbl, answerEntry, feedbackLbl, explanationLbl, progressBar, statsEntry, topicList := buildUI(st)
	w.SetContent(content)

	if len(st.ag.GetCatalog().Topics) > 0 {
		refresh := func(code string) {
			st.task = st.ag.Generate(code)
			updateUI(st, timerLbl, questionLbl, answerEntry, feedbackLbl, explanationLbl, progressBar, statsEntry)
		}
		refresh(st.ag.GetCatalog().Topics[0].Code)

		topicList.OnSelected = func(id widget.ListItemID) {
			if id < len(st.ag.GetCatalog().Topics) {
				refresh(st.ag.GetCatalog().Topics[id].Code)
			}
		}
	}

	answerEntry.OnSubmitted = func(string) {
		submitAnswer(st, answerEntry, feedbackLbl, explanationLbl, progressBar, statsEntry, questionLbl, answerEntry)
	}

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
	if len(st.ag.GetCatalog().Topics) > 0 {
		st.task = st.ag.Generate(st.ag.GetCatalog().Topics[0].Code)
	}
	return st
}

func buildUI(s *state) (fyne.CanvasObject, *widget.Label, *widget.Label, *widget.Entry, *widget.Label, *widget.Label, *widget.ProgressBar, *widget.Entry, *widget.List) {
	title := NewTitleLabel("Linux Tutor")
	subtitle := NewSubheaderLabel("LPIC learning flow with adaptive practice and timed tests")
	timer := widget.NewLabel("60s")
	progress := widget.NewProgressBar()
	progress.SetValue(0.5)

	answer := widget.NewEntry()
	answer.SetPlaceHolder("Type your answer here")

	feedback := widget.NewLabel("")
	explanation := widget.NewLabel("")
	explanation.Wrapping = fyne.TextWrapWord

	question := NewMutedText("No question loaded yet.")
	lesson := NewMutedText("")

	stats := widget.NewMultiLineEntry()
	stats.Disable()

	topics := widget.NewList(
		func() int { return len(s.ag.GetCatalog().Topics) },
		func() fyne.CanvasObject { return widget.NewLabel("topic") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i < len(s.ag.GetCatalog().Topics) {
				t := s.ag.GetCatalog().Topics[i]
				o.(*widget.Label).SetText(fmt.Sprintf("%s %s", t.Code, t.Title))
			}
		},
	)

	topicScroll := container.NewVScroll(topics)
	topicScroll.SetMinSize(fyne.NewSize(360, 700))

	refresh := func(code string) {
		s.task = s.ag.Generate(code)
		lesson.SetText(lessonText(s.task.Topic))
		question.SetText(renderQuestion(s.task))
		feedback.SetText("")
		explanation.SetText("")
		answer.SetText("")
		stats.SetText(renderStats(s))
		progress.SetValue(float64(s.correct) / math.Max(1, float64(s.correct+s.wrong+1)))
	}

	if len(s.ag.GetCatalog().Topics) > 0 {
		refresh(s.ag.GetCatalog().Topics[0].Code)
	}

	topics.OnSelected = func(id widget.ListItemID) {
		if id < len(s.ag.GetCatalog().Topics) {
			refresh(s.ag.GetCatalog().Topics[id].Code)
		}
	}

	checkButton := widget.NewButton("Check answer", func() {
		submitAnswer(s, answer, feedback, explanation, progress, stats, question, answer)
	})

	nextButton := widget.NewButton("Next question", func() {
		s.task = s.nextAdaptive()
		lesson.SetText(lessonText(s.task.Topic))
		question.SetText(renderQuestion(s.task))
		feedback.SetText("")
		explanation.SetText("")
		answer.SetText("")
		stats.SetText(renderStats(s))
		progress.SetValue(float64(s.correct) / math.Max(1, float64(s.correct+s.wrong+1)))
	})

	testStart := widget.NewButton("Start 60s test", func() {
		s.testMode = true
		s.testRemaining = 60
		timer.SetText("60s")
		feedback.SetText("Test started. Focus on accuracy and speed.")
	})

	controls := container.NewVBox(checkButton, nextButton, testStart)

	mainCard := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
		lesson,
		widget.NewSeparator(),
		question,
		answer,
		feedback,
		explanation,
		controls,
	)
	progressCard := container.NewVBox(progress, stats)
	tabs := container.NewAppTabs(
		container.NewTabItem("Practice", mainCard),
		container.NewTabItem("Progress", progressCard),
	)

	left := container.NewVBox(NewSectionLabel("Topics"), topicScroll)
	right := container.NewVBox(timer, tabs)
	root := container.NewHSplit(left, right)
	root.SetOffset(0.27)

	return root, timer, question, answer, feedback, explanation, progress, stats, topics
}

func updateUI(s *state, timer *widget.Label, question *widget.Label, answer *widget.Entry, feedback *widget.Label, explanation *widget.Label, progress *widget.ProgressBar, stats *widget.Entry) {
	timer.SetText(fmt.Sprintf("%ds", s.testRemaining))
	question.SetText(renderQuestion(s.task))
	feedback.SetText("")
	explanation.SetText("")
	answer.SetText("")
	progress.SetValue(float64(s.correct) / math.Max(1, float64(s.correct+s.wrong+1)))
	stats.SetText(renderStats(s))
}

func submitAnswer(s *state, answer *widget.Entry, feedback *widget.Label, explanation *widget.Label, progress *widget.ProgressBar, stats *widget.Entry, question *widget.Label, answerEntry *widget.Entry) {
	ans := strings.TrimSpace(answer.Text)
	if ans == "" || s.task.ID == "" {
		return
	}

	r := s.ag.Evaluate(s.task, ans)
	s.add(r.ScoreDelta, ans)
	feedback.SetText(renderFeedback(r, s.task, ans))
	explanation.SetText(agent.ExplainWithEnv(context.Background(), s.task, ans, r))
	progress.SetValue(float64(s.correct) / math.Max(1, float64(s.correct+s.wrong+1)))
	stats.SetText(renderStats(s))

	if r.ScoreDelta > 0 {
		s.task = s.nextAdaptive()
		question.SetText(renderQuestion(s.task))
		answerEntry.SetText("")
	}
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

	s.attempts = append(s.attempts, domain.Attempt{TopicCode: s.task.Topic.Code, Prompt: s.task.Prompt, Answer: ans, Notes: fmt.Sprintf("%d", delta), ScoreDelta: delta, CreatedAt: time.Now()})
	_ = s.repo.Save(s.correct, s.wrong)
	_ = s.repo.SaveAttempt(domain.Attempt{TopicCode: s.task.Topic.Code, Prompt: s.task.Prompt, Answer: ans, Notes: fmt.Sprintf("%d", delta), ScoreDelta: delta, CreatedAt: time.Now()})
}

func (s *state) nextAdaptive() domain.Question {
	tops := s.ag.GetCatalog().Topics
	if len(tops) == 0 {
		return s.task
	}
	if s.task.ID != "" {
		for i, t := range tops {
			if t.Code == s.task.Topic.Code {
				next := tops[(i+1)%len(tops)]
				return s.ag.Generate(next.Code)
			}
		}
	}
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
	if best != "" {
		return s.ag.Generate(best)
	}
	return s.ag.Generate(tops[0].Code)
}

func lessonText(t domain.Topic) string {
	return fmt.Sprintf("%s — %sArea: %s", t.Code, t.Title, t.Area)
}

func renderQuestion(task domain.Question) string {
	if task.ID == "" {
		return "No question loaded yet."
	}
	lines := []string{task.Prompt, fmt.Sprintf("Topic: %s", task.Topic.Title), fmt.Sprintf("Area: %s", task.Topic.Area), fmt.Sprintf("Kind: %s", task.Kind)}
	if len(task.Choices) > 0 {
		lines = append(lines, "Choices:")
		for i, c := range task.Choices {
			lines = append(lines, fmt.Sprintf(" %d) %s", i+1, c))
		}
	}
	return strings.Join(lines, "	")
}

func renderFeedback(r domain.AnswerResult, task domain.Question, ans string) string {
	status := "Wrong"
	switch r.ScoreDelta {
	case 10:
		status = "Correct"
	case 5:
		status = "Partially correct"
	}
	return fmt.Sprintf("%s. Score +%d. Your answer: %s. Expected: %s", status, r.ScoreDelta, ans, task.Expected)
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

	return fmt.Sprintf("Correct: %dWrong: %dScore: %dAttempts: %d	Weak topics:	%sAreas:	%s	Topic stats:	%s", s.correct, s.wrong, s.score, len(s.attempts), strings.Join(weak, "	"), strings.Join(areas, "	"), strings.Join(topics, "	"))
}
