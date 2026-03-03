package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"steam-bubbles/internal/app"
	"steam-bubbles/internal/config"
	"steam-bubbles/internal/domain"
	infraGemini "steam-bubbles/internal/infrastructure/gemini"
	infraSteam "steam-bubbles/internal/infrastructure/steam"
	"steam-bubbles/internal/models"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "print version information")
	flag.BoolVar(showVersion, "v", false, "print version information")
	flag.Parse()

	if *showVersion {
		fmt.Printf("recco version=%s commit=%s date=%s\n", version, commit, date)
		return
	}

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
