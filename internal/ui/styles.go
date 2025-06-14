package ui

import "github.com/charmbracelet/lipgloss"

var (
	borderColor  = lipgloss.AdaptiveColor{Light: "#e6f1f8", Dark: "#2f404a"}
	spinnerColor = lipgloss.AdaptiveColor{Light: "#a0efe1", Dark: "#a0efe1"}
	textColor    = lipgloss.AdaptiveColor{Light: "#F8EFBA", Dark: "#F8EFBA"}

	textStyle = lipgloss.NewStyle().Foreground(textColor)
)
