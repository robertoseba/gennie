package output

import (
	"bufio"
	"fmt"
	"strings"
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
