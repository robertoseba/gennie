package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertoseba/gennie/internal/core/usecases/complete"
)

type statusModel struct {
	spinner.Model
	message    string
	style      lipgloss.Style
	errorStyle lipgloss.Style
	isActive   bool
	profile    string
	model      string
}

type promptModel struct {
	*huh.Confirm
	isActive bool
}

type helpView struct {
	text  string
	style lipgloss.Style
}

type ui struct {
	prompt       *promptModel
	help         *helpView
	status       *statusModel
	markdownView *markdownModel
}

func (u *ui) setSize(w, h int) {
	u.help.style = u.help.style.Width(w)
	u.status.style = u.status.style.Width(w)
	u.status.errorStyle = u.status.errorStyle.Width(w)
	u.markdownView.setSize(w, h-u.help.style.GetHeight()-u.status.style.GetHeight()-4)

}

type model struct {
	ui              *ui
	width           int
	height          int
	ctx             context.Context
	cancel          context.CancelFunc
	responseChan    <-chan complete.Response
	isDoneAnswering bool
	startedAt       time.Time
	finishedAt      time.Time
}

func (m *model) setSize(w, h int) {
	m.width = min(w, 120)
	m.height = h - 2
	m.ui.setSize(m.width, m.height)
}

func newStatusView(height int) *statusModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("15")).
		Padding(0, 2).
		Height(height)

	return &statusModel{
		Model:      spin,
		message:    "Starting...",
		style:      statusStyle,
		errorStyle: statusStyle.Foreground(lipgloss.Color("205")).Bold(true),
		isActive:   true,
	}
}

func NewPromptModel() *promptModel {
	prompt := huh.NewConfirm()
	prompt.Negative("Cancel").Affirmative("Continue")
	prompt.WithTheme(huh.ThemeDracula())
	return &promptModel{
		Confirm:  prompt,
		isActive: false,
	}
}

func newHelpView(height int) *helpView {
	return &helpView{
		text: "Press q or ctrl+c to exit | j,k or arrows to scroll | f to follow up question",
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BBBBBB")).
			Height(height),
	}
}

func newUi() *ui {
	ui := &ui{
		prompt:       NewPromptModel(),
		help:         newHelpView(1),
		status:       newStatusView(1),
		markdownView: newMarkdownModel(),
	}

	return ui
}

func NewModel(responseChan <-chan complete.Response) model {
	ctx, cancel := context.WithCancel(context.Background())

	return model{
		ctx:          ctx,
		cancel:       cancel,
		width:        80,
		height:       20,
		ui:           newUi(),
		responseChan: responseChan,
		startedAt:    time.Now(),
	}
}

type DoneAnswer string

func waitForContent(ctx context.Context, responseChan <-chan complete.Response) tea.Cmd {
	return func() tea.Msg {
		select {
		case content, ok := <-responseChan:
			if !ok {
				return DoneAnswer("done")
			}
			return content

		case <-ctx.Done():
			return DoneAnswer("done")
		}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForContent(m.ctx, m.responseChan),
		m.ui.status.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case DoneAnswer:
		m.isDoneAnswering = true
		m.finishedAt = time.Now()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.cancel()
			return m, tea.Quit

			// case "p":
			// 	m.status = statusAskConfirm
			// 	m.prompt = m.prompt.Title("Can I run this tool?")
			// 	cmd = m.prompt.Focus()
			// 	cmds = append(cmds, cmd)
		}

	case complete.Response:
		// Append new content
		switch msg.Type {
		case complete.RtLlmAnswer:
			m.ui.markdownView.appendContent(msg.Data)
		case complete.RtLoading:
			m.ui.status.message = msg.Data
		case complete.RtModel:
			m.ui.status.model = msg.Data
		case complete.RtProfile:
			m.ui.status.profile = msg.Data
		}
		cmds = append(cmds, waitForContent(m.ctx, m.responseChan))

	case spinner.TickMsg:
		if !m.isDoneAnswering {
			m.ui.status.Model, cmd = m.ui.status.Model.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Always update viewport
	// if m.status == statusAskConfirm {
	// 	prompt, cmd := m.prompt.Update(msg)
	// 	m.prompt = prompt.(*huh.Confirm)
	// 	if cmd != nil {
	// 		cmds = append(cmds, cmd)
	// 	}
	// 	return m, tea.Batch(cmds...)
	// }

	m.ui.markdownView.Model, cmd = m.ui.markdownView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var help string
	statusBar := m.ui.status.style.Render(fmt.Sprintf("%s %s | Model: %s | Profile: %s", m.ui.status.View(), m.ui.status.message, m.ui.status.model, m.ui.status.profile))
	help = m.ui.help.style.Render(m.ui.help.text)

	if m.isDoneAnswering {
		statusBar = m.ui.status.style.Render(fmt.Sprintf("Model: %s | Profile: %s", m.ui.status.model, m.ui.status.profile))
	}

	// Main viewport
	// var viewportView string
	// // if m.status != statusAskConfirm {
	// 	viewportView = m.markdownView.View()
	// } else {
	// 	m.status = statusIdle
	// 	m.prompt = m.prompt.Title("Confirm Action")
	// 	viewportView = lipgloss.Place(m.markdownView.Width, m.markdownView.Height, lipgloss.Center, lipgloss.Center, m.prompt.View())
	// }

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		m.ui.markdownView.View(),
		help,
	)
}

func Run(contentChan <-chan complete.Response) error {
	p := tea.NewProgram(
		NewModel(contentChan),
		// tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start ui: %w", err)
	}
	return nil
}
