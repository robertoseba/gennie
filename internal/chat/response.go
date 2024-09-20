package chat

import "time"

// Data sent to the model
type Input struct {
	Content   string
	Role      string
	Timestamp time.Time
}

// Data returned by the model
type Output struct {
	Content   string
	Role      string
	Timestamp time.Time
}

type Response struct {
	Answer   Output
	Question Input
}
