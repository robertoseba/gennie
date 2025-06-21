package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/robertoseba/gennie/internal/core/usecases/complete"
)

type (
	DoneAnswer  string
	promptModel struct {
		*huh.Confirm
		isActive bool
	}
)

var maxWidth = 120

type ui struct {
	prompt   *promptModel
	help     *helpView
	status   *statusModel
	markdown *markdownModel
}

func (u *ui) setSize(w, h int) {
	u.help.style = u.help.style.Width(w)
	u.status.setSize(w, h)
	u.markdown.setSize(w, h-u.help.style.GetHeight()-u.status.borderStyle.GetHeight()-4)
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
	m.width = min(w, maxWidth)
	m.height = h - 2
	m.ui.setSize(m.width, m.height)
}

func newPromptModel() *promptModel {
	prompt := huh.NewConfirm()
	prompt.Negative("Cancel").Affirmative("Continue")
	prompt.WithTheme(huh.ThemeDracula())
	return &promptModel{
		Confirm:  prompt,
		isActive: false,
	}
}

func newUi() *ui {
	ui := &ui{
		prompt:   newPromptModel(),
		help:     newHelpView(),
		status:   newStatusView(),
		markdown: newMarkdownModel(),
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
	}
}

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

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.cancel()
			return m, tea.Quit
		}

	case DoneAnswer:
		m.isDoneAnswering = true
		m.ui.status.finishNow()

	case spinner.TickMsg:
		if !m.isDoneAnswering {
			m.ui.status.Model, cmd = m.ui.status.Model.Update(msg)
			cmds = append(cmds, cmd)
		}

	case complete.Response:
		switch msg.Type {
		case complete.RtLlmAnswer:
			m.ui.markdown.appendContent(msg.Data)
		case complete.RtLoading:
			m.ui.status.message = msg.Data
		case complete.RtModel:
			m.ui.status.model = msg.Data
		case complete.RtProfile:
			m.ui.status.profile = msg.Data
		case complete.RtApprovalReq:
			m.ui.prompt.Confirm = m.ui.prompt.Title("Approval Request").Description(msg.Data)
			m.ui.prompt.isActive = true
			cmd = m.ui.prompt.Focus()
			cmds = append(cmds, cmd)
		case complete.RtError:
			m.ui.markdown.appendContent(fmt.Sprintf("---\n## Error\n %s", msg.Err))
			return m, tea.Quit
		}
		cmds = append(cmds, waitForContent(m.ctx, m.responseChan))
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

	if m.isDoneAnswering && !m.isViewportActive() {
		return m, tea.Quit
	}

	m.ui.markdown.viewport, cmd = m.ui.markdown.viewport.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var help, content string
	if m.isViewportActive() {
		help = m.ui.help.View()
		content = m.ui.markdown.View()
	} else {
		content = m.ui.markdown.ViewNoViewport()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.ui.status.View(),
		content,
		help,
	)
}

func (m *model) isViewportActive() bool {
	return m.ui.markdown.height() > m.height-10
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
