package openai

type ClientI interface {
	GetChatGPTResponse(
		userLines []string,
		botLines []string,
		systemPrompt string,
		linesToDrop int,
		model GPTModel,
	) ChatGPTResponse
}
