package models

type ModelEnum string

const (
	OpenAIMini   ModelEnum = "gpt-4o-mini"
	OpenAI       ModelEnum = "gpt-4o"
	ClaudeSonnet ModelEnum = "sonnet"
	Maritaca     ModelEnum = "maritaca"
	Groq         ModelEnum = "groq"
	Ollama       ModelEnum = "ollama"
)

const DefaultModel = OpenAIMini

var availableModels = map[ModelEnum]string{
	OpenAIMini:   "GPT-4o-mini (OPENAI)",
	OpenAI:       "GPT-4o (OPENAI)",
	ClaudeSonnet: "Claude Sonnet 3.5 (ANTHROPIC)",
	Maritaca:     "Maritaca (BR)",
	Groq:         "Groq (LLAMA-3.3-70B)",
	Ollama:       "Ollama",
}

func (m ModelEnum) String() string {
	return availableModels[m]
}

func (m ModelEnum) Slug() string {
	return string(m)
}

func ParseFrom(modelSlug string) (ModelEnum, bool) {
	_, ok := availableModels[ModelEnum(modelSlug)]
	if !ok {
		return DefaultModel, false
	}

	return ModelEnum(modelSlug), true
}

func ListModels() map[ModelEnum]string {
	return availableModels
}
