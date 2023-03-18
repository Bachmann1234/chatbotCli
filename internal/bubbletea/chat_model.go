package bubbletea

import (
	"bufio"
	"dev/mattbachmann/chatbotcli/internal/bots"
	"dev/mattbachmann/chatbotcli/internal/integrations/openai"
	"dev/mattbachmann/chatbotcli/internal/presentation"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BotMsg struct {
	Content string
}

type ChatModel struct {
	userLines                    []string
	botLines                     []string
	currentLine                  textinput.Model
	quitting                     bool
	spinner                      spinner.Model
	width                        int
	height                       int
	systemPrompt                 string
	linesToRemoveFromChatRequest int
	openAIClient                 openai.ClientI
	chatStartTime                time.Time
	chatBot                      bots.ChatBotI
}

func WriteLine(sb *strings.Builder, message string, user presentation.User) {
	sb.WriteString(user.Style.Render(fmt.Sprintf("%s%s", user.Prompt, message)))
	sb.WriteString("\n")
}

func WriteUserLine(sb *strings.Builder, message string) {
	WriteLine(sb, message, presentation.HumanUser)
}

func WriteBotLine(sb *strings.Builder, message string) {
	WriteLine(sb, message, presentation.BotUser)
}

func InitialModel(systemPrompt string, modelName string) ChatModel {
	userInput := textinput.New()
	userInput.TextStyle = presentation.HumanUser.Style
	userInput.Prompt = presentation.HumanUser.Prompt
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return ChatModel{
		systemPrompt:  systemPrompt,
		userLines:     []string{},
		botLines:      []string{},
		currentLine:   userInput,
		quitting:      false,
		spinner:       s,
		chatBot:       GetModel(modelName),
		chatStartTime: time.Now(),
	}
}

func GetModel(name string) bots.ChatBotI {
	model := openai.GetGPTModel(name)
	if model == nil {
		model = bots.GetChatBot(name)
	}
	if model == nil {
		panic(fmt.Sprintf("Unknown model %s", name))
	}
	return model
}

func (m ChatModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func isUserTurn(m ChatModel) bool {
	return len(m.userLines) == len(m.botLines)
}

func isBotTurn(m ChatModel) bool {
	return !isUserTurn(m)
}

func (m ChatModel) View() string {
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
			sb.WriteString(presentation.HumanUser.Style.Render(m.currentLine.View()))
		}
	}
	width := m.width
	if width > presentation.MaxWidth {
		width = presentation.MaxWidth
	}
	return wordwrap.String(sb.String(), width)
}

func (m ChatModel) formatChatForMarkdown() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# System Prompt: %s\n", m.systemPrompt))
	sb.WriteString(fmt.Sprintf("## %s\n\n", m.chatStartTime.Format("2006-January-02 15:04:05")))
	for i := 0; i < len(m.userLines); i++ {
		sb.WriteString(fmt.Sprintf("### Human \n %s%s\n\n", presentation.HumanUser.Prompt, m.userLines[i]))
		if i < len(m.botLines) {
			sb.WriteString(fmt.Sprintf("### Bot \n %s%s\n\n", presentation.BotUser.Prompt, m.botLines[i]))
		}
	}
	return sb.String()
}

func (m ChatModel) getFilename() string {
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

func (m ChatModel) WriteChatToFile() tea.Cmd {
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

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case BotMsg:
		m.botLines = append(m.botLines, msg.Content)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return m, cmd
}

func (m ChatModel) DoBotMessage() tea.Msg {
	chatGptResponse := m.chatBot.GetBotResponse(m.userLines, m.botLines, m.systemPrompt)
	return BotMsg{
		Content: chatGptResponse,
	}
}
