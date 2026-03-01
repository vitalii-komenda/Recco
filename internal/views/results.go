package views

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"steam-bubbles/internal/domain"
)

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

	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Reason", Width: 90},
		{Title: "Platform", Width: 20},
	}

	var rows []table.Row
	for _, rec := range recommendations {
		rows = append(rows, table.Row{rec.Name, rec.Reason, rec.Platform})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),

		table.WithHeight(len(recommendations)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorTeal).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(colorWhite).
		Background(colorPurple).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			row := m.table.SelectedRow()
			if row != nil {
				name := row[0]
				searchURL := "https://store.steampowered.com/search/?term=" + url.QueryEscape(name)
				openBrowser(searchURL)
			}
		case "up", "k":
			if m.table.Cursor() == 0 && len(m.table.Rows()) > 0 {
				m.table.SetCursor(len(m.table.Rows()) - 1)
				return m, nil
			}
		case "down", "j":
			if m.table.Cursor() == len(m.table.Rows())-1 && len(m.table.Rows()) > 0 {
				m.table.SetCursor(0)
				return m, nil
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	helpText := PlaytimeStyle.Render("  enter to open in Steam • ↑/↓ to scroll • q to quit")
	return RecsBoxStyle.Render(m.table.View()) + "\n" + helpText + "\n"
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		// Silently fallback if it fails, or log it if there's a good place for it
		_ = err
	}
}
