package components

import (
	"bufio"
	"dev/mattbachmann/chatbotcli/internal/bot_metadata"
	"dev/mattbachmann/chatbotcli/internal/bots"
	"dev/mattbachmann/chatbotcli/internal/integrations/openai"
	"dev/mattbachmann/chatbotcli/internal/presentation"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ChatModel struct {
	systemPrompt   string
	userLines      []string
	botLines       []bots.BotResponse
	currentMessage textarea.Model
	viewport       viewport.Model
	metadata       tea.Model
	width          int
	quitting       bool
	spinner        spinner.Model
	chatStartTime  time.Time
	chatBot        bots.ChatBotI
	writingMessage bool
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

	ta := textarea.New()
	ta.Focus()

	ta.SetWidth(presentation.MaxWidth)
	ta.SetHeight(25)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(presentation.MaxWidth, 10)
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	metadata := bot_metadata.New()

	return ChatModel{
		systemPrompt:   systemPrompt,
		userLines:      []string{},
		botLines:       []bots.BotResponse{},
		metadata:       metadata,
		currentMessage: ta,
		viewport:       vp,
		quitting:       false,
		spinner:        s,
		chatBot:        GetAIModel(modelName),
		chatStartTime:  time.Now(),
		writingMessage: true,
		width:          presentation.MaxWidth,
	}
}

func GetAIModel(name string) bots.ChatBotI {
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
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

func isUserTurn(m ChatModel) bool {
	return len(m.userLines) == len(m.botLines)
}

func isBotTurn(m ChatModel) bool {
	return !isUserTurn(m)
}

func (m ChatModel) View() string {
	var sb strings.Builder
	if isBotTurn(m) {
		sb.WriteString(m.spinner.View())
	} else {
		sb.WriteString(m.metadata.View())
	}
	if m.writingMessage {
		sb.WriteString("\n")
		sb.WriteString(m.currentMessage.View())
	} else {
		sb.WriteString(m.viewport.View())
	}
	return sb.String()
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
	chatTime := m.chatStartTime
	filename := fmt.Sprintf(
		"%s-%d-%d-%d-%s.txt",
		chatTime.Format("2006-01-02"), chatTime.Hour(), chatTime.Minute(), chatTime.Second(), m.getFilename(),
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

func (m ChatModel) renderConversation() string {
	var sb strings.Builder
	WriteBotLine(&sb, fmt.Sprintf("Prompt: %s", m.systemPrompt))
	for index, message := range m.userLines {
		WriteUserLine(&sb, message)
		if index < len(m.botLines) {
			WriteBotLine(&sb, m.botLines[index].Content)
		}
	}
	width := m.width
	if width > presentation.MaxWidth {
		width = presentation.MaxWidth
	}
	return wordwrap.String(sb.String(), width)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)
	if m.writingMessage && isUserTurn(m) {
		m.currentMessage, tiCmd = m.currentMessage.Update(msg)
	} else {
		m.viewport, vpCmd = m.viewport.Update(msg)
	}
	m.metadata, cmd = m.metadata.Update(msg)
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		width := msg.Width
		if width > presentation.MaxWidth {
			width = presentation.MaxWidth
		}
		if isUserTurn(m) {
			m.currentMessage.SetWidth(width)
		}
		m.width = width
		m.viewport.Width = width
		m.viewport.Height = msg.Height - 5

	case tea.KeyMsg:
		currentKey := msg.String()

		switch currentKey {

		case "ctrl+c":
			m.quitting = true
			return m, m.WriteChatToFile()

		case "enter":
			if isUserTurn(m) {
				if m.writingMessage && m.currentMessage.Value() != "" {
					m.userLines = append(m.userLines, m.currentMessage.Value())
					m.currentMessage.Reset()
					m.viewport.GotoBottom()
					m.viewport.SetContent(m.renderConversation())
					m.writingMessage = false
					return m, m.DoBotMessage
				} else if !m.writingMessage {
					m.writingMessage = true
				} else {
					m.writingMessage = false
				}
			}
			return m, cmd
		}
	case bots.BotResponse:
		m.botLines = append(m.botLines, msg)
		m.currentMessage.Reset()
		m.viewport.GotoBottom()
		m.viewport.SetContent(m.renderConversation())
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return m, tea.Batch(tiCmd, vpCmd, cmd)
}

func (m ChatModel) DoBotMessage() tea.Msg {
	return m.chatBot.GetBotResponse(m.userLines, m.botLines, m.systemPrompt)
}
