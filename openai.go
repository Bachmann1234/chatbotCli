package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIClientI interface {
	getChatGPTResponse(userLines []string, botLines []string, systemPrompt string, linesToDrop int) ChatGPTResponse
}

type OpenAIClient struct {
	apiKey string
}

const DefaultTokenThreshold = 3_700 // Max tokens is 4,096. We need some buffer for the response

type GPTModel struct {
	Name      string
	MaxTokens int
}

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

func getModel(name string) GPTModel {
	switch name {
	case "3.5":
		return GPTModel{
			Name:      "gpt-3.5-turbo",
			MaxTokens: 3_700,
		}
	case "4":
		return GPTModel{
			Name:      "gpt-4",
			MaxTokens: 8_192,
		}
	default:
		panic(fmt.Sprintf("Invalid model name %s", name))
	}
}

func (openAIClient OpenAIClient) getChatGPTResponse(userLines []string, botLines []string, systemPrompt string, linesToDrop int) ChatGPTResponse {
	client := &http.Client{}
	chatGptRequest := ChatGPTRequest{
		Model:    "gpt-3.5-turbo",
		Messages: constructMessages(userLines, botLines, systemPrompt, linesToDrop),
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
	req.Header.Add("Authorization", "Bearer "+openAIClient.apiKey)
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

func constructMessages(userLines []string, botLines []string, systemPrompt string, linesToDrop int) []ChatGPTMessage {
	var messages []ChatGPTMessage
	messages = append(messages, ChatGPTMessage{systemPrompt, "system"})
	for i := linesToDrop; i < len(userLines); i++ {
		messages = append(messages, ChatGPTMessage{userLines[i], "user"})
		if i < len(botLines) {
			messages = append(messages, ChatGPTMessage{botLines[i], "assistant"})
		}
	}
	return messages

}
