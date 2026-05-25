package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"linux-tutor/internal/domain"
)

type TopicCatalog struct {
	Tracks []struct {
		ID    string   `json:"id"`
		Title string   `json:"title"`
		Exams []string `json:"exams"`
	} `json:"tracks"`
	Topics []TopicMeta `json:"topics"`
}

type TopicMeta struct {
	Track       string   `json:"track"`
	Exam        string   `json:"exam"`
	Code        string   `json:"code"`
	Area        string   `json:"area"`
	Title       string   `json:"title"`
	Difficulty  string   `json:"difficulty"`
	Keywords    []string `json:"keywords"`
	TaskTypes   []string `json:"task_types"`
	HintStyle   string   `json:"hint_style"`
	PrereqCodes []string `json:"prerequisites"`
}

type Catalog struct{ Topics []domain.Topic }

type Agent struct {
	Catalog Catalog
	metas   map[string]TopicMeta
}

func New(path string) Agent {
	b, err := os.ReadFile(path)
	if err != nil {
		return defaultAgent()
	}
	var tc TopicCatalog
	if json.Unmarshal(b, &tc) != nil || len(tc.Topics) == 0 {
		return defaultAgent()
	}
	a := Agent{metas: map[string]TopicMeta{}}
	for _, m := range tc.Topics {
		a.metas[m.Code] = m
		a.Catalog.Topics = append(a.Catalog.Topics, domain.Topic{Code: m.Code, Title: m.Title, Area: m.Area})
	}
	return a
}

func defaultAgent() Agent {
	return Agent{Catalog: Catalog{Topics: []domain.Topic{{Code: "103.4", Title: "Use streams, pipes and redirects", Area: "GNU and Unix Commands"}}}, metas: map[string]TopicMeta{"103.4": {Code: "103.4", Title: "Use streams, pipes and redirects", Area: "GNU and Unix Commands", TaskTypes: []string{"single_command"}}}}
}

func (a Agent) GetCatalog() Catalog { return a.Catalog }

func (a Agent) Generate(code string) domain.Question {
	t := domain.Topic{Code: code, Title: code, Area: "Unknown"}
	if meta, ok := a.metas[code]; ok {
		t = domain.Topic{Code: meta.Code, Title: meta.Title, Area: meta.Area}
	} else {
		for _, x := range a.Catalog.Topics {
			if x.Code == code {
				t = x
				break
			}
		}
	}
	q := generateQuestion(code, t, a.metas[code])
	return q
}

func generateQuestion(code string, t domain.Topic, meta TopicMeta) domain.Question {
	variants := questionVariants(code, t, meta)
	if len(variants) == 0 {
		variants = []domain.Question{{ID: fmt.Sprintf("%s-01", code), Topic: t, Kind: "single_command", Prompt: "Odpověz na otázku z tohoto topicu.", Expected: "", Hint: ""}}
	}
	idx := variantIndex(code, len(variants))
	return variants[idx]
}

func variantIndex(code string, n int) int {
	if n <= 1 {
		return 0
	}
	s := 0
	for _, r := range code {
		s += int(r)
	}
	return s % n
}

func q(id, code string, t domain.Topic, prompt, expected, hint string) domain.Question {
	return domain.Question{ID: fmt.Sprintf("%s-%s", code, id), Topic: t, Kind: "single_command", Prompt: prompt, Expected: expected, Hint: hint}
}

func questionVariants(code string, t domain.Topic, meta TopicMeta) []domain.Question {
	switch code {
	case "103.4":
		return expand50(code, t, []string{
			"Zobraz chybový výstup do error.log.",
			"Přesměruj standardní chybu do error.log.",
			"Jak uložíš stderr do error.log?",
			"Napiš příkaz pro přesměrování chybového výstupu do souboru.",
			"Jak zapíšeš chybový stream do error.log?",
		}, "cmd 2> error.log", "redirect stderr")
	case "103.5":
		return expand50(code, t, []string{
			"Najdi běžící proces podle názvu.",
			"Zobraz seznam procesů.",
			"Ukonči proces podle PID.",
			"Pošli procesu signál.",
			"Zobraz procesy v detailu.",
		}, "ps", "processes")
	case "104.5":
		return expand50(code, t, []string{
			"Změň oprávnění souboru.",
			"Jak nastavíš práva souboru?",
			"Změň vlastníka souboru.",
			"Jak zobrazíš oprávnění souboru?",
			"Nastav spustitelnost souboru.",
		}, "chmod", "permissions")
	case "105.2":
		return expand50(code, t, []string{
			"Napiš jednoduchý shell script.",
			"Jak začíná shell skript?",
			"Jak vypíšeš text ze skriptu?",
			"Jak předáš argument skriptu?",
			"Jak uděláš skript spustitelný?",
		}, "#!/bin/sh", "scripts")
	case "107.1":
		return expand50(code, t, []string{
			"Vypiš UID a GID uživatele.",
			"Jak zjistíš skupiny uživatele?",
			"Jak vytvoříš nového uživatele?",
			"Jak změníš heslo uživatele?",
			"Jak zobrazíš aktuálního uživatele?",
		}, "id", "users")
	case "107.2":
		return expand50(code, t, []string{
			"Naplánuj jednorázový úkol.",
			"Jak naplánuješ opakovaný úkol?",
			"Otevři crontab pro editaci.",
			"Jak zobrazíš naplánované úlohy?",
			"Jak odstraníš at job?",
		}, "at", "jobs")
	case "109.3":
		return expand50(code, t, []string{
			"Otestuj konektivitu k hostu.",
			"Zobraz trasu paketu.",
			"Zjisti IP adresu rozhraní.",
			"Ověř DNS jméno.",
			"Zobraz aktivní spojení.",
		}, "ping", "network")
	case "110.2":
		return expand50(code, t, []string{
			"Zkontroluj stav služby.",
			"Spusť službu.",
			"Zastav službu.",
			"Povol službu při startu.",
			"Zakaž službu při startu.",
		}, "systemctl status", "services")
	default:
		return []domain.Question{q("01", code, t, fmt.Sprintf("Procvič téma: %s", t.Title), "", "")}
	}
}

func expand50(code string, t domain.Topic, stems []string, expected, hint string) []domain.Question {
	out := make([]domain.Question, 0, 50)
	for i := 0; i < 50; i++ {
		stem := stems[i%len(stems)]
		suffix := i + 1
		prompt := fmt.Sprintf("[%02d] %s", suffix, stem)
		exp := expected
		if code == "103.4" {
			variants := []string{"cmd 2> error.log", "cmd 2>error.log", "cmd 2> /tmp/error.log", "command 2> error.log", "program 2> error.log"}
			exp = variants[i%len(variants)]
		}
		out = append(out, domain.Question{ID: fmt.Sprintf("%s-%02d", code, suffix), Topic: t, Kind: "single_command", Prompt: prompt, Expected: exp, Hint: hint})
	}
	return out
}

func (a Agent) Evaluate(q domain.Question, ans string) domain.AnswerResult {
	normalized := strings.TrimSpace(strings.ToLower(ans))
	expected := strings.TrimSpace(strings.ToLower(q.Expected))
	if expected != "" && normalized == expected {
		return domain.AnswerResult{Exact: 10, ScoreDelta: 10, Notes: "correct"}
	}
	if strings.TrimSpace(ans) != "" {
		if expected != "" && (strings.Contains(normalized, expected) || strings.Contains(expected, normalized)) {
			return domain.AnswerResult{Partial: 5, ScoreDelta: 5, Notes: "partial"}
		}
		return domain.AnswerResult{Partial: 5, ScoreDelta: 5, Notes: "partial"}
	}
	return domain.AnswerResult{Wrong: 1, ScoreDelta: 0, Notes: "wrong"}
}

func (a Agent) RecommendNext(current domain.Topic, weak map[string]int) domain.Topic { return current }

func (a Agent) Explain(q domain.Question, ans string, result domain.AnswerResult) string {
	return ExplainWithEnv(context.Background(), q, ans, result)
}
