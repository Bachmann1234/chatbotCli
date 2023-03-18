package bots

type ChatBotI interface {
	GetBotResponse(
		userLines []string,
		botLines []string,
		systemPrompt string,
	) string
}
