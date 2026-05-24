package config

import (
	"os"
	"strings"
)

type Config struct {
	LLMProvider   string
	GitHubToken   string
	GitHubModel   string
	GitHubVersion string
}

func Load() Config {
	return Config{
		LLMProvider:   env("LLM_PROVIDER", "github"),
		GitHubToken:   strings.TrimSpace(os.Getenv("GITHUB_TOKEN")),
		GitHubModel:   env("GITHUB_MODEL", "openai/gpt-4o-mini"),
		GitHubVersion: env("GITHUB_API_VERSION", "2026-03-10"),
	}
}

func env(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}
