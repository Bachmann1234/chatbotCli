package main

import (
	"bytes"
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"os"
)

func SayGoodBye() tea.Cmd {
	return tea.Quit
}

func (m chatModel) DoBotMessage() tea.Msg {
	client := &http.Client{}
	chatGbtRequest := ChatGBTRequest{
		Model:    "gpt-3.5-turbo",
		Messages: constructMessages(m.userLines, m.botLines, m.systemPrompt),
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
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return botMsg(data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string))
}

type botMsg string

type ChatGBTRequest struct {
	Model    string           `json:"model"`
	Messages []ChatGBTMessage `json:"messages"`
}

type ChatGBTMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func constructMessages(userLines []string, botLines []string, systemPrompt string) []ChatGBTMessage {
	var messages []ChatGBTMessage
	messages = append(messages, ChatGBTMessage{systemPrompt, "system"})
	for i := 0; i < len(userLines); i++ {
		messages = append(messages, ChatGBTMessage{userLines[i], "user"})
		if i < len(botLines) {
			messages = append(messages, ChatGBTMessage{botLines[i], "assistant"})
		}
	}
	return messages

}
