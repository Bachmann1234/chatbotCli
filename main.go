package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
)

type chatLog struct {
	userLines   []string
	botLines    []string
	currentLine textinput.Model
	quitting    bool
}

func initialModel() chatLog {
	userInput := textinput.New()
	userInput.TextStyle = humanUser.style
	userInput.Prompt = humanUser.prompt
	return chatLog{
		userLines:   []string{},
		botLines:    []string{},
		currentLine: userInput,
		quitting:    false,
	}
}

func (m chatLog) Init() tea.Cmd {
	return textinput.Blink
}

func WriteLine(sb *strings.Builder, message string, user User) {
	sb.WriteString(user.style.Render(fmt.Sprintf("%s%s", user.prompt, message)))
	sb.WriteString("\n")
}

func WriteUserLine(sb *strings.Builder, message string) {
	WriteLine(sb, message, humanUser)
}

func WriteBotLine(sb *strings.Builder, message string) {
	WriteLine(sb, message, botUser)
}

func (m chatLog) View() string {
	var sb strings.Builder
	WriteBotLine(&sb, "Hello!")
	for index, message := range m.userLines {
		WriteUserLine(&sb, message)
		WriteBotLine(&sb, m.botLines[index])
	}
	if m.quitting {
		WriteBotLine(&sb, "Goodbye!")
	} else {
		sb.WriteString(humanUser.style.Render(m.currentLine.View()))
	}
	return sb.String()
}

func (m chatLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {

	case tea.KeyMsg:
		currentKey := msg.String()

		switch currentKey {

		case "ctrl+c":
			m.quitting = true
			return m, SayGoodBye()

		case "enter":
			m.userLines = append(m.userLines, m.currentLine.Value())
			m.currentLine = textinput.New()
			m.botLines = append(m.botLines, "That's so cool!")
			return m, cmd
		}
		if !m.currentLine.Focused() {
			m.currentLine.Focus()
		}
		m.currentLine, cmd = m.currentLine.Update(msg)
	}
	return m, cmd
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
