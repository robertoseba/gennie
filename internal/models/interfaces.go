package models

type IModel interface {
	ask(question string, history *ChatHistory) Response
}
