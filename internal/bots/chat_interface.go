package bots

type ChatBotI interface {
	GetBotResponse(
		userLines []string,
		botLines []string,
		systemPrompt string,
	) BotResponse
}

type BotResponse struct {
	Content  string
	Metadata map[string]string
}

func GetChatBot(name string) ChatBotI {
	switch name {
	case "lorem":
		return getLoremBot()
	default:
		return nil
	}
}
