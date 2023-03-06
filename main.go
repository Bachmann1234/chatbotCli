package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

const CompPrompt = "◈"
const UserPrompt = ">"

var BotStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#15fd00"))

var UserStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fc0ff5"))

type chatLog struct {
	userLines   []string
	botLines    []string
	currentLine string
	quitting    bool
}

func initialModel() chatLog {
	return chatLog{
		userLines:   []string{},
		botLines:    []string{},
		currentLine: "",
		quitting:    false,
	}
}

func (m chatLog) Init() tea.Cmd {
	return nil
}

func WriteUserLine(sb *strings.Builder, message string) {
	sb.WriteString(UserStyle.Render(fmt.Sprintf("%s %s", UserPrompt, message)))
	sb.WriteString("\n")
}

func WriteBotLine(sb *strings.Builder, message string) {
	sb.WriteString(BotStyle.Render(fmt.Sprintf("%s %s", CompPrompt, message)))
	sb.WriteString("\n")
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
		WriteUserLine(&sb, m.currentLine)
	}
	return sb.String()
}

func (m chatLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		currentKey := msg.String()

		switch currentKey {

		case "ctrl+c":
			m.quitting = true
			return m, SayGoodBye()

		case "enter":
			m.userLines = append(m.userLines, m.currentLine)
			m.botLines = append(m.botLines, "That's so cool!")
			m.currentLine = ""
			return m, nil
		case "tab":
			m.currentLine += "    "
			return m, nil
		case "backspace":
			if len(m.currentLine) > 0 {
				m.currentLine = m.currentLine[:len(m.currentLine)-1]
			}
			return m, nil
		default:
			m.currentLine += currentKey
			return m, nil
		}
	}
	return m, nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
