package output

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestNewPrinter(t *testing.T) {
	p := NewPrinter(nil, nil)
	if p == nil {
		t.Errorf("Expected a Printer object, got nil")
	}
}

func TestPrintsWithCodeStyling(t *testing.T) {
	output := bytes.NewBufferString("")

	p := NewPrinter(output, nil)
	p.margin = ""

	p.PrintWithCodeStyling("This is a test\n```\ncode text\n```\n", Red)

	expected := fmt.Sprintf("This is a test\n%s%scode text\n%s", Reset, Red, Reset)

	printed := output.String()

	if printed != expected {
		t.Errorf("Expected '%#v', got '%#v'", expected, printed)
	}
}

func TestPrintLine(t *testing.T) {
	output := bytes.NewBufferString("")

	p := NewPrinter(output, nil)
	p.width = 10
	p.margin = ""

	p.PrintLine(Red)

	expected := fmt.Sprintf("%s%s%s\n", Red, strings.Repeat("\u2014", 10), Reset)

	printed := output.String()

	if printed != expected {
		t.Errorf("Expected '%#v', got '%#v'", expected, printed)
	}
}

func TestPrintWithMargin(t *testing.T) {
	output := bytes.NewBufferString("")

	margin := strings.Repeat(" ", 5)

	p := NewPrinter(output, nil)
	p.margin = margin

	p.Print("This is a test", "")

	expected := fmt.Sprintf("%sThis is a test%s\n%s", margin, margin, Reset)

	printed := output.String()

	if printed != expected {
		t.Errorf("Expected '%#v', got '%#v'", expected, printed)
	}
}

func TestPrinWithMarginAndColor(t *testing.T) {
	output := bytes.NewBufferString("")

	margin := strings.Repeat(" ", 5)

	p := NewPrinter(output, nil)
	p.margin = margin

	p.Print("This is a test", Red)

	expected := fmt.Sprintf("%s%sThis is a test%s\n%s", Red, margin, margin, Reset)

	printed := output.String()

	if printed != expected {
		t.Errorf("Expected '%#v', got '%#v'", expected, printed)
	}
}

func TestWrapWithMargins(t *testing.T) {
	text := "This is a test with a very long line that should be split in multiple lines"
	p := NewPrinter(nil, nil)
	p.width = 15
	p.marginSize = 2

	lines := p.wrapWithMargins(text, []string{})

	expected := []string{
		"This is a",
		"test with",
		"a very",
		"long line",
		"that",
		"should be",
		"split in",
		"multiple",
		"lines",
	}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Expected line '%s', got '%s'", expected[i], line)
		}

		if len(line) > 15 {
			t.Errorf("Expected line to be 15 characters long, got %d", len(line))
		}
	}

}

func TestGenMargin(t *testing.T) {
	p := NewPrinter(nil, nil)
	p.marginSize = 5

	if len(p.margin) != 5 {
		t.Errorf("Expected margin to be 5 characters long, got %d", len(p.margin))
	}

	if p.margin != "     " {
		t.Errorf("Expected margin to be '     ', got '%s'", p.margin)
	}
}
