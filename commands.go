package main

import (
	"bufio"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
	"time"
)

func (m chatModel) formatChatForFile() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("System Prompt: %s\n***\n\n", m.systemPrompt))
	for i := 0; i < len(m.userLines); i++ {
		sb.WriteString(fmt.Sprintf("%s%s\n\n", humanUser.prompt, m.userLines[i]))
		if i < len(m.botLines) {
			sb.WriteString(fmt.Sprintf("%s%s\n\n", botUser.prompt, m.botLines[i]))
		}
	}
	return sb.String()
}

func (m chatModel) WriteChatToFile() tea.Cmd {
	f, _ := os.Create(fmt.Sprintf("./%d.txt", time.Now().Unix()))
	w := bufio.NewWriter(f)
	_, err := w.WriteString(m.formatChatForFile())
	if err != nil {
		panic("Could not write to file")
	}
	err = w.Flush()
	if err != nil {
		panic("Could not flush buffer")
	}
	return tea.Quit
}

type botMsg ChatGBTResponse

func (m chatModel) DoBotMessage() tea.Msg {
	chatGbtResponse := m.openAIClient.getChatGBTResponse(m.userLines, m.botLines, m.systemPrompt, m.linesToRemoveFromChatRequest)
	if chatGbtResponse.Usage.TotalTokens > m.tokenThresholdBeforeDropping {
		m.linesToRemoveFromChatRequest += 1
	}

	return botMsg(chatGbtResponse)
}
