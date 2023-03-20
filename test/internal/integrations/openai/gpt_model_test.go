package openai

import (
	"dev/mattbachmann/chatbotcli/internal/bots"
	"dev/mattbachmann/chatbotcli/internal/integrations/openai"
	"testing"
)

type MockOpenAIClient struct {
}

func (openAIClient MockOpenAIClient) GetChatGPTResponse(
	_ []string,
	_ []bots.BotResponse,
	_ string,
	_ int,
	_ openai.GPTModel,
) openai.ChatGPTResponse {
	return openai.ChatGPTResponse{
		Id:      "bla",
		Object:  "bla",
		Created: 0,
		Choices: []openai.ChatGPTChoice{{
			Index: 0,
			Message: openai.ChatGPTMessage{
				Content: "Howdy!",
				Role:    "assistant",
			},
			FinishReason: "stop",
		}},
		Usage: openai.ChatGPTUsage{
			PromptTokens:     4000,
			CompletionTokens: 4000,
			TotalTokens:      8000,
		},
	}
}

func TestIncrementingLinesToRemoveWhenUsedTokensHigh(t *testing.T) {
	model := openai.GPTModel{
		Name:             "test",
		MaxTokens:        1000,
		Client:           MockOpenAIClient{},
		PricePer1KTokens: "0.00",
	}

	response := model.GetBotResponse(
		[]string{"hi", "howdy"},
		[]bots.BotResponse{
			{
				Content: "hello",
				Metadata: map[string]string{
					openai.TokensUsed:  "10",
					openai.MessagesCut: "0",
				},
			},
			{
				Content: "Cool message",
				Metadata: map[string]string{
					openai.TokensUsed:  "1000",
					openai.MessagesCut: "0",
				},
			},
		},
		"System prompt",
	)

	if response.Metadata[openai.MessagesCut] != "1" {
		t.Errorf("Expected lines to remove to be 1, got %s", response.Metadata[openai.MessagesCut])
	}

}

func TestLeavingTokensAloneWhenUsedLow(t *testing.T) {
	model := openai.GPTModel{
		Name:             "test",
		MaxTokens:        1000,
		Client:           MockOpenAIClient{},
		PricePer1KTokens: "0.00",
	}

	response := model.GetBotResponse(
		[]string{"hi"},
		[]bots.BotResponse{{
			Content: "hello",
			Metadata: map[string]string{
				openai.TokensUsed:  "10",
				openai.MessagesCut: "0",
			},
		}},
		"System prompt",
	)

	if response.Metadata[openai.MessagesCut] != "0" {
		t.Errorf("Expected lines to remove to be 0, got %s", response.Metadata[openai.MessagesCut])
	}

}
