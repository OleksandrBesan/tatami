package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#06B6D4")
	secondaryColor = lipgloss.Color("#22D3EE")
	mutedColor     = lipgloss.Color("#6B7280")
	successColor   = lipgloss.Color("#10B981")
	errorColor     = lipgloss.Color("#EF4444")

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Selected item style
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Normal item style
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	// Muted text style
	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Success style
	successStyle = lipgloss.NewStyle().
			Foreground(successColor)

	// Box style for panels
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Input label style
	labelStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)
)
