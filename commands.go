package main

import (
	"bufio"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (m chatModel) formatChatForMarkdown() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# System Prompt: %s\n", m.systemPrompt))
	sb.WriteString(fmt.Sprintf("## %s\n\n", m.chatStartTime.Format("2006-January-02 15:04:05")))
	for i := 0; i < len(m.userLines); i++ {
		sb.WriteString(fmt.Sprintf("### Human \n %s%s\n\n", humanUser.prompt, m.userLines[i]))
		if i < len(m.botLines) {
			sb.WriteString(fmt.Sprintf("### Bot \n %s%s\n\n", botUser.prompt, m.botLines[i]))
		}
	}
	return sb.String()
}

func (m chatModel) getFilename() string {
	var line string
	size := 30
	if len(m.userLines) > 0 {
		if len(m.userLines[0]) > size {
			line = m.userLines[0][:size]
		} else {
			line = m.userLines[0]
		}
	} else {
		if len(m.systemPrompt) > size {
			line = m.systemPrompt[:size]
		} else {
			line = m.systemPrompt
		}
	}
	return strings.ReplaceAll(line, " ", "_")
}

func (m chatModel) WriteChatToFile() tea.Cmd {
	now := time.Now()
	filename := fmt.Sprintf(
		"%s-%d-%d-%d-%s.txt",
		now.Format("2006-January-02"), now.Hour(), now.Minute(), now.Second(), m.getFilename(),
	)
	path := filepath.Join(os.Getenv("CHATBOT_LOGS"), filename)
	f, _ := os.Create(path)
	w := bufio.NewWriter(f)
	_, err := w.WriteString(m.formatChatForMarkdown())
	if err != nil {
		panic("Could not write to file")
	}
	err = w.Flush()
	if err != nil {
		panic("Could not flush buffer")
	}
	return tea.Quit
}

type botMsg ChatGPTResponse

func (m chatModel) DoBotMessage() tea.Msg {
	chatGptResponse := m.openAIClient.getChatGPTResponse(
		m.userLines,
		m.botLines,
		m.systemPrompt,
		m.linesToRemoveFromChatRequest,
		m.GPTModel,
	)
	if chatGptResponse.Usage.TotalTokens > m.GPTModel.MaxTokens {
		m.linesToRemoveFromChatRequest += 1
	}

	return botMsg(chatGptResponse)
}
