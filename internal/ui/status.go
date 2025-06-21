package ui

import (
	"time"

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
	startedAt   time.Time
	finishedAt  time.Time
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
		startedAt:   time.Now(),
		finishedAt:  time.Now(),
	}
}

func (m *statusModel) finishNow() {
	m.finishedAt = time.Now()
	m.isActive = false
	m.message = ""
}

func (m *statusModel) setSize(w, h int) {
	m.borderStyle = m.borderStyle.Width(w)
	m.errorStyle = m.errorStyle.Width(w)
}

func (m *statusModel) View() string {
	if !m.isActive {
		elapsed := m.finishedAt.Sub(m.startedAt)
		return m.borderStyle.Render(m.textStyle.Render("Finished in  -> ", elapsed.String(), " | ", m.model, m.profile))
	}
	return m.borderStyle.Render(m.Model.View() + m.textStyle.Render(m.message, " -> ", m.model, m.profile))
}
