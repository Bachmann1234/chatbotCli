package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func SayGoodBye() tea.Cmd {
	return tea.Quit
}

type botMsg string

func (m chatModel) DoBotMessage() tea.Msg {
	chatGbtResponse := getChatGBTResponse(m.userLines, m.botLines, m.systemPrompt, m.linesToRemoveFromChatRequest)
	if chatGbtResponse.Usage.TotalTokens > tokenThreshold {
		m.linesToRemoveFromChatRequest += 1
	}

	return botMsg(chatGbtResponse.Choices[0].Message.Content)
}
