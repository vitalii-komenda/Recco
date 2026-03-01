// Package models defines the ELM-style application state (the "Model").
// The model is the single source of truth for everything the app knows at runtime.
package models

import (
	"steam-bubbles/internal/config"
	"steam-bubbles/internal/domain"
)

type Step int

const (
	StepCredentials Step = iota // waiting for user to supply API keys
	StepFetchGames              // fetching Steam library
	StepSelectGames             // user picking games from list
	StepFetchRecs               // asking Gemini for recommendations
	StepShowResults             // rendering final results
)

// AppModel is the central application state — the ELM "Model".
// All mutable state lives here; nothing else should hold shared state.
type AppModel struct {
	// Current step in the application flow.
	Step Step

	// Configuration (API keys, Steam ID).
	Config config.Config

	// Full library returned by Steam, sorted by playtime desc.
	Library []domain.Game

	// Games the user selected in the multi-select step.
	Selected []domain.Game

	// Recommendations returned by Gemini.
	Recommendations []domain.Recommendation

	// Any error that terminated the current step — checked by the app runner.
	Err error
}

func New(cfg config.Config) *AppModel {
	return &AppModel{
		Step:   StepCredentials,
		Config: cfg,
	}
}
