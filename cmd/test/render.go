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
}

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
	u.markdownView.Width = w + 2
	u.markdownView.Height = h - u.help.style.GetHeight() - u.status.style.GetHeight() - 4
	u.status.style = u.status.style.Width(w)
	u.status.errorStyle = u.status.errorStyle.Width(w)
}

type model struct {
	ui              *ui
	width           int
	height          int
	ctx             context.Context
	cancel          context.CancelFunc
	responseChan    <-chan complete.Response
	isDoneAnswering bool
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

func newMarkdownModel() *markdownModel {
	mv := &markdownModel{
		Model:   viewport.New(80, 20),
		content: &strings.Builder{},
	}
	mv.Style = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("15")).
		Padding(0, 1)

	return mv
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
			m.ui.markdownView.content.WriteString(msg.Data)
			m.ui.markdownView.SetContent(m.ui.markdownView.content.String())
			m.ui.markdownView.GotoBottom()
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

	m.ui.markdownView.Model, cmd = m.ui.markdownView.Model.Update(msg)
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

func main() {
	contentChan := make(chan complete.Response, 10)

	go func() {
		defer close(contentChan)
		for i := 0; i < 100; i++ {
			if i == 40 {
				contentChan <- complete.Response{
					Data: fmt.Sprintf("Asking the model"),
					Err:  nil,
					Type: complete.RtLoading,
				}
			}
			if i == 1 {
				contentChan <- complete.Response{
					Data: fmt.Sprintf("gpt-4.1-mini"),
					Err:  nil,
					Type: complete.RtModel,
				}
			}
			if i == 2 {
				contentChan <- complete.Response{
					Data: fmt.Sprintf("personal"),
					Err:  nil,
					Type: complete.RtProfile,
				}
			}
			contentChan <- complete.Response{
				Data: fmt.Sprintf("%d - message\n", i),
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
