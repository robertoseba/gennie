package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

const codePrefix = "```"

type Printer struct {
	width      int
	height     int
	marginSize int
	margin     string
	Stdout     io.Writer
	Stderr     io.Writer
}

func NewPrinter(out io.Writer, err io.Writer) *Printer {
	width, height := GetTermSize()
	var margin int = 5

	if width < 100 {
		margin = 1
	}

	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}

	return &Printer{
		width:      width,
		height:     height,
		marginSize: margin,
		margin:     strings.Repeat(" ", margin),
		Stdout:     out,
		Stderr:     err,
	}
}

func GetTermSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 100
		height = 80
	}
	return width, height
}

func (p *Printer) PrintWithCodeStyling(answer string, codeColor Color) {
	scanner := bufio.NewScanner(strings.NewReader(answer))
	isCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(strings.Trim(line, " "), codePrefix) {
			isCodeBlock = !isCodeBlock
			continue
		}

		if isCodeBlock {
			p.Print(line, codeColor)
			continue
		}

		p.Print(line, "")
	}
}

func (p *Printer) PrintLine(color Color) {
	if color == "" {
		color = Gray
	}

	lineChar := "\u2014"
	line := strings.Repeat(lineChar, p.width)
	fmt.Fprintf(p.Stdout, "%s%s%s\n", color, line, Reset)
}

func (p *Printer) Print(message string, color Color) {
	lines := p.wrapWithMargins(message, []string{})

	fmt.Fprintf(p.Stdout, "%s", color)

	for _, text := range lines {
		fmt.Fprintf(p.Stdout, "%s%s%s\n", p.margin, text, p.margin)
	}

	fmt.Fprintf(p.Stdout, "%s", Reset)
}

func (p *Printer) wrapWithMargins(text string, initial []string) []string {
	if len(text)+p.marginSize*2 <= p.width {
		return append(initial, text)
	}

	cut := text[:p.width-p.marginSize*2]

	if cut[len(cut)-1] != ' ' {
		idx := strings.LastIndex(cut, " ")
		if idx == -1 {
			idx = len(cut)
		}
		initial = append(initial, cut[:idx])
		return p.wrapWithMargins(text[idx+1:], initial)
	}

	initial = append(initial, cut)
	return p.wrapWithMargins(text[p.width-p.marginSize*2:], initial)
}
