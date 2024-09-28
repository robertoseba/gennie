// Inspired/Forked from https://github.com/Nexidian/gocliselect
// I decided to create my own version without dependencies for the menu

package output

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var up byte = 65
var down byte = 66
var escape byte = 27
var enter byte = 13
var ctrlC byte = 3

var arrowOn = fmt.Sprintf(" %s>%s", Cyan, Reset)
var arrowOff = fmt.Sprintf(" %s %s", Cyan, Reset)

var keys = map[byte]bool{
	up:   true,
	down: true,
}

type Menu struct {
	Prompt    string
	CursorPos int
	MenuItems []*MenuItem
}

type MenuItem struct {
	Text    string
	ID      string
	SubMenu *Menu
}

func NewMenu(prompt string) *Menu {
	return &Menu{
		Prompt:    prompt,
		MenuItems: make([]*MenuItem, 0),
	}
}

func (m *Menu) AddItem(option string, id string) *Menu {
	menuItem := &MenuItem{
		Text: option,
		ID:   id,
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return m
}

func (m *Menu) renderMenuItems(redraw bool) {
	if redraw {
		// Move the cursor up n lines where n is the number of options, setting the new
		// location to start printing from, effectively redrawing the option list
		//
		// This is done by sending a VT100 escape code to the terminal
		// @see http://www.climagic.org/mirrors/VT100_Escape_Codes.html
		fmt.Printf("\033[%dA", len(m.MenuItems)-1)
	}

	for index, menuItem := range m.MenuItems {
		var newline = "\n"
		if index == len(m.MenuItems)-1 {
			// Adding a new line on the last option will move the cursor position out of range
			// For out redrawing
			newline = ""
		}

		menuItemText := menuItem.Text
		if index == m.CursorPos {
			fmt.Printf("\r%s %s%s%s%s", arrowOn, Cyan, menuItemText, Reset, newline)
			continue
		}

		fmt.Printf("\r%s %s%s", arrowOff, menuItemText, newline)
	}
}

func (m *Menu) clearMenu() {
	fmt.Printf("\033[%dA", len(m.MenuItems))
	for i := 0; i < len(m.MenuItems)+1; i++ {
		fmt.Print("\033[G")
		fmt.Print("\033[K")
		fmt.Print("\033[G")
		fmt.Print("\n")
	}

	fmt.Printf("\033[%dA", len(m.MenuItems)+1)
}

func (m *Menu) Display(selectedIdx int) string {
	defer showCursor()

	fmt.Printf(" %s%s%s\n", Yellow, m.Prompt, Reset)
	m.CursorPos = selectedIdx
	m.renderMenuItems(false)

	hideCursor()

	for {
		keyCode := getInput()
		switch keyCode {
		case escape, ctrlC:
			m.clearMenu()
			return ""
		case enter:
			menuItem := m.MenuItems[m.CursorPos]
			m.clearMenu()
			return menuItem.ID
		case up:
			m.CursorPos = (m.CursorPos + len(m.MenuItems) - 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		case down:
			m.CursorPos = (m.CursorPos + 1) % len(m.MenuItems)
			m.renderMenuItems(true)
		}
	}
}

func getInput() byte {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var read int

	readBytes := make([]byte, 3)

	read, _ = os.Stdin.Read(readBytes)

	if read == 3 {
		if _, ok := keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}

	return 0
}

func showCursor() {
	fmt.Printf("\033[?25h")
}

func hideCursor() {
	fmt.Printf("\033[?25l")
}
