package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type statusModel struct {
	spinner.Model
	message     string
	borderStyle lipgloss.Style
	textStyle   lipgloss.Style
	errorStyle  lipgloss.Style
	isActive    bool
	profile     string
	model       string
}

func newStatusView() *statusModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(spinnerColor)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Height(1).
		Padding(0, 2)

	return &statusModel{
		Model:       spin,
		message:     "Starting...",
		borderStyle: borderStyle,
		textStyle:   textStyle,
		errorStyle:  textStyle.Foreground(lipgloss.Color("205")).Bold(true),
		isActive:    true,
	}
}

func (m *statusModel) setSize(w, h int) {
	m.borderStyle = m.borderStyle.Width(w)
	m.errorStyle = m.errorStyle.Width(w)
}

func (m *statusModel) View() string {
	return m.borderStyle.Render(m.Model.View() + m.textStyle.Render(m.message, " -> ", m.model, m.profile))
}
