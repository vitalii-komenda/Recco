package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"

	"steam-bubbles/internal/domain"
)

const maxSelectableGames = 200

func RunGameSelector(library []domain.Game) ([]domain.Game, error) {
	cap := len(library)
	if cap > maxSelectableGames {
		cap = maxSelectableGames
	}

	options := buildOptions(library[:cap])

	var selectedNames []string
	var customGames string

	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select games you enjoy").
				Description("Space to toggle · Enter to confirm · ↑↓ to move · q to quit").
				Options(options...).
				Height(18).
				Value(&selectedNames),
			huh.NewInput().
				Title("Any other favorites?").
				Description("Type games missing from the list (comma separated)").
				Value(&customGames),
		),
	).WithKeyMap(km)

	if err := form.Run(); err != nil {
		return nil, err
	}

	selected := resolveSelected(library, selectedNames)

	if strings.TrimSpace(customGames) != "" {
		for _, customName := range strings.Split(customGames, ",") {
			cleanedName := strings.TrimSpace(customName)
			if cleanedName != "" {
				selected = append(selected, domain.Game{Name: cleanedName})
			}
		}
	}

	return selected, nil
}

// buildOptions converts domain.Game slice into huh options with playtime labels.
func buildOptions(games []domain.Game) []huh.Option[string] {
	opts := make([]huh.Option[string], len(games))
	for i, g := range games {
		label := g.Name
		if g.PlaytimeMinutes > 0 {
			label += "  " + PlaytimeStyle.Render("("+g.FormattedPlaytime()+")")
		}
		opts[i] = huh.NewOption(label, g.Name)
	}
	return opts
}

// resolveSelected maps selected names back to their full domain.Game structs,
// preserving the order in which the user toggled them.
func resolveSelected(library []domain.Game, names []string) []domain.Game {
	index := make(map[string]domain.Game, len(library))
	for _, g := range library {
		index[g.Name] = g
	}

	selected := make([]domain.Game, 0, len(names))
	for _, n := range names {
		if g, ok := index[n]; ok {
			selected = append(selected, g)
		}
	}
	return selected
}

// PrintLibrarySummary prints a styled count of the loaded library.
func PrintLibrarySummary(count int) {
	fmt.Printf("\n%s\n\n",
		SectionStyle.Render(fmt.Sprintf("Found %d games in your library", count)),
	)
}
