package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"linux-tutor/internal/domain"
)

type GitHubModelsClient struct {
	Token   string
	Model   string
	Version string
	BaseURL string
	HTTP    *http.Client
}

func NewGitHubModelsClient(token, model, version string) *GitHubModelsClient {
	if version == "" {
		version = "2026-03-10"
	}
	return &GitHubModelsClient{
		Token:   token,
		Model:   model,
		Version: version,
		BaseURL: "https://models.github.ai/inference",
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *GitHubModelsClient) Explain(ctx context.Context, task domain.Question, answer string, result domain.AnswerResult) string {
	if strings.TrimSpace(c.Token) == "" {
		return localExplain(task, answer, result)
	}

	prompt := fmt.Sprintf(`Vysvětli stručně a jasně odpověď na otázku z Linux/LPIC.

Téma: %s
Otázka: %s
Odpověď uživatele: %s
Výsledek: %d
Nápověda: %s
Správná odpověď / očekávání: %s

Napiš 2-4 krátké věty v češtině.`, task.Topic.Title, task.Prompt, answer, result.ScoreDelta, task.Hint, task.Expected)

	body := map[string]any{
		"model":       c.Model,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"temperature": 0.2,
		"max_tokens":  220,
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/chat/completions", bytes.NewReader(b))
	if err != nil {
		return localExplain(task, answer, result)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-GitHub-Api-Version", c.Version)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return localExplain(task, answer, result)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return localExplain(task, answer, result)
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return localExplain(task, answer, result)
	}
	if len(parsed.Choices) == 0 {
		return localExplain(task, answer, result)
	}
	text := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if text == "" {
		return localExplain(task, answer, result)
	}
	return text
}

func localExplain(task domain.Question, answer string, result domain.AnswerResult) string {
	if task.ID == "" {
		return "No explanation available."
	}
	if result.ScoreDelta > 0 {
		return fmt.Sprintf("Správně. Klíč je v tématu %s. %s", task.Topic.Title, firstSentence(task.Hint, task.Expected, "Dobrá práce."))
	}
	return fmt.Sprintf("Ne úplně. Správný směr je: %s. %s", nonEmpty(task.Expected, "zkus si projít nápovědu"), firstSentence(task.Hint, "Pamatuj si základní princip.", ""))
}

func firstSentence(parts ...string) string {
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			return p
		}
	}
	return ""
}

func nonEmpty(parts ...string) string {
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			return p
		}
	}
	return ""
}

func ExplainWithEnv(ctx context.Context, task domain.Question, answer string, result domain.AnswerResult) string {
	c := NewGitHubModelsClient(os.Getenv("GITHUB_TOKEN"), os.Getenv("GITHUB_MODEL"), os.Getenv("GITHUB_API_VERSION"))
	return c.Explain(ctx, task, answer, result)
}
