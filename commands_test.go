package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"testing"
)

type MockOpenAIClient struct {
	apiKey string
}

func (openAIClient MockOpenAIClient) getChatGBTResponse(userLines []string, botLines []string, systemPrompt string, linesToDrop int) ChatGBTResponse {
	return ChatGBTResponse{
		Id:      "bla",
		Object:  "bla",
		Created: 0,
		Choices: []ChatGBTChoice{{
			Index: 0,
			Message: ChatGBTMessage{
				Content: "Howdy!",
				Role:    "assistant",
			},
			FinishReason: "stop",
		}},
		Usage: ChatGBTUsage{
			PromptTokens:     4000,
			CompletionTokens: 4000,
			TotalTokens:      8000,
		},
	}
}

func TestIncrementingLinesToRemoveWhenUsedTokensHigh(t *testing.T) {
	model := chatModel{
		userLines:                    []string{"user"},
		botLines:                     []string{"bot"},
		currentLine:                  textinput.New(),
		quitting:                     false,
		spinner:                      spinner.New(),
		systemPrompt:                 "System prompt",
		linesToRemoveFromChatRequest: 0,
		tokenThresholdBeforeDropping: DefaultTokenThreshold,
		openAIClient: MockOpenAIClient{
			apiKey: "Test",
		},
	}

	newModel, cmd := model.Update(
		botMsg(
			model.openAIClient.getChatGBTResponse(
				model.userLines,
				model.botLines,
				model.systemPrompt,
				model.linesToRemoveFromChatRequest),
		),
	)

	if newModel.(chatModel).linesToRemoveFromChatRequest != 1 {
		t.Errorf("Expected linesToRemoveFromChatRequest to be 1, got %d", model.linesToRemoveFromChatRequest)
	}

	if cmd != nil {
		t.Error("Expected cmd to be nil")
	}

}
