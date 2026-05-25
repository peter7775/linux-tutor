package ai

import (
	"encoding/json"
	"fmt"
	"os"

	"linux-tutor/internal/domain"
)

type TopicCatalog struct {
	Tracks []struct {
		ID    string   `json:"id"`
		Title string   `json:"title"`
		Exams []string `json:"exams"`
	} `json:"tracks"`
	Topics []struct {
		Track       string   `json:"track"`
		Exam        string   `json:"exam"`
		Code        string   `json:"code"`
		Area        string   `json:"area"`
		Title       string   `json:"title"`
		Difficulty  string   `json:"difficulty"`
		HintStyle   string   `json:"hint_style"`
		Keywords    []string `json:"keywords"`
		TaskTypes   []string `json:"task_types"`
		PrereqCodes []string `json:"prerequisites"`
	} `json:"topics"`
}

func LoadCatalog(path string) ([]domain.Topic, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tc TopicCatalog
	if err := json.Unmarshal(b, &tc); err != nil {
		return nil, fmt.Errorf("decode catalog: %w", err)
	}
	out := make([]domain.Topic, 0, len(tc.Topics))
	for _, t := range tc.Topics {
		out = append(out, domain.Topic{Code: t.Code, Title: t.Title, Area: t.Area})
	}
	return out, nil
}
