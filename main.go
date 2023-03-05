package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
)

const CompPrompt = "◈"
const UserPrompt = ">"

type chatLog struct {
	userLines   []string
	botLines    []string
	currentLine string
}

func initialModel() chatLog {
	return chatLog{
		userLines: []string{},

		botLines:    []string{},
		currentLine: "",
	}
}

func (m chatLog) Init() tea.Cmd {
	return nil
}

func (m chatLog) View() string {
	var sb strings.Builder

	for index, message := range m.userLines {
		sb.WriteString(fmt.Sprintf("%s %s\n", UserPrompt, message))
		sb.WriteString(fmt.Sprintf("%s %s\n", CompPrompt, m.botLines[index]))
	}
	sb.WriteString(fmt.Sprintf(
		"%s %s", UserPrompt, m.currentLine,
	))
	return sb.String()
}

func (m chatLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		currentKey := msg.String()

		switch currentKey {

		case "ctrl+c":
			return m, tea.Quit

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
