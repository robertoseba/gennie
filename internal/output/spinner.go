package output

import (
	"fmt"
	"time"
)

type Spinner struct {
	frames    []string
	index     int
	done      chan struct{}
	message   string
	isRunning bool
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:    []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
		index:     0,
		done:      make(chan struct{}),
		message:   message,
		isRunning: false,
	}
}

func (s *Spinner) SetMessage(message string) {
	s.message = message
	fmt.Print("\r\033[K")
	if s.index == 0 {
		fmt.Printf("\r%s%s %s%s", Cyan, s.frames[s.index], s.message, Gray)
	} else {
		fmt.Printf("\r%s%s %s%s", Cyan, s.frames[s.index-1], s.message, Gray)
	}
}

func (s *Spinner) IsRunning() bool {
	return s.isRunning
}

func (s *Spinner) Start() {
	fmt.Printf("\033[?25l") // Hide the cursor
	s.isRunning = true
	go func() {
		for {
			select {
			case <-s.done:
				return
			default:
				fmt.Printf("\r%s%s %s%s", Cyan, s.frames[s.index], s.message, Gray)
				s.index = (s.index + 1) % len(s.frames)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.done <- struct{}{}
	s.isRunning = false
	fmt.Print("\r") // Clear the spinner
	fmt.Println()
	fmt.Printf("\033[?25h") // Show the cursor
}
