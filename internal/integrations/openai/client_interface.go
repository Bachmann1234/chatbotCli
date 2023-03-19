package openai

import "dev/mattbachmann/chatbotcli/internal/bots"

type ClientI interface {
	GetChatGPTResponse(
		userLines []string,
		botLines []bots.BotResponse,
		systemPrompt string,
		linesToDrop int,
		model GPTModel,
	) ChatGPTResponse
}
