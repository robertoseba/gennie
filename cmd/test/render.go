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
)

type contentMsg string

type status int

const (
	statusIdle status = iota
	statusAsking
	statusWaiting
	statusProcessing
	statusError
	statusAskConfirm
)

func (s status) String() string {
	switch s {
	case statusAsking:
		return "Asking the model"
	case statusWaiting:
		return "Waiting for answer"
	case statusProcessing:
		return "Processing response"
	case statusError:
		return "Error occurred"
	case statusAskConfirm:
		return "Should I continue?"
	default:
		return "Ready"
	}
}

type model struct {
	viewport     viewport.Model
	content      *strings.Builder
	contentChan  <-chan string
	spinner      spinner.Model
	status       status
	prompt       *huh.Confirm
	showPrompt   bool
	ctx          context.Context
	cancel       context.CancelFunc
	windowWidth  int
	windowHeight int
}

const maxWidth = 80

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("15")).
			Margin(0, 0, 1, 0).
			Width(maxWidth)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("60")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("60")).
			Margin(1, 0, 0, 0).
			Width(maxWidth - 2)

	errorStatusStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("196")).
				Foreground(lipgloss.Color("15")).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder())

	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func initialModel(contentChan <-chan string) model {
	ctx, cancel := context.WithCancel(context.Background())

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	vp := viewport.New(maxWidth, 20)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("15")).
		Width(maxWidth+2).
		Padding(0, 1).
		Margin(0)
	vp.SetContent("Waiting for content...\n")

	prompt := huh.NewConfirm()
	prompt.Negative("Cancel").Affirmative("Continue")

	return model{
		viewport:     vp,
		content:      &strings.Builder{},
		contentChan:  contentChan,
		spinner:      s,
		status:       statusIdle,
		ctx:          ctx,
		cancel:       cancel,
		prompt:       prompt,
		windowWidth:  maxWidth,
		windowHeight: 20,
	}
}

func waitForContent(ctx context.Context, contentChan <-chan string) tea.Cmd {
	return func() tea.Msg {
		select {
		case content, ok := <-contentChan:
			if !ok {
				return nil // Channel closed
			}
			return contentMsg(content)
		case <-ctx.Done():
			return nil
		}
	}
}

// Simulate asking the model
func askModel() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return statusMsg{statusWaiting}
	})
}

// Simulate waiting for response
func waitForResponse() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return statusMsg{statusProcessing}
	})
}

// Simulate processing
func processResponse() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return statusMsg{statusIdle}
	})
}

func confirmAction() tea.Cmd {
	return func() tea.Msg {
		return statusMsg{statusIdle}
	}
}

type statusMsg struct {
	status status
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForContent(m.ctx, m.contentChan),
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = min(msg.Width, maxWidth)
		m.windowHeight = msg.Height

		// Update viewport size (leaving space for status bar)
		headerHeight := 4
		statusBarStyle = statusBarStyle.UnsetWidth().Width(m.windowWidth - 2)
		HeaderStyle = HeaderStyle.UnsetWidth().Width(m.windowWidth - 2)
		m.viewport.Width = m.windowWidth
		m.viewport.Height = msg.Height - headerHeight - statusBarStyle.GetHeight() - 4

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.cancel()
			return m, tea.Quit

		case "j", "down":
			m.viewport.ScrollDown(2)
		case "k", "up":
			m.viewport.ScrollUp(2)
		case "c":
			m.content.Reset()
			m.viewport.SetContent("")
		case "r":
			// Simulate requesting something from model
			m.status = statusAsking
			cmds = append(cmds, askModel())
		case "e":
			// Simulate error
			m.status = statusError
		case "s":
			// Reset to idle
			m.status = statusIdle
			cmds = append(cmds, func() tea.Msg {
				return contentMsg("Resetting content...\n")
			})

		case "p":
			m.status = statusAskConfirm
			m.prompt.Title("Can I run this tool?")
			return m, nil
		}

	case contentMsg:
		// Append new content
		m.content.WriteString(string(msg))
		m.viewport.SetContent(m.content.String())
		m.viewport.GotoBottom()

		// Continue listening for more content
		if m.status == statusIdle {
			cmds = append(cmds, waitForContent(m.ctx, m.contentChan))
		}

	case statusMsg:
		m.status = msg.status

		// Chain status transitions
		switch msg.status {
		case statusWaiting:
			cmds = append(cmds, waitForResponse())
		case statusProcessing:
			cmds = append(cmds, processResponse())
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Always update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// Header
	headerText := "Press j/k to scroll, r to request, e for error, s to reset, q to quit"
	header := HeaderStyle.Render(headerText)

	// Main viewport
	var viewportView string
	if m.status != statusAskConfirm {
		viewportView = m.viewport.View()
	} else {
		m.prompt.Run()

		m.status = statusIdle
		viewportView = m.viewport.View()
	}

	// Status bar
	statusBar := m.renderStatusBar()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		viewportView,
		statusBar,
	)
}

func (m model) renderStatusBar() string {
	var statusText string
	var style lipgloss.Style

	switch m.status {
	case statusError:
		style = errorStatusStyle
		statusText = m.status.String()
	case statusIdle:
		style = statusBarStyle
		statusText = fmt.Sprintf("%s | Lines: %d | Scroll: %d/%d",
			m.status.String(),
			strings.Count(m.content.String(), "\n"),
			m.viewport.YOffset,
			max(0, m.viewport.TotalLineCount()-m.viewport.Height))
	case statusAskConfirm:
		style = statusBarStyle
		statusText = "Waiting for confirmation..."
	default:
		style = statusBarStyle
		statusText = fmt.Sprintf("%s %s", m.spinner.View(), m.status.String())
	}

	// Pad the status bar to full width
	paddedStatus := statusText + strings.Repeat(" ", max(0, m.windowWidth-lipgloss.Width(statusText)))

	return style.Render(paddedStatus)
}

func main() {
	// Create content channel
	contentChan := make(chan string, 10)

	// Start content producer
	go func() {
		defer close(contentChan)
		for i := 0; i < 100; i++ {
			contentChan <- fmt.Sprintf("Message %d at %s\n", i+1, time.Now().Format("15:04:05"))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	p := tea.NewProgram(
		initialModel(contentChan),
		tea.WithMouseCellMotion(),
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
