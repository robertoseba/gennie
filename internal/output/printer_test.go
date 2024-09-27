package output

import (
	"testing"
)

func TestNewPrinter(t *testing.T) {
	p := NewPrinter()
	if p == nil {
		t.Errorf("Expected a Printer object, got nil")
	}
}

func TestSplitLine(t *testing.T) {
	text := "This is a test with a very long line that should be split in multiple lines"
	p := NewPrinter()
	p.width = 15
	p.marginSize = 2

	lines := p.wrapWithMargins(text, []string{})
	for _, line := range lines {
		if len(line) > 15 {
			t.Errorf("Expected line to be 15 characters long, got %d", len(line))
		}
	}
}

func TestGenMargin(t *testing.T) {
	p := NewPrinter()
	p.marginSize = 5

	if len(p.margin) != 5 {
		t.Errorf("Expected margin to be 5 characters long, got %d", len(p.margin))
	}

	if p.margin != "     " {
		t.Errorf("Expected margin to be '     ', got '%s'", p.margin)
	}
}
