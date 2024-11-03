package output

import (
	"fmt"
	"strings"
)

type Question struct {
	question      *strings.Builder
	color         Color
	previousValue string
}

func NewQuestion(input string) *Question {
	question := strings.Builder{}
	question.WriteString(input)

	return &Question{
		question: &question,
		color:    Gray,
	}
}

func (q *Question) WithColor(color Color) *Question {
	q.color = color
	return q
}

func (q *Question) WithPrevious(previousValue string, IsMasked bool) *Question {
	q.previousValue = previousValue

	q.question.WriteString(" (Enter to use: ")

	if IsMasked {
		q.question.WriteString(previousValue[0:4])
		q.question.WriteString("...")
	} else {
		q.question.WriteString(previousValue)
	}
	q.question.WriteString(")")

	return q
}

func (q *Question) Ask(p *Printer) string {
	fmt.Fprintf(p.Stdout, "%s", q.color)
	fmt.Fprintf(p.Stdout, "%s\n", q.question.String())
	fmt.Fprintf(p.Stdout, "%s>%s ", Yellow, Reset)

	var input string
	fmt.Scanln(&input)

	if input == "" && q.previousValue != "" {
		return q.previousValue
	}

	return input
}
