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

func New(path string) Agent { b, err := os.ReadFile(path); if err != nil { return Agent{Catalog: Catalog{Topics: defaultTopics()}} }; var c Catalog; if json.Unmarshal(b,&c)!=nil || len(c.Topics)==0 { c.Topics=defaultTopics() }; return Agent{Catalog:c} }
func (a Agent) topic(code string) domain.Topic { for _, t := range a.Catalog.Topics { if t.Code == code { return t } }; return domain.Topic{Code: code, Title: "Unknown", Area: "Unknown"} }

func (a Agent) Generate(code string) domain.Task {
	t := a.topic(code)
	m := map[string][]domain.Task{
		"103.4": {{ID:"103.4-1", Topic:t, Kind:"single_command", Prompt:"Zobraz chybový výstup do error.log.", Expected:"cmd 2> error.log", Hint:"LPIC 103.4: redirect stderr"}, {ID:"103.4-2", Topic:t, Kind:"multi_command", Prompt:"Spoj standardní výstup a chybový výstup do výsledku.txt.", Expected:"cmd > result.txt 2>&1", Hint:"LPIC 103.4: redirection"}},
		"103.5": {{ID:"103.5-1", Topic:t, Kind:"scenario", Prompt:"Proces se zasekl. Jaký příkaz použiješ pro zobrazení PID a stavu?", Expected:"ps", Hint:"LPIC 103.5: processes"}, {ID:"103.5-2", Topic:t, Kind:"single_command", Prompt:"Ukonči proces podle PID 1234.", Expected:"kill 1234", Hint:"LPIC 103.5: kill"}},
		"104.5": {{ID:"104.5-1", Topic:t, Kind:"fill_blank", Prompt:"Doplň příkaz pro nastavení práv 640: ____ 640 file", Expected:"chmod", Hint:"LPIC 104.5: permissions"}, {ID:"104.5-2", Topic:t, Kind:"multiple_choice", Prompt:"Které oprávnění dává vlastníkovi čtení, zápis a spouštění?", Choices:[]string{"700","644","755"}, Expected:"700", Hint:"LPIC 104.5: permissions"}},
		"105.2": {{ID:"105.2-1", Topic:t, Kind:"ordering", Prompt:"Seřaď kroky: [write script, shebang, make executable]", Expected:"write script > shebang > make executable", Hint:"LPIC 105.2: scripts"}, {ID:"105.2-2", Topic:t, Kind:"fill_blank", Prompt:"Doplň proměnnou pro první argument skriptu: echo $____", Expected:"1", Hint:"LPIC 105.2: positional parameters"}},
		"107.1": {{ID:"107.1-1", Topic:t, Kind:"multiple_choice", Prompt:"Který soubor obsahuje uživatelské účty?", Choices:[]string{"/etc/passwd","/etc/hosts","/etc/fstab"}, Expected:"/etc/passwd", Hint:"LPIC 107.1: user accounts"}, {ID:"107.1-2", Topic:t, Kind:"scenario", Prompt:"Přidej uživatele do skupiny wheel. Jaké dva příkazy použiješ?", Expected:"usermod;groupadd", Hint:"LPIC 107.1: users and groups"}},
		"107.2": {{ID:"107.2-1", Topic:t, Kind:"single_command", Prompt:"Naplánuj jednorázovou úlohu přes cron pro uživatele.", Expected:"crontab -e", Hint:"LPIC 107.2: cron"}, {ID:"107.2-2", Topic:t, Kind:"scenario", Prompt:"Spusť skript každý den v 06:30. Jaký zápis používá crontab?", Expected:"30 6 * * *", Hint:"LPIC 107.2: scheduling"}},
		"109.3": {{ID:"109.3-1", Topic:t, Kind:"multi_command", Prompt:"Nejdřív zjisti IP, pak otestuj spojení. Odděl příkazy ;", Expected:"ip a;ping", Hint:"LPIC 109.3: troubleshooting"}, {ID:"109.3-2", Topic:t, Kind:"scenario", Prompt:"DNS nefunguje. Který soubor nejprve zkontroluješ?", Expected:"/etc/resolv.conf", Hint:"LPIC 109.3: DNS troubleshooting"}},
		"110.2": {{ID:"110.2-1", Topic:t, Kind:"single_command", Prompt:"Zobraz pravidla firewallu na hostu.", Expected:"iptables -L", Hint:"LPIC 110.2: host security"}, {ID:"110.2-2", Topic:t, Kind:"multiple_choice", Prompt:"Který příkaz typicky zobrazí stav SSH služby?", Choices:[]string{"systemctl status ssh","grep ssh /etc/passwd","ls /root"}, Expected:"systemctl status ssh", Hint:"LPIC 110.2: host security"}},
	}
	list := m[code]; if len(list)==0 { return domain.Task{ID:"t0", Topic:t, Kind:"single_command", Prompt:"Zobraz pracovní adresář.", Expected:"pwd", Hint:"LPIC basics"} }
	return list[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(list))]
}

func (a Agent) Next(code string) (domain.Topic, error) { for i, t := range a.Catalog.Topics { if t.Code == code && i+1 < len(a.Catalog.Topics) { return a.Catalog.Topics[i+1], nil } }; if len(a.Catalog.Topics)==0 { return domain.Topic{}, errors.New("catalog empty") }; return a.Catalog.Topics[rand.Intn(len(a.Catalog.Topics))], nil }

func (a Agent) Evaluate(task domain.Task, input string) (bool, string) {
	norm := strings.TrimSpace(strings.ToLower(input)); exp := strings.TrimSpace(strings.ToLower(task.Expected)); ok := false
	switch task.Kind {
	case "single_command", "fill_blank": ok = norm == exp || strings.HasPrefix(norm, exp+" ")
	case "multi_command": ok = strings.Contains(norm, "ip a") && strings.Contains(norm, "ping") || strings.Contains(norm, "2>&1")
	case "multiple_choice": ok = norm == exp
	case "ordering": ok = norm == exp
	case "scenario": ok = strings.Contains(norm, exp)
	default: ok = norm == exp
	}
	if ok { return true, fmt.Sprintf("Správně (%s / %s).", task.Topic.Code, task.Kind) }
	if len(task.Choices) > 0 { return false, fmt.Sprintf("Špatně. Správná odpověď je: %s", task.Expected) }
	return false, fmt.Sprintf("Špatně. Očekávám: %s", task.Expected)
}

func defaultTopics() []domain.Topic { return []domain.Topic{{Code: "103.4", Title: "Use streams, pipes and redirects", Area: "GNU and Unix Commands"}} }
