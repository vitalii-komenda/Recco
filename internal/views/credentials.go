package views

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"steam-bubbles/internal/config"
)

func CredentialFields(cfg *config.Config) []huh.Field {
	var fields []huh.Field

	if cfg.SteamAPIKey == "" {
		fields = append(fields,
			huh.NewInput().
				Title("Steam API Key").
				Description("Get yours at https://steamcommunity.com/dev/apikey").
				EchoMode(huh.EchoModePassword).
				Value(&cfg.SteamAPIKey),
		)
	}
	if cfg.SteamID == "" {
		fields = append(fields,
			huh.NewInput().
				Title("Steam ID (64-bit)").
				Description("e.g. 76561198xxxxxxxxx  –  find yours at steamid.io").
				Value(&cfg.SteamID),
		)
	}
	if cfg.GeminiKey == "" {
		fields = append(fields,
			huh.NewInput().
				Title("Gemini API Key").
				Description("Get yours at https://aistudio.google.com/app/apikey").
				EchoMode(huh.EchoModePassword).
				Value(&cfg.GeminiKey),
		)
	}

	return fields
}

func RunCredentialForm(cfg *config.Config) error {
	fields := CredentialFields(cfg)
	if len(fields) == 0 {
		return nil
	}

	fmt.Println(TitleStyle.Render(" 🎮  Recco "))

	return huh.NewForm(huh.NewGroup(fields...)).Run()
}
