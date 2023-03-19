package openai

import (
	"bytes"
	"dev/mattbachmann/chatbotcli/internal/bots"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const MessagesCut = "messagesCut"
const TokensUsed = "tokensUsed"
const Buffer = 500

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

func (gptModel GPTModel) GetBotResponse(userLines []string, botLines []bots.BotResponse, systemPrompt string) bots.BotResponse {
	chatGPTResponse, messagesCut := gptModel.getChatGPTResponse(
		userLines,
		botLines,
		systemPrompt,
	)
	return bots.BotResponse{
		Content: chatGPTResponse.Choices[0].Message.Content,
		Metadata: map[string]string{
			MessagesCut: strconv.Itoa(messagesCut),
			TokensUsed:  strconv.Itoa(chatGPTResponse.Usage.TotalTokens),
		},
	}
}

func determineMessagesToCut(botLines []bots.BotResponse, maxTokens int) int {
	messagesToCut := 0
	if len(botLines) > 0 {
		lastBotLine := botLines[len(botLines)-1]
		lastLinesToCut, err := strconv.Atoi(lastBotLine.Metadata[MessagesCut])
		if err != nil {
			panic(fmt.Sprintf("Bad metadata for messagesCut %s", lastBotLine.Metadata[MessagesCut]))
		}
		messagesToCut = lastLinesToCut

		lastTokensUsed, err := strconv.Atoi(lastBotLine.Metadata[TokensUsed])
		if err != nil {
			panic(fmt.Sprintf("Bad metadata for tokensUsed %s", lastBotLine.Metadata[TokensUsed]))
		}

		if lastTokensUsed >= maxTokens-Buffer {
			messagesToCut += 1
		}
	}

	return messagesToCut

}

func (gptModel GPTModel) getChatGPTResponse(
	userLines []string,
	botLines []bots.BotResponse,
	systemPrompt string,
) (ChatGPTResponse, int) {
	client := &http.Client{}
	messagesToCut := determineMessagesToCut(botLines, gptModel.MaxTokens)
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
	return chatGPTResponse, messagesToCut
}
