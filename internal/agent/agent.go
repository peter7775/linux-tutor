package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"linux-tutor/internal/domain"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Catalog struct { Topics []domain.Topic `json:"topics"` }

type Agent struct { Catalog Catalog }

func New(path string) Agent {
	b, err := os.ReadFile(path)
	if err != nil { return Agent{Catalog: Catalog{Topics: defaultTopics()}} }
	var c Catalog
	if json.Unmarshal(b, &c) != nil || len(c.Topics) == 0 { c.Topics = defaultTopics() }
	return Agent{Catalog: c}
}

func (a Agent) topic(code string) domain.Topic {
	for _, t := range a.Catalog.Topics { if t.Code == code { return t } }
	return domain.Topic{Code: code, Title: "Unknown", Area: "Unknown"}
}

func (a Agent) Generate(code string) domain.Task {
	t := a.topic(code)
	switch code {
	case "103.4":
		return domain.Task{ID: "t1", Topic: t, Kind: "single_command", Prompt: "Zobraz chybový výstup příkazu do error.log.", Expected: "cmd 2> error.log", Hint: "LPIC 103.4: redirect stderr"}
	case "103.5":
		return domain.Task{ID: "t2", Topic: t, Kind: "scenario", Prompt: "Proces se zasekl. Jaký příkaz použiješ pro zobrazení PID a stavu?", Expected: "ps", Hint: "LPIC 103.5: processes"}
	case "104.5":
		return domain.Task{ID: "t3", Topic: t, Kind: "fill_blank", Prompt: "Doplň příkaz pro nastavení práv 640: ____ 640 file", Expected: "chmod", Hint: "LPIC 104.5: permissions"}
	case "105.2":
		return domain.Task{ID: "t4", Topic: t, Kind: "ordering", Prompt: "Seřaď kroky shell skriptu: [shebang, make executable, write script]", Expected: "write script > shebang > make executable", Hint: "LPIC 105.2: simple scripts"}
	case "107.1":
		return domain.Task{ID: "t5", Topic: t, Kind: "multiple_choice", Prompt: "Který soubor typicky obsahuje informace o uživatelích?", Choices: []string{"/etc/passwd", "/etc/hosts", "/var/log/syslog"}, Expected: "/etc/passwd", Hint: "LPIC 107.1: user accounts"}
	case "109.3":
		return domain.Task{ID: "t6", Topic: t, Kind: "multi_command", Prompt: "Jeden příkaz nestačí: nejdřív zjisti IP, potom otestuj spojení. Napiš obě části oddělené znakem ;", Expected: "ip a; ping", Hint: "LPIC 109.3: troubleshooting"}
	default:
		return domain.Task{ID: "t0", Topic: t, Kind: "single_command", Prompt: "Zobraz pracovní adresář.", Expected: "pwd", Hint: "LPIC command line basics"}
	}
}

func (a Agent) Next(code string) (domain.Topic, error) {
	for i, t := range a.Catalog.Topics { if t.Code == code && i+1 < len(a.Catalog.Topics) { return a.Catalog.Topics[i+1], nil } }
	if len(a.Catalog.Topics) == 0 { return domain.Topic{}, errors.New("catalog empty") }
	return a.Catalog.Topics[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(a.Catalog.Topics))], nil
}

func (a Agent) Evaluate(task domain.Task, input string) (bool, string) {
	norm := strings.TrimSpace(strings.ToLower(input))
	exp := strings.TrimSpace(strings.ToLower(task.Expected))
	ok := false
	switch task.Kind {
	case "single_command", "fill_blank":
		ok = norm == exp || strings.HasPrefix(norm, exp+" ")
	case "multi_command":
		ok = strings.Contains(norm, "ip a") && strings.Contains(norm, "ping")
	case "multiple_choice":
		ok = norm == exp
	case "ordering":
		ok = norm == exp
	case "scenario":
		ok = strings.Contains(norm, exp)
	default:
		ok = norm == exp
	}
	if ok { return true, fmt.Sprintf("Správně (%s / %s).", task.Topic.Code, task.Kind) }
	if task.Kind == "multiple_choice" {
		return false, fmt.Sprintf("Špatně. Správná odpověď je: %s", task.Expected)
	}
	return false, fmt.Sprintf("Špatně. Očekávám: %s", task.Expected)
}

func defaultTopics() []domain.Topic { return []domain.Topic{{Code: "103.4", Title: "Use streams, pipes and redirects", Area: "GNU and Unix Commands"}} }
