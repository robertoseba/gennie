package ui

import "github.com/charmbracelet/lipgloss"

type helpView struct {
	text  string
	style lipgloss.Style
}

func newHelpView() *helpView {
	return &helpView{
		text: "Press q or ctrl+c to exit | j,k or arrows to scroll | f to follow up question",
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BBBBBB")).
			Height(1),
	}
}

func (h *helpView) View() string {
	if h.text == "" {
		return ""
	}
	return h.style.Render(h.text)
}
