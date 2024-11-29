package models

import "github.com/robertoseba/gennie/internal/conversation"

type IModel interface {
	/**
	* Complete chat receives a chat history with the last chat being the one that needs to be completed.
	* It also receives a system prompt that can be used to generate the answer.
	* Once succeded the model will fill out the answer in the last conversation.
	 */
	CompleteChat(chatHistory *conversation.Conversation, systemPrompt string) error
}

/**
* The Provider is only responsible for preparing the payload,
* formatting it accordinly to the model's requirements
* and parsing the response back to the system.
 */
type IModelProvider interface {
	PreparePayload(chatHistory *conversation.Conversation, systemPrompt string) (string, error)
	ParseResponse(response []byte) (string, error)
	GetHeaders() map[string]string
	GetUrl() string
}
