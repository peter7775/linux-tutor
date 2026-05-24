package agent

import (
	"encoding/json"
	"os"
	"strings"

	"linux-tutor/internal/domain"
)

type Catalog struct{ Topics []domain.Topic }
type Agent struct{ Catalog Catalog }

type lpicCatalog struct {
	Topics []struct {
		Code  string `json:"code"`
		Title string `json:"title"`
		Area  string `json:"area"`
	} `json:"topics"`
}

func New(path string) Agent {
	b, err := os.ReadFile(path)
	if err != nil {
		return Agent{Catalog: Catalog{Topics: []domain.Topic{}}}
	}
	var c lpicCatalog
	if err := json.Unmarshal(b, &c); err != nil {
		return Agent{Catalog: Catalog{Topics: []domain.Topic{}}}
	}
	tops := make([]domain.Topic, 0, len(c.Topics))
	for _, t := range c.Topics {
		tops = append(tops, domain.Topic{Code: strings.TrimSpace(t.Code), Title: strings.TrimSpace(t.Title), Area: strings.TrimSpace(t.Area)})
	}
	return Agent{Catalog: Catalog{Topics: tops}}
}

func (a Agent) GetCatalog() Catalog { return a.Catalog }

func (a Agent) Generate(code string) domain.Question {
	t := domain.Topic{Code: code, Title: code, Area: "Unknown"}
	for _, x := range a.Catalog.Topics {
		if x.Code == code {
			t = x
			break
		}
	}
	return domain.Question{ID: code + "-1", Topic: t, Kind: "single_command", Prompt: promptFor(code), Expected: expectedFor(code), Hint: hintFor(code)}
}

func (a Agent) Evaluate(q domain.Question, ans string) domain.AnswerResult {
	return ScoreAnswer(q.Expected, ans)
}

func (a Agent) RecommendNext(current domain.Topic, weak map[string]int) domain.Topic { return current }

func promptFor(code string) string {
	switch code {
	case "103.4":
		return "Zobraz pracovní adresář."
	case "103.5":
		return "Najdi a ukonči proces."
	case "104.5":
		return "Změň oprávnění souboru."
	case "105.2":
		return "Skript vypiš pozdrav."
	case "107.1":
		return "Vypiš informace o uživateli."
	case "107.2":
		return "Naplánuj jednorázový úkol."
	case "109.3":
		return "Otestuj síťové spojení."
	case "110.2":
		return "Zkontroluj stav služby."
	default:
		return "Odpověz na otázku z tohoto tématu."
	}
}

func expectedFor(code string) string {
	switch code {
	case "103.4":
		return "pwd"
	case "103.5":
		return "ps"
	case "104.5":
		return "chmod"
	case "105.2":
		return "shell script"
	case "107.1":
		return "id"
	case "107.2":
		return "at"
	case "109.3":
		return "ping"
	case "110.2":
		return "systemctl"
	default:
		return ""
	}
}

func hintFor(code string) string {
	switch code {
	case "103.4":
		return "basic shell"
	case "103.5":
		return "processes"
	case "104.5":
		return "permissions"
	case "105.2":
		return "scripts"
	case "107.1":
		return "users"
	case "107.2":
		return "jobs"
	case "109.3":
		return "network"
	case "110.2":
		return "services"
	default:
		return ""
	}
}
