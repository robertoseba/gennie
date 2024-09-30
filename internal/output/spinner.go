package output

import (
	"fmt"
	"time"
)

type Spinner struct {
	frames  []string
	index   int
	done    chan struct{}
	message string
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:  []string{"|", "/", "-", "\\"},
		index:   0,
		done:    make(chan struct{}),
		message: message,
	}
}

func (s *Spinner) Start() {
	fmt.Printf("\033[?25l") // Hide the cursor
	go func() {
		for {
			select {
			case <-s.done:
				return
			default:
				fmt.Printf("\r%s %s", s.frames[s.index], s.message)
				s.index = (s.index + 1) % len(s.frames)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.done <- struct{}{}
	fmt.Print("\r")         // Clear the spinner
	fmt.Printf("\033[?25h") // Show the cursor
}
