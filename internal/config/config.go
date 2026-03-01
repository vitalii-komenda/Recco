package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SteamAPIKey string
	SteamID     string
	GeminiKey   string
	GeminiModel string
}

func (c Config) Validate() error {
	if c.SteamAPIKey == "" {
		return fmt.Errorf("STEAM_API_KEY is required")
	}
	if c.SteamID == "" {
		return fmt.Errorf("STEAM_ID is required")
	}
	if c.GeminiKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is required")
	}
	return nil
}

func Load() Config {
	_ = godotenv.Load()

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-1.5-flash"
	}

	return Config{
		SteamAPIKey: os.Getenv("STEAM_API_KEY"),
		SteamID:     os.Getenv("STEAM_ID"),
		GeminiKey:   os.Getenv("GEMINI_API_KEY"),
		GeminiModel: model,
	}
}
