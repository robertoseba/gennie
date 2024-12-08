package mocks

func NewMockOpenAIResponse(answer string) []byte {
	return []byte(`
	{
    "id": "chatcmpl-abc123",
    "object": "chat.completion",
    "created": 1677858242,
    "model": "gpt-4o-mini",
    "usage": {
        "prompt_tokens": 13,
        "completion_tokens": 7,
        "total_tokens": 20,
        "completion_tokens_details": {
            "reasoning_tokens": 0,
            "accepted_prediction_tokens": 0,
            "rejected_prediction_tokens": 0
        }
    },
    "choices": [
        {
            "message": {
                "role": "assistant",
                "content":"` + answer + `" 
            },
            "logprobs": null,
            "finish_reason": "stop",
            "index": 0
        }
    ]}`)
}
