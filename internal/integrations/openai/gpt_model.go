package openai

import (
	"dev/mattbachmann/chatbotcli/internal/bots"
	"fmt"
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
			Client: Client{
				ApiKey: apiKey,
			},
		}
	case "gpt4":
		return GPTModel{
			Name:      "gpt-4",
			MaxTokens: 8_192,
			Client: Client{
				ApiKey: apiKey,
			},
		}
	default:
		return nil
	}
}

type GPTModel struct {
	Name      string
	MaxTokens int
	Client    ClientI
}

func (gptModel GPTModel) GetBotResponse(userLines []string, botLines []bots.BotResponse, systemPrompt string) bots.BotResponse {
	messagesToCut := determineMessagesToCut(botLines, gptModel.MaxTokens)
	chatGPTResponse := gptModel.Client.GetChatGPTResponse(
		userLines,
		botLines,
		systemPrompt,
		messagesToCut,
		gptModel,
	)
	return bots.BotResponse{
		Content: chatGPTResponse.Choices[0].Message.Content,
		Metadata: map[string]string{
			MessagesCut: strconv.Itoa(messagesToCut),
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
