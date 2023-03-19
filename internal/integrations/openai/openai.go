package openai

import (
	"bytes"
	"dev/mattbachmann/chatbotcli/internal/bots"
	"encoding/json"
	"io"
	"net/http"
)

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

type Client struct {
	ApiKey string
}

func (openaiClient Client) GetChatGPTResponse(
	userLines []string,
	botLines []bots.BotResponse,
	systemPrompt string,
	messagesToCut int,
	gptModel GPTModel,
) ChatGPTResponse {
	client := &http.Client{}
	messages := ConstructMessages(userLines, botLines, systemPrompt, messagesToCut)
	chatGptRequest := ChatGPTRequest{
		Model:    gptModel.Name,
		Messages: messages,
	}
	postBody, err := json.Marshal(chatGptRequest)
	if err != nil {
		panic(err)
	}
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		requestBody,
	)
	req.Header.Add("Authorization", "Bearer "+openaiClient.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var chatGPTResponse ChatGPTResponse
	err = json.Unmarshal(body, &chatGPTResponse)
	if err != nil {
		panic(err)
	}
	return chatGPTResponse
}

func ConstructMessages(userLines []string, botLines []bots.BotResponse, systemPrompt string, messagesToCut int) []ChatGPTMessage {
	var messages []ChatGPTMessage
	messages = append(messages, ChatGPTMessage{systemPrompt, "system"})
	for i := messagesToCut; i < len(userLines); i++ {
		messages = append(messages, ChatGPTMessage{userLines[i], "user"})
		if i < len(botLines) {
			messages = append(messages, ChatGPTMessage{botLines[i].Content, "assistant"})
		}
	}
	return messages

}
