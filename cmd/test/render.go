package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertoseba/gennie/internal/core/usecases/complete"
)

type markdownModel struct {
	viewport.Model
	content *strings.Builder
	style   lipgloss.Style
}

type statusModel struct {
	spinner.Model
	message    string
	style      lipgloss.Style
	errorStyle lipgloss.Style
	isActive   bool
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
	u.markdownView.Width = w
	u.markdownView.Height = h - u.help.style.GetHeight() - u.status.style.GetHeight() - 4
	u.status.style.Width(w)
	u.status.errorStyle.Width(w)
}

type model struct {
	ui           *ui
	width        int
	height       int
	ctx          context.Context
	cancel       context.CancelFunc
	responseChan <-chan complete.Response
}

func (m *model) setSize(w, h int) {
	m.width = min(w, 120)
	m.height = h - 2
	m.ui.setSize(m.width, m.height)
}

var HeaderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("60")).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("60")).
	Margin(1, 0, 0, 0).
	Padding(0, 1)

func newStatusView(height int) *statusModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("15")).
		Margin(0, 0, 1, 0).
		Padding(0, 1).
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

func newMarkdownModel() *markdownModel {
	return &markdownModel{
		Model:   viewport.New(80, 20),
		content: &strings.Builder{},
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("15")).
			Padding(0, 1).
			Margin(0),
	}
}

func newHelpView(height int) *helpView {
	return &helpView{
		text: "Press q or ctrl+c to exit | j,k or arrows to scroll | f to follow up question",
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Foreground(lipgloss.Color("9")).
			Height(height),
	}
}

func newUi() *ui {
	ui := &ui{
		prompt:       NewPromptModel(),
		help:         newHelpView(4),
		status:       newStatusView(4),
		markdownView: newMarkdownModel(),
	}

	return ui
}

func newModel(responseChan <-chan complete.Response) model {
	ctx, cancel := context.WithCancel(context.Background())

	return model{
		ctx:          ctx,
		cancel:       cancel,
		width:        80,
		height:       20,
		ui:           newUi(),
		responseChan: responseChan,
	}
}

func waitForContent(ctx context.Context, responseChan <-chan complete.Response) tea.Cmd {
	return func() tea.Msg {
		select {
		case content, ok := <-responseChan:
			if !ok {
				return nil
			}
			return content

		case <-ctx.Done():
			return nil
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

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.cancel()
			return m, tea.Quit

			// case "j", "down":
			// 	m.markdownView.ScrollDown(2)
			// case "k", "up":
			// 	m.markdownView.ScrollUp(2)

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
			m.ui.markdownView.content.WriteString(msg.Data)
			m.ui.markdownView.SetContent(m.ui.markdownView.content.String())
			m.ui.markdownView.GotoBottom()
		}

	case spinner.TickMsg:
		m.ui.status.Model, cmd = m.ui.status.Model.Update(msg)
		cmds = append(cmds, cmd)
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

	m.ui.markdownView.Model, cmd = m.ui.markdownView.Model.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Header
	help := m.ui.help.style.Render(m.ui.help.text)

	// Main viewport
	// var viewportView string
	// // if m.status != statusAskConfirm {
	// 	viewportView = m.markdownView.View()
	// } else {
	// 	m.status = statusIdle
	// 	m.prompt = m.prompt.Title("Confirm Action")
	// 	viewportView = lipgloss.Place(m.markdownView.Width, m.markdownView.Height, lipgloss.Center, lipgloss.Center, m.prompt.View())
	// }

	// Status bar
	// statusBar := m.renderStatusBar()
	//
	return lipgloss.JoinVertical(
		lipgloss.Left,
		// statusBar,
		help,
		m.ui.markdownView.View(),
	)
}

func (m model) renderStatusBar() string {
	return ""
	// var statusText string
	// var style lipgloss.Style
	//
	// switch m.status {
	// case statusError:
	// 	style = errorStatusStyle
	// 	statusText = m.status.String()
	// case statusIdle:
	// 	style = statusBarStyle
	// 	statusText = fmt.Sprintf("%s | Lines: %d | Scroll: %d/%d",
	// 		m.status.String(),
	// 		strings.Count(m.content.String(), "\n"),
	// 		m.markdownView.YOffset,
	// 		max(0, m.markdownView.TotalLineCount()-m.markdownView.Height))
	// case statusAskConfirm:
	// 	style = statusBarStyle
	// 	statusText = "Waiting for confirmation..."
	// default:
	// 	style = statusBarStyle
	// 	statusText = fmt.Sprintf("%s %s", m.spinner.View(), m.status.String())
	// }
	//
	// // Pad the status bar to full width
	// paddedStatus := statusText + strings.Repeat(" ", max(0, m.width-lipgloss.Width(statusText)))
	//
	// return style.Render(paddedStatus)
}

func main() {
	// Create content channel
	contentChan := make(chan complete.Response, 10)

	// Start content producer
	go func() {
		defer close(contentChan)
		for i := 0; i < 100; i++ {
			contentChan <- complete.Response{
				Data: fmt.Sprintf("%d - message", i),
				Err:  nil,
				Type: complete.RtLlmAnswer,
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	p := tea.NewProgram(
		newModel(contentChan),
		// tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
