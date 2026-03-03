package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
)

type Config struct {
	SteamAPIKey string
	SteamID     string
	GeminiKey   string
	GeminiModel string
}

var serviceName = "recco"

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

func (c Config) SaveInKeyring() error {
	if err := keyring.Set(serviceName, "STEAM_API_KEY", c.SteamAPIKey); err != nil {
		return fmt.Errorf("failed to save STEAM_API_KEY to keyring: %w", err)
	}
	if err := keyring.Set(serviceName, "STEAM_ID", c.SteamID); err != nil {
		return fmt.Errorf("failed to save STEAM_ID to keyring: %w", err)
	}
	if err := keyring.Set(serviceName, "GEMINI_API_KEY", c.GeminiKey); err != nil {
		return fmt.Errorf("failed to save GEMINI_API_KEY to keyring: %w", err)
	}

	return nil
}

func getFromEnvOrKeyring(name string) string {
	val := os.Getenv(name)
	if val != "" {
		return val
	}

	val, err := keyring.Get(serviceName, name)
	if err != nil {
		fmt.Printf("keyring not found: %s \n", err.Error())
		return ""
	}
	return val
}

func Load() Config {
	_ = godotenv.Load()

	model := getFromEnvOrKeyring("GEMINI_MODEL")
	if model == "" {
		model = "gemini-flash-latest"
	}

	steamAPIKey := getFromEnvOrKeyring("STEAM_API_KEY")
	steamID := getFromEnvOrKeyring("STEAM_ID")
	geminiKey := getFromEnvOrKeyring("GEMINI_API_KEY")

	return Config{
		SteamAPIKey: steamAPIKey,
		SteamID:     steamID,
		GeminiKey:   geminiKey,
		GeminiModel: model,
	}
}
