package agent

import (
	"strings"
	"linux-tutor/internal/domain"
)

type Agent struct{}

func New() Agent { return Agent{} }

func (a Agent) Generate(topic string) domain.ShellTask {
	switch topic {
	case "103.5":
		return domain.ShellTask{TopicCode: topic, Prompt: "Zobraz běžící procesy.", Expected: "ps", Hint: "LPIC 103.5: processes"}
	case "103.4":
		return domain.ShellTask{TopicCode: topic, Prompt: "Přesměruj chybový výstup do souboru error.log.", Expected: "cmd 2> error.log", Hint: "LPIC 103.4: redirection"}
	case "104.5":
		return domain.ShellTask{TopicCode: topic, Prompt: "Nastav oprávnění souboru na čitelné pro vlastníka a skupinu, ostatní bez přístupu.", Expected: "chmod 640 file", Hint: "LPIC 104.5: permissions"}
	default:
		return domain.ShellTask{TopicCode: "103.1", Prompt: "Zobraz aktuální pracovní adresář.", Expected: "pwd", Hint: "LPIC command line basics"}
	}
}

func (a Agent) Evaluate(task domain.ShellTask, input string) (bool, string) {
	norm := strings.TrimSpace(strings.ToLower(input))
	exp := strings.TrimSpace(strings.ToLower(task.Expected))
	if norm == exp { return true, "Správně." }
	if strings.Contains(exp, "ps") && strings.HasPrefix(norm, "ps") { return true, "Přijato." }
	return false, "Špatně. Správná odpověď je: " + task.Expected
}
