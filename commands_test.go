package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"testing"
)

type MockOpenAIClient struct {
	apiKey string
}

func (openAIClient MockOpenAIClient) getChatGPTResponse(
	userLines []string,
	botLines []string,
	systemPrompt string,
	linesToDrop int,
	model GPTModel,
) ChatGPTResponse {
	return ChatGPTResponse{
		Id:      "bla",
		Object:  "bla",
		Created: 0,
		Choices: []ChatGPTChoice{{
			Index: 0,
			Message: ChatGPTMessage{
				Content: "Howdy!",
				Role:    "assistant",
			},
			FinishReason: "stop",
		}},
		Usage: ChatGPTUsage{
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
		GPTModel:                     getModel("3.5"),
		openAIClient: MockOpenAIClient{
			apiKey: "Test",
		},
	}

	newModel, cmd := model.Update(
		botMsg(
			model.openAIClient.getChatGPTResponse(
				model.userLines,
				model.botLines,
				model.systemPrompt,
				model.linesToRemoveFromChatRequest,
				model.GPTModel,
			),
		),
	)

	if newModel.(chatModel).linesToRemoveFromChatRequest != 1 {
		t.Errorf("Expected linesToRemoveFromChatRequest to be 1, got %d", model.linesToRemoveFromChatRequest)
	}

	if cmd != nil {
		t.Error("Expected cmd to be nil")
	}

}
