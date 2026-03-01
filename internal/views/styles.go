package views

import "github.com/charmbracelet/lipgloss"

var (
	colorPurple = lipgloss.Color("#7C3AED")
	colorTeal   = lipgloss.Color("#14B8A6")
	colorYellow = lipgloss.Color("#FBBF24")
	colorSubtle = lipgloss.Color("#6B7280")
	colorWhite  = lipgloss.Color("#F9FAFB")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPurple).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorTeal).
			Padding(0, 2).
			MarginBottom(1)

	SectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorTeal).
			MarginTop(1).
			MarginBottom(1)

	GameItemStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			PaddingLeft(2)

	PlaytimeStyle = lipgloss.NewStyle().
			Foreground(colorSubtle).
			Italic(true)

	RecsBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPurple).
			Padding(1, 2).
			MarginTop(1)

	RecsTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow).
			MarginBottom(1)

	BulletStyle = lipgloss.NewStyle().
			Foreground(colorTeal)

	BoldInlineStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow)
)
