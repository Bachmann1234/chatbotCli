package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const maxTokens = 4_096
const tokenThreshold = 3_700

type ChatGBTRequest struct {
	Model    string           `json:"model"`
	Messages []ChatGBTMessage `json:"messages"`
}

type ChatGBTMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatGBTResponse struct {
	Id      string          `json:"id"`
	Object  string          `json:"object"`
	Created uint64          `json:"created"`
	Choices []ChatGBTChoice `json:"choices"`
	Usage   ChatGBTUsage    `json:"usage"`
}

type ChatGBTChoice struct {
	Index        int            `json:"index"`
	Message      ChatGBTMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

type ChatGBTUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func getChatGBTResponse(userLines []string, botLines []string, systemPrompt string, linesToDrop int) ChatGBTResponse {
	client := &http.Client{}
	chatGbtRequest := ChatGBTRequest{
		Model:    "gpt-3.5-turbo",
		Messages: constructMessages(userLines, botLines, systemPrompt, linesToDrop),
	}
	postBody, err := json.Marshal(chatGbtRequest)
	if err != nil {
		panic(err)
	}
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		requestBody,
	)
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
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
	var chatGBTResponse ChatGBTResponse
	err = json.Unmarshal(body, &chatGBTResponse)
	if err != nil {
		panic(err)
	}
	return chatGBTResponse
}

func constructMessages(userLines []string, botLines []string, systemPrompt string, linesToDrop int) []ChatGBTMessage {
	var messages []ChatGBTMessage
	messages = append(messages, ChatGBTMessage{systemPrompt, "system"})
	for i := linesToDrop; i < len(userLines); i++ {
		messages = append(messages, ChatGBTMessage{userLines[i], "user"})
		if i < len(botLines) {
			messages = append(messages, ChatGBTMessage{botLines[i], "assistant"})
		}
	}
	return messages

}
