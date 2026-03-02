package views

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"steam-bubbles/internal/domain"
)

type recItem struct {
	rec domain.Recommendation
}

func (r recItem) Title() string {
	return lipgloss.NewStyle().Bold(true).Foreground(colorWhite).Render(r.rec.Name) +
		lipgloss.NewStyle().Foreground(colorTeal).Italic(true).Render("  ["+r.rec.Platform+"]")
}

func (r recItem) Description() string { return r.rec.Reason }
func (r recItem) FilterValue() string { return r.rec.Name }

type recsModel struct {
	list list.Model
}

func (m recsModel) Init() tea.Cmd { return nil }

func (m recsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(recItem); ok {
				searchURL := "https://store.steampowered.com/search/?term=" + url.QueryEscape(item.rec.Name)
				openBrowser(searchURL)
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-4, msg.Height-4)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m recsModel) View() string {
	return RecsBoxStyle.Render(m.list.View())
}

func RenderResults(selected []domain.Game, recommendations []domain.Recommendation) {
	fmt.Println()
	fmt.Println(SectionStyle.Render("🎮  Games you selected:"))

	for _, g := range selected {
		pt := PlaytimeStyle.Render("(" + g.FormattedPlaytime() + ")")
		fmt.Println(GameItemStyle.Render("• " + g.Name + "  " + pt))
	}

	fmt.Println()
	fmt.Println(RecsTitleStyle.Render("✨  Gemini Recommendations"))

	if len(recommendations) == 0 {
		fmt.Println("No recommendations found.")
		return
	}

	items := make([]list.Item, len(recommendations))
	for i, rec := range recommendations {
		items[i] = recItem{rec: rec}
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.SetSpacing(1)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(colorWhite).
		BorderForeground(colorPurple)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#D1D5DB")).
		BorderForeground(colorPurple)

	l := list.New(items, delegate, 100, 30)
	l.Title = "✨  Gemini Recommendations"
	l.Styles.Title = RecsTitleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	m := recsModel{list: l}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func openBrowser(rawURL string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", rawURL).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL).Start()
	case "darwin":
		err = exec.Command("open", rawURL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		_ = err
	}
}
