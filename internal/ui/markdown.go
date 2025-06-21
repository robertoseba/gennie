package ui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type markdownModel struct {
	viewport        viewport.Model
	content         *strings.Builder
	renderedContent string
	renderer        *glamour.TermRenderer
}

func newMarkdownModel() *markdownModel {
	mv := &markdownModel{
		viewport: viewport.New(80, 40),
		content:  &strings.Builder{},
	}
	mv.viewport.Style = lipgloss.NewStyle().
		// Border(lipgloss.NormalBorder()).
		// BorderForeground(borderColor).
		Padding(0, 1)

	var err error
	mv.renderer, err = glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithPreservedNewLines(), glamour.WithWordWrap(maxWidth))
	if err != nil {
		log.Fatal("Failed to initialize markdown renderer")
	}

	return mv
}

func (m *markdownModel) setSize(w, h int) {
	m.viewport.Width = w + 2
	m.viewport.Height = h
}

func (m *markdownModel) height() int {
	return lipgloss.Height(m.renderedContent)
}

// Appends the content to the viewport and renders the markdown
// If we have a problem rendering the markdown we fallback to
// the basic string
func (m *markdownModel) appendContent(newContent string) {
	m.content.WriteString(newContent)
	var err error
	m.renderedContent, err = m.renderer.Render(m.content.String())
	if err != nil {
		m.renderedContent = m.content.String()
	}

	m.viewport.SetContent(m.renderedContent)
	m.viewport.GotoBottom()
}

func (m *markdownModel) View() string {
	if m.height() > m.viewport.Height {
		return m.viewport.View()
	}
	return m.renderedContent
}
