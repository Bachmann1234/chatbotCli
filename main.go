package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"os"
	"strings"
	"time"
)

type chatModel struct {
	userLines                    []string
	botLines                     []string
	currentLine                  textinput.Model
	quitting                     bool
	spinner                      spinner.Model
	width                        int
	height                       int
	systemPrompt                 string
	linesToRemoveFromChatRequest int
	openAIClient                 OpenAIClientI
	chatStartTime                time.Time
	GPTModel                     GPTModel
}

func initialModel(systemPrompt string, modelName string) chatModel {
	userInput := textinput.New()
	userInput.TextStyle = humanUser.style
	userInput.Prompt = humanUser.prompt
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return chatModel{
		userLines:                    []string{},
		botLines:                     []string{},
		currentLine:                  userInput,
		quitting:                     false,
		spinner:                      s,
		systemPrompt:                 systemPrompt,
		linesToRemoveFromChatRequest: 0,
		GPTModel:                     getModel(modelName),
		openAIClient: OpenAIClient{
			apiKey: os.Getenv("OPENAI_API_KEY"),
		},
		chatStartTime: time.Now(),
	}
}

func isUserTurn(m chatModel) bool {
	return len(m.userLines) == len(m.botLines)
}

func isBotTurn(m chatModel) bool {
	return !isUserTurn(m)
}

func (m chatModel) Init() tea.Cmd {
	return m.spinner.Tick
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

func (m chatModel) View() string {
	var sb strings.Builder
	WriteBotLine(&sb, fmt.Sprintf("Prompt: %s", m.systemPrompt))
	for index, message := range m.userLines {
		WriteUserLine(&sb, message)
		if index < len(m.botLines) {
			WriteBotLine(&sb, m.botLines[index])
		}
	}
	if m.quitting {
		WriteBotLine(&sb, "Goodbye!")
	} else {
		if isBotTurn(m) {
			sb.WriteString(m.spinner.View())
		} else {
			sb.WriteString(humanUser.style.Render(m.currentLine.View()))
		}
	}
	width := m.width
	if width > maxWidth {
		width = maxWidth
	}
	return wordwrap.String(sb.String(), width)
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {

	case tea.KeyMsg:
		currentKey := msg.String()

		switch currentKey {

		case "ctrl+c":
			m.quitting = true
			return m, m.WriteChatToFile()

		case "enter":
			if isUserTurn(m) {
				m.userLines = append(m.userLines, m.currentLine.Value())
				m.currentLine = textinput.New()
				return m, m.DoBotMessage
			}
			return m, cmd
		}
		if !m.currentLine.Focused() {
			m.currentLine.Focus()
		}
		m.currentLine, cmd = m.currentLine.Update(msg)
	case botMsg:
		m.botLines = append(m.botLines, msg.Choices[0].Message.Content)
		if msg.Usage.TotalTokens > m.GPTModel.MaxTokens {
			m.linesToRemoveFromChatRequest += 1
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return m, cmd
}

func main() {
	systemPrompt := flag.String(
		"prompt",
		"You are a helpful assistant",
		"The initial prompted hinting at the personality of the chatbot",
	)
	modelName := flag.String(
		"model",
		"3.5",
		"which model to use, 3.5 or 4 right now",
	)
	flag.Parse()
	p := tea.NewProgram(initialModel(*systemPrompt, *modelName))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
