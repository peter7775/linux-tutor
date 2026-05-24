package ai

type AIConfig struct { Enabled bool; Model string; APIKey string; BaseURL string }
func DefaultAIConfig() AIConfig { return AIConfig{Enabled:true, Model:"claude-sonnet-4-6", BaseURL:"https://api.anthropic.com"} }
