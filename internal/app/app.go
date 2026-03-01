// Package app contains the application orchestrator — the ELM-style Update loop.
// It drives state transitions by reading the current model step, executing side
// effects (I/O, API calls), and advancing the model to the next step.
package app

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"

	"steam-bubbles/internal/config"
	"steam-bubbles/internal/domain"
	"steam-bubbles/internal/models"
	"steam-bubbles/internal/views"
)

type ClientFactory func(cfg config.Config) (repo domain.GameRepository, rec domain.Recommender)
type Runner struct {
	factory ClientFactory
}

func New(factory ClientFactory) *Runner {
	return &Runner{factory: factory}
}

func (r *Runner) Run(ctx context.Context, m *models.AppModel) error {
	var repo domain.GameRepository
	var rec domain.Recommender

	for {
		switch m.Step {
		case models.StepCredentials:
			if err := handleCredentials(m); err != nil {
				return err
			}
			repo, rec = r.factory(m.Config)

		case models.StepFetchGames:
			if err := handleFetchGames(ctx, m, repo); err != nil {
				return err
			}

		case models.StepSelectGames:
			if err := handleSelectGames(m); err != nil {
				return err
			}

		case models.StepFetchRecs:
			if err := handleFetchRecs(ctx, m, rec); err != nil {
				return err
			}

		case models.StepShowResults:
			handleShowResults(m)
			return nil
		}
	}
}

func handleCredentials(m *models.AppModel) error {
	if err := views.RunCredentialForm(&m.Config); err != nil {
		return fmt.Errorf("credential input aborted: %w", err)
	}
	if err := m.Config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	advance(m, models.StepFetchGames)
	return nil
}

func handleFetchGames(ctx context.Context, m *models.AppModel, repo domain.GameRepository) error {
	fmt.Println()
	fmt.Println(views.TitleStyle.Render(" 🎮  Recco "))

	var fetchErr error

	_ = spinner.New().
		Title(" Loading your Steam games…").
		Action(func() {
			m.Library, fetchErr = repo.GetOwnedGames(ctx, m.Config.SteamID)
		}).
		Run()

	if fetchErr != nil {
		return fmt.Errorf("failed to fetch Steam library: %w", fetchErr)
	}

	views.PrintLibrarySummary(len(m.Library))
	advance(m, models.StepSelectGames)
	return nil
}

func handleSelectGames(m *models.AppModel) error {
	selected, err := views.RunGameSelector(m.Library)
	if err != nil {
		return fmt.Errorf("game selection aborted: %w", err)
	}

	if len(selected) == 0 {
		fmt.Fprintln(os.Stdout, "No games selected. Bye! 👋")
		os.Exit(0)
	}

	m.Selected = selected
	advance(m, models.StepFetchRecs)
	return nil
}

func handleFetchRecs(ctx context.Context, m *models.AppModel, rec domain.Recommender) error {
	var recErr error

	_ = spinner.New().
		Title(" Asking Gemini for recommendations…").
		Action(func() {
			m.Recommendations, recErr = rec.Recommend(ctx, m.Selected)
		}).
		Run()

	if recErr != nil {
		return fmt.Errorf("failed to get recommendations: %w", recErr)
	}

	advance(m, models.StepShowResults)
	return nil
}

func handleShowResults(m *models.AppModel) {
	views.RenderResults(m.Selected, m.Recommendations)
	fmt.Println()
}

func advance(m *models.AppModel, next models.Step) {
	m.Step = next
}
