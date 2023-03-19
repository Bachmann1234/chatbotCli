package openai

import "dev/mattbachmann/chatbotcli/internal/bots"

type ChatGPTRequest struct {
	Model    string           `json:"model"`
	Messages []ChatGPTMessage `json:"messages"`
}

type ChatGPTMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatGPTResponse struct {
	Id      string          `json:"id"`
	Object  string          `json:"object"`
	Created uint64          `json:"created"`
	Choices []ChatGPTChoice `json:"choices"`
	Usage   ChatGPTUsage    `json:"usage"`
}

type ChatGPTChoice struct {
	Index        int            `json:"index"`
	Message      ChatGPTMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

type ChatGPTUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func ConstructMessages(userLines []string, botLines []bots.BotResponse, systemPrompt string) []ChatGPTMessage {
	var messages []ChatGPTMessage
	messages = append(messages, ChatGPTMessage{systemPrompt, "system"})
	for i := 0; i < len(userLines); i++ {
		messages = append(messages, ChatGPTMessage{userLines[i], "user"})
		if i < len(botLines) {
			messages = append(messages, ChatGPTMessage{botLines[i].Content, "assistant"})
		}
	}
	return messages

}
