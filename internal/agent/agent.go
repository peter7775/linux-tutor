package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"linux-tutor/internal/domain"
	"os"
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

func (a Agent) Generate(topicCode string) domain.ShellTask {
	switch topicCode {
	case "103.5":
		return domain.ShellTask{TopicCode: topicCode, Prompt: "Zobraz běžící procesy.", Expected: "ps", Hint: "LPIC 103.5: processes"}
	case "103.4":
		return domain.ShellTask{TopicCode: topicCode, Prompt: "Přesměruj chybový výstup do souboru error.log.", Expected: "cmd 2> error.log", Hint: "LPIC 103.4: redirection"}
	case "104.5":
		return domain.ShellTask{TopicCode: topicCode, Prompt: "Nastav oprávnění souboru na 640.", Expected: "chmod 640 file", Hint: "LPIC 104.5: permissions"}
	default:
		return domain.ShellTask{TopicCode: "103.1", Prompt: "Zobraz pracovní adresář.", Expected: "pwd", Hint: "LPIC command line basics"}
	}
}

func (a Agent) NextFromCatalog(after string) (domain.Topic, error) {
	for i, t := range a.Catalog.Topics {
		if t.Code == after && i+1 < len(a.Catalog.Topics) { return a.Catalog.Topics[i+1], nil }
	}
	if len(a.Catalog.Topics) > 0 { return a.Catalog.Topics[0], nil }
	return domain.Topic{}, errors.New("catalog empty")
}

func (a Agent) Evaluate(task domain.ShellTask, input string) (bool, string) {
	if input == task.Expected { return true, fmt.Sprintf("Správně (%s).", task.TopicCode) }
	return false, fmt.Sprintf("Špatně (%s). Správná odpověď je: %s", task.TopicCode, task.Expected)
}

func defaultTopics() []domain.Topic { return []domain.Topic{{Code: "103.5", Title: "Create, monitor and kill processes", Area: "GNU and Unix Commands"}} }
