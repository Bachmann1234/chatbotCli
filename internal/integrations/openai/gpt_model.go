package openai

import (
	"dev/mattbachmann/chatbotcli/internal/bots"
	"fmt"
	"github.com/shopspring/decimal"
	"os"
	"strconv"
)

const MessagesCut = "messagesCut"
const TokensUsed = "tokensUsed"
const ConversationCost = "conversationCost"
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
			PricePer1KTokens: "0.002",
		}
	case "gpt4":
		return GPTModel{
			Name:      "gpt-4",
			MaxTokens: 8_192,
			Client: Client{
				ApiKey: apiKey,
			},
			PricePer1KTokens: "0.06",
		}
	case "gpt4-32":
		return GPTModel{
			Name:      "gpt-4-32k",
			MaxTokens: 32_768,
			Client: Client{
				ApiKey: apiKey,
			},
			PricePer1KTokens: "0.12",
		}
	default:
		return nil
	}
}

type GPTModel struct {
	Name             string
	MaxTokens        int
	Client           ClientI
	PricePer1KTokens string
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
	totalConversationCost := gptModel.computeTotalConversationCost(chatGPTResponse.Usage.TotalTokens, botLines)
	return bots.BotResponse{
		Content: chatGPTResponse.Choices[0].Message.Content,
		Metadata: map[string]string{
			MessagesCut:      strconv.Itoa(messagesToCut),
			TokensUsed:       strconv.Itoa(chatGPTResponse.Usage.TotalTokens),
			ConversationCost: totalConversationCost,
		},
	}
}

func (gptModel GPTModel) computeTotalConversationCost(totalTokens int, botLines []bots.BotResponse) string {
	costPer1kTokens, err := decimal.NewFromString(gptModel.PricePer1KTokens)
	if err != nil {
		panic(err)
	}
	tokensInRequest := decimal.NewFromInt(int64(totalTokens))
	costForRequest := costPer1kTokens.Mul(tokensInRequest).Div(decimal.NewFromInt(1000))

	if len(botLines) > 0 {
		lastBotLine := botLines[len(botLines)-1]
		lastConversationCost, err := decimal.NewFromString(lastBotLine.Metadata[ConversationCost])
		if err != nil {
			panic(fmt.Sprintf("Bad bot_metadata for conversationCost %s", lastBotLine.Metadata[ConversationCost]))
		}
		costForRequest = costForRequest.Add(lastConversationCost)
	}

	return costForRequest.String()
}

func determineMessagesToCut(botLines []bots.BotResponse, maxTokens int) int {
	messagesToCut := 0
	if len(botLines) > 0 {
		lastBotLine := botLines[len(botLines)-1]
		lastLinesToCut, err := strconv.Atoi(lastBotLine.Metadata[MessagesCut])
		if err != nil {
			panic(fmt.Sprintf("Bad bot_metadata for messagesCut %s", lastBotLine.Metadata[MessagesCut]))
		}
		messagesToCut = lastLinesToCut

		lastTokensUsed, err := strconv.Atoi(lastBotLine.Metadata[TokensUsed])
		if err != nil {
			panic(fmt.Sprintf("Bad bot_metadata for tokensUsed %s", lastBotLine.Metadata[TokensUsed]))
		}

		if lastTokensUsed >= maxTokens-Buffer {
			messagesToCut += 1
		}
	}
	return messagesToCut
}
