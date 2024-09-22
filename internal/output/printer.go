package output

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const codePrefix = "```"

func PrintAnswer(answer string) {
	scanner := bufio.NewScanner(strings.NewReader(answer))
	isCodeBlock := false

	fmt.Println()

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(strings.Trim(line, " "), codePrefix) {
			isCodeBlock = !isCodeBlock
			continue
		}

		if isCodeBlock {
			fmt.Printf("\t%s%s%s\n", Yellow, line, Reset)
			continue
		}

		fmt.Println(line)
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

func PrintLine() {
	lineChar := "\u2014"
	w, _ := GetTermSize()

	for i := 0; i < w; i++ {
		fmt.Print(lineChar)
	}
	fmt.Println()
}
