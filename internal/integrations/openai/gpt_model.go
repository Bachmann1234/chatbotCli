package openai

import (
	"bytes"
	"dev/mattbachmann/chatbotcli/internal/bots"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func GetGPTModel(name string) bots.ChatBotI {
	apiKey := os.Getenv("OPENAI_API_KEY")
	switch name {
	case "gpt3_5":
		return GPTModel{
			Name:      "gpt-3.5-turbo",
			MaxTokens: 4_096,
			ApiKey:    apiKey,
		}
	case "gpt4":
		return GPTModel{
			Name:      "gpt-4",
			MaxTokens: 8_192,
			ApiKey:    apiKey,
		}
	default:
		return nil
	}
}

type GPTModel struct {
	Name      string
	MaxTokens int
	ApiKey    string
}

func (gptModel GPTModel) GetBotResponse(userLines []string, botLines []string, systemPrompt string) string {
	chatGPTResponse := gptModel.getChatGPTResponse(
		userLines,
		botLines,
		systemPrompt,
	)
	return chatGPTResponse.Choices[0].Message.Content
}

func (gptModel GPTModel) getChatGPTResponse(
	userLines []string,
	botLines []string,
	systemPrompt string,
) ChatGPTResponse {
	client := &http.Client{}
	chatGptRequest := ChatGPTRequest{
		Model:    gptModel.Name,
		Messages: ConstructMessages(userLines, botLines, systemPrompt),
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
	req.Header.Add("Authorization", "Bearer "+gptModel.ApiKey)
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
