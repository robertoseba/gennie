package ui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type markdownModel struct {
	viewport.Model
	content  *strings.Builder
	renderer *glamour.TermRenderer
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

	var err error
	mv.renderer, err = glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(80), glamour.WithPreservedNewLines())
	if err != nil {
		log.Fatal("Failed to initialize markdown renderer")
	}

	return mv
}

func (m *markdownModel) setSize(w, h int) {
	m.Width = w + 2
	m.Height = h

	//Replaces glamour renderer with new size
	var err error
	m.renderer, err = glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(w-2), glamour.WithPreservedNewLines())

	//TODO: review log fatal
	if err != nil {
		log.Fatal("Failed to initialize markdown renderer")
	}
}

// Appends the content to the viewport and renders the markdown
// If we have a problem rendering the markdown we fallback to
// the basic string
func (m *markdownModel) appendContent(newContent string) {
	m.content.WriteString(newContent)
	mdRender, err := m.renderer.Render(m.content.String())
	if err != nil {
		mdRender = m.content.String()
	}
	m.SetContent(mdRender)
	m.GotoBottom()
}
