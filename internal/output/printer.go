package output

import (
	"bufio"
	"fmt"
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
}

func NewPrinter() *Printer {
	width, height := GetTermSize()
	var margin int = 5

	if width < 100 {
		margin = 1
	}

	return &Printer{
		width:      width,
		height:     height,
		marginSize: margin,
		margin:     strings.Repeat(" ", margin),
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

func (p *Printer) PrintAnswer(answer string) {
	scanner := bufio.NewScanner(strings.NewReader(answer))
	isCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(strings.Trim(line, " "), codePrefix) {
			isCodeBlock = !isCodeBlock
			continue
		}

		if isCodeBlock {
			p.Print(line, Yellow)
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

	for i := 0; i < p.width; i++ {
		fmt.Printf("%s%s%s", color, lineChar, Reset)
	}
	fmt.Print("\n")
}

func (p *Printer) PrintDetails(details string) {
	p.Print(details, Cyan)
}

func (p *Printer) Print(message string, color Color) {
	fullMessage := fmt.Sprintf("%s%s%s", color, message, Reset)

	lines := p.splitLine(fullMessage, []string{})

	fmt.Printf("%s", color)
	for _, text := range lines {
		fmt.Printf("%s%s%s\n", p.margin, text, p.margin)
	}
	fmt.Printf("%s", Reset)
}

func (p *Printer) splitLine(text string, initial []string) []string {
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
		return p.splitLine(text[idx+1:], initial)
	}

	initial = append(initial, cut)
	return p.splitLine(text[p.width-p.marginSize*2:], initial)
}
