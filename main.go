package main

import (
	"context"
	"fmt"
	"os"

	"steam-bubbles/internal/app"
	"steam-bubbles/internal/config"
	"steam-bubbles/internal/domain"
	infraGemini "steam-bubbles/internal/infrastructure/gemini"
	infraSteam "steam-bubbles/internal/infrastructure/steam"
	"steam-bubbles/internal/models"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	factory := func(cfg config.Config) (domain.GameRepository, domain.Recommender) {
		return infraSteam.NewClient(cfg.SteamAPIKey),
			infraGemini.NewClient(cfg.GeminiKey, cfg.GeminiModel)
	}

	model := models.New(cfg)

	runner := app.New(factory)
	if err := runner.Run(ctx, model); err != nil {
		fmt.Fprintf(os.Stderr, "❌  %v\n", err)
		os.Exit(1)
	}
}
