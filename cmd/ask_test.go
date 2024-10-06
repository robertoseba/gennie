package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/robertoseba/gennie/internal/cache"
	"github.com/robertoseba/gennie/internal/chat"
	mock_httpclient "github.com/robertoseba/gennie/internal/httpclient/mock"
	"github.com/robertoseba/gennie/internal/output"
	"github.com/robertoseba/gennie/internal/profile"
	"go.uber.org/mock/gomock"
)

func TestHasAllFlags(t *testing.T) {
	r := NewAskCmd(nil, nil, nil)
	if r.Use != "ask [question for the llm model]" {
		t.Errorf("Expected 'ask' but got %s", r.Use)
	}
	expectedFlags := []string{"followup", "append", "model", "profile"}

	for _, f := range expectedFlags {
		if r.Flags().Lookup(f) == nil {
			t.Errorf("Expected flag %s not found", f)
		}
	}
}

func TestSavesChatToCache(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock_httpclient.NewMockIHttpClient(ctrl)

	httpResponse := mockOpenAiResponse("The meaning of life is 42")

	mockClient.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(httpResponse), nil)

	out := bytes.NewBufferString("")

	printer := output.NewPrinter(out, nil)

	cache := setupTestCache()

	c := NewAskCmd(cache, printer, mockClient)

	c.SetArgs([]string{"ask what is the meaning of life?"})
	c.Execute()

	actualResponse, _ := cache.ChatHistory.LastChat()

	if actualResponse.GetAnswer() != "The meaning of life is 42" {
		t.Errorf("Expected 'The meaning of life is 42' but got %s", actualResponse.GetAnswer())
	}
}

func TestUsesModelFromFlag(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock_httpclient.NewMockIHttpClient(ctrl)

	httpResponse := mockOpenAiResponse("The meaning of life is 42")

	body := `{"model":"gpt-4o","messages":[{"role":"system","content":"You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necesary."},{"role":"user","content":"ask what is the meaning of life?"}]}`
	mockClient.EXPECT().Post(gomock.Any(), body, gomock.Any()).Return([]byte(httpResponse), nil)

	out := bytes.NewBufferString("")

	printer := output.NewPrinter(out, nil)

	cache := setupTestCache()

	c := NewAskCmd(cache, printer, mockClient)

	c.SetArgs([]string{"ask what is the meaning of life?", "--model", "gpt-4o"})
	c.Execute()

	//Saves the model from the flag to the cache
	if cache.GetCurrModelSlug() != "gpt-4o" {
		t.Errorf("Expected model to be 'gpt-4o' but got %s", cache.GetCurrModelSlug())
	}
}

func TestAppendsFileContentToQuestion(t *testing.T) {

	fileContent := "This is a file content"

	f, err := os.Create("./test.csv")
	if err != nil {
		t.Fatal("Failed to create file")
	}
	f.WriteString(fileContent)
	f.Close()

	ctrl := gomock.NewController(t)

	mockClient := mock_httpclient.NewMockIHttpClient(ctrl)

	httpResponse := mockOpenAiResponse("The meaning of life is 42")

	body := `{"model":"gpt-4o-mini","messages":[{"role":"system","content":"You are a helpful cli assistant. Try to answer in a concise way providing the most relevant information. And examples when necesary."},{"role":"user","content":"ask what is the meaning of life?\nThis is a file content"}]}`
	mockClient.EXPECT().Post(gomock.Any(), body, gomock.Any()).Return([]byte(httpResponse), nil)

	out := bytes.NewBufferString("")

	printer := output.NewPrinter(out, nil)

	cache := setupTestCache()

	c := NewAskCmd(cache, printer, mockClient)

	c.SetArgs([]string{"ask what is the meaning of life?", "--append", "./test.csv"})
	c.Execute()

	os.Remove("./test.csv")
}

func TestResetsChatHistoryIfNotFollowUp(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock_httpclient.NewMockIHttpClient(ctrl)

	httpResponse := mockOpenAiResponse("The meaning of life is 42")

	mockClient.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(httpResponse), nil)

	out := bytes.NewBufferString("")

	printer := output.NewPrinter(out, nil)

	cache := setupTestCache()
	oldChat := chat.Chat{}
	oldChat.AddQuestion("Initial question")
	oldChat.AddAnswer("Answer to initial question")
	cache.ChatHistory.AddChat(oldChat)

	c := NewAskCmd(cache, printer, mockClient)

	c.SetArgs([]string{"ask what is the meaning of life?", "--model", "gpt-4o"})
	c.Execute()

	if cache.ChatHistory.Len() != 1 {
		t.Errorf("Expected chat history to have 1 item but got %d", cache.ChatHistory.Len())
	}

	chat, _ := cache.ChatHistory.LastChat()

	if chat.GetAnswer() != "The meaning of life is 42" || chat.GetQuestion() != "ask what is the meaning of life?" {
		t.Errorf("Expected chat history to have only the last question and answer but got %v", chat)
	}

}

func TestFollowUpAppendsToChatHistory(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock_httpclient.NewMockIHttpClient(ctrl)

	httpResponse := mockOpenAiResponse("The meaning of life is 42")

	body := `{"model":"gpt-4o-mini","messages":[{"role":"system","content":"you are a assistant"},{"role":"user","content":"Initial question"},{"role":"assistant","content":"Answer to initial question"},{"role":"user","content":"ask what is the meaning of life?"}]}`
	mockClient.EXPECT().Post(gomock.Any(), body, gomock.Any()).Return([]byte(httpResponse), nil)

	out := bytes.NewBufferString("")

	printer := output.NewPrinter(out, nil)

	cache := setupTestCache()
	cache.CurrProfile.Data = "you are a assistant"
	oldChat := chat.Chat{}
	oldChat.AddQuestion("Initial question")
	oldChat.AddAnswer("Answer to initial question")
	cache.ChatHistory.AddChat(oldChat)

	c := NewAskCmd(cache, printer, mockClient)

	c.SetArgs([]string{"ask what is the meaning of life?", "--followup"})
	c.Execute()

	if cache.ChatHistory.Len() != 2 {
		t.Errorf("Expected chat history to have 2 item but got %d", cache.ChatHistory.Len())
	}

	chats := cache.ChatHistory.Responses

	if chats[0].GetAnswer() != "Answer to initial question" || chats[1].GetAnswer() != "The meaning of life is 42" {
		t.Errorf("Failed to append to chat history. Expected chat history to have both answers, but got: %v", chats)
	}

}

func TestUsesProfileFromFlag(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock_httpclient.NewMockIHttpClient(ctrl)

	httpResponse := mockOpenAiResponse("The meaning of life is 41")

	mockClient.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(httpResponse), nil)

	out := bytes.NewBufferString("")

	printer := output.NewPrinter(out, nil)

	cache := setupTestCache()
	cache.CachedProfiles = map[string]profile.ProfileInfo{
		"stub_profile": {Filepath: "../internal/cache/test/stub.profile.json"},
	}

	c := NewAskCmd(cache, printer, mockClient)

	c.SetArgs([]string{"ask what is the meaning of life?", "-p=stub_profile"})
	c.Execute()

	if cache.GetCurrProfile().Data != "just a profile stub for testing" {
		t.Errorf("Expected profile to be loaded from stub but got %s", cache.GetCurrProfile())
	}
}

func TestUsesProfileFromCacheIfNotSpecified(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockClient := mock_httpclient.NewMockIHttpClient(ctrl)

	httpResponse := mockOpenAiResponse("The meaning of life is 41")

	mockClient.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(httpResponse), nil)

	out := bytes.NewBufferString("")

	printer := output.NewPrinter(out, nil)

	cache := setupTestCache()
	cache.CurrProfile = profile.Profile{Slug: "test"}

	c := NewAskCmd(cache, printer, mockClient)

	c.SetArgs([]string{"ask what is the meaning of life?"})
	c.Execute()

	if cache.GetCurrProfile().Slug != "test" {
		t.Errorf("Expected profile to be test but got %s", cache.GetCurrProfile().Slug)
	}
}

func setupTestCache() *cache.Storage {
	return cache.NewStorage(".temp")
}

func mockOpenAiResponse(answer string) string {
	base := `{
		"choices": [
			{
			"finish_reason": "stop",
			"index": 0,
			"message": {
				"content": "%s",
				"role": "assistant"
			},
			"logprobs": null
			}
		],
		"created": 1677664795,
		"id": "chatcmpl",
		"model": "gpt-4o-mini",
		"object": "chat.completion",
		"usage": {
			"completion_tokens": 17,
			"prompt_tokens": 57,
			"total_tokens": 74,
			"completion_tokens_details": {
			"reasoning_tokens": 0
			}
		}
		}`

	return fmt.Sprintf(base, answer)
}
