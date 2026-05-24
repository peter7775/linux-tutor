package ai

import ("encoding/json"; "fmt"; "os"; "linux-tutor/internal/domain")

type TopicCatalog struct { Tracks []struct{ ID, Title string; Exams []string } `json:"tracks"`; Topics []struct{ Track, Exam, Code, Area, Title, Difficulty, HintStyle string; Keywords, TaskTypes, PrereqCodes []string } `json:"topics"` }
func LoadCatalog(path string) ([]domain.Topic, error) { b, err := os.ReadFile(path); if err != nil { return nil, err }; var tc TopicCatalog; if err := json.Unmarshal(b, &tc); err != nil { return nil, fmt.Errorf("decode catalog: %w", err) }; out := make([]domain.Topic, 0, len(tc.Topics)); for _, t := range tc.Topics { out = append(out, domain.Topic{Code:t.Code, Title:t.Title, Area:t.Area}) }; return out, nil }
