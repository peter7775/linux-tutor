package agent

import (
	"encoding/json"
	"errors"
	"linux-tutor/internal/domain"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Catalog struct {
	Topics []domain.Topic `json:"topics"`
}
type Agent struct{ Catalog Catalog }

func New(path string) Agent {
	b, err := os.ReadFile(path)
	if err != nil {
		return Agent{Catalog: Catalog{Topics: defaultTopics()}}
	}
	var c Catalog
	if json.Unmarshal(b, &c) != nil || len(c.Topics) == 0 {
		c.Topics = defaultTopics()
	}
	return Agent{Catalog: c}
}
func (a Agent) topic(code string) domain.Topic {
	for _, t := range a.Catalog.Topics {
		if t.Code == code {
			return t
		}
	}
	return domain.Topic{Code: code, Title: "Unknown", Area: "Unknown"}
}
func (a Agent) Generate(code string) domain.Task {
	t := a.topic(code)
	switch code {
	case "103.4":
		return domain.Task{ID: "103.4-1", Topic: t, Kind: "single_command", Prompt: "Zobraz chybový výstup do error.log.", Expected: "cmd 2> error.log", Hint: "LPIC 103.4: redirect stderr"}
	case "103.5":
		return domain.Task{ID: "103.5-1", Topic: t, Kind: "scenario", Prompt: "Proces se zasekl. Jaký příkaz použiješ pro zobrazení PID a stavu?", Expected: "ps", Hint: "LPIC 103.5: processes"}
	case "104.5":
		return domain.Task{ID: "104.5-1", Topic: t, Kind: "fill_blank", Prompt: "Doplň příkaz pro nastavení práv 640: ____ 640 file", Expected: "chmod", Hint: "LPIC 104.5: permissions"}
	case "105.2":
		return domain.Task{ID: "105.2-1", Topic: t, Kind: "ordering", Prompt: "Seřaď kroky: [write script, shebang, make executable]", Expected: "write script > shebang > make executable", Hint: "LPIC 105.2: scripts"}
	case "107.1":
		return domain.Task{ID: "107.1-1", Topic: t, Kind: "multiple_choice", Prompt: "Který soubor obsahuje uživatelské účty?", Choices: []string{"/etc/passwd", "/etc/hosts", "/etc/fstab"}, Expected: "/etc/passwd", Hint: "LPIC 107.1: user accounts"}
	case "107.2":
		return domain.Task{ID: "107.2-1", Topic: t, Kind: "scenario", Prompt: "Spusť skript každý den v 06:30. Jaký zápis používá crontab?", Expected: "30 6 * * *", Hint: "LPIC 107.2: scheduling"}
	case "109.3":
		return domain.Task{ID: "109.3-1", Topic: t, Kind: "multi_command", Prompt: "Nejdřív zjisti IP, pak otestuj spojení. Odděl příkazy ;", Expected: "ip a;ping", Hint: "LPIC 109.3: troubleshooting"}
	case "110.2":
		return domain.Task{ID: "110.2-1", Topic: t, Kind: "multiple_choice", Prompt: "Který příkaz typicky zobrazí stav SSH služby?", Choices: []string{"systemctl status ssh", "grep ssh /etc/passwd", "ls /root"}, Expected: "systemctl status ssh", Hint: "LPIC 110.2: host security"}
	default:
		return domain.Task{ID: "t0", Topic: t, Kind: "single_command", Prompt: "Zobraz pracovní adresář.", Expected: "pwd", Hint: "LPIC basics"}
	}
}
func (a Agent) Next(code string) (domain.Topic, error) {
	for i, t := range a.Catalog.Topics {
		if t.Code == code && i+1 < len(a.Catalog.Topics) {
			return a.Catalog.Topics[i+1], nil
		}
	}
	if len(a.Catalog.Topics) == 0 {
		return domain.Topic{}, errors.New("catalog empty")
	}
	return a.Catalog.Topics[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(a.Catalog.Topics))], nil
}
func (a Agent) Evaluate(task domain.Task, input string) domain.AnswerResult {
	norm := strings.TrimSpace(strings.ToLower(input))
	exp := strings.TrimSpace(strings.ToLower(task.Expected))
	r := domain.AnswerResult{Exact: 10, Partial: 5, Wrong: 0}
	switch task.Kind {
	case "single_command", "fill_blank":
		if norm == exp || strings.HasPrefix(norm, exp+" ") {
			r.ScoreDelta = r.Exact
			r.Notes = "correct"
		} else {
			r.ScoreDelta = r.Wrong
			r.Notes = "wrong"
		}
	case "multi_command":
		if (strings.Contains(norm, "ip a") && strings.Contains(norm, "ping")) || strings.Contains(norm, "2>&1") {
			r.ScoreDelta = r.Exact
			r.Notes = "correct"
		} else if strings.Contains(norm, "ip") || strings.Contains(norm, "ping") {
			r.ScoreDelta = r.Partial
			r.Notes = "partial"
		} else {
			r.ScoreDelta = r.Wrong
			r.Notes = "wrong"
		}
	case "multiple_choice":
		if norm == exp {
			r.ScoreDelta = r.Exact
			r.Notes = "correct"
		} else {
			r.ScoreDelta = r.Wrong
			r.Notes = "wrong"
		}
	case "ordering":
		if norm == exp {
			r.ScoreDelta = r.Exact
			r.Notes = "correct"
		} else if len(norm) > 0 {
			r.ScoreDelta = r.Partial
			r.Notes = "partial"
		} else {
			r.ScoreDelta = r.Wrong
			r.Notes = "wrong"
		}
	case "scenario":
		if strings.Contains(norm, exp) {
			r.ScoreDelta = r.Exact
			r.Notes = "correct"
		} else if len(norm) > 0 {
			r.ScoreDelta = r.Partial
			r.Notes = "partial"
		} else {
			r.ScoreDelta = r.Wrong
			r.Notes = "wrong"
		}
	default:
		if norm == exp {
			r.ScoreDelta = r.Exact
			r.Notes = "correct"
		} else {
			r.ScoreDelta = r.Wrong
			r.Notes = "wrong"
		}
	}
	return r
}
func defaultTopics() []domain.Topic {
	return []domain.Topic{{Code: "103.4", Title: "Use streams, pipes and redirects", Area: "GNU and Unix Commands"}}
}
