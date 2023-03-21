package presentation

import "github.com/charmbracelet/lipgloss"

const Width = 80
const BoxWidth = 120
const CompPrompt = "◎ "
const UserPrompt = "▶ "

var MetadataStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#b8aec8")).MarginBottom(1).MarginTop(1)

var botColor = lipgloss.Color("#15fd00")
var BotStyle = lipgloss.NewStyle().
	Foreground(botColor).
	Padding(1).
	Margin(1).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(botColor)

var userColor = lipgloss.Color("#fc0ff5")
var UserStyle = lipgloss.NewStyle().
	Foreground(userColor).
	Padding(1).
	Margin(1).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(userColor)

var PromptColor = lipgloss.Color("#6df1d8")
var PromptStyle = lipgloss.NewStyle().
	Foreground(PromptColor).
	MarginTop(1).
	MarginBottom(1)

type User struct {
	Prompt string
	Style  lipgloss.Style
}

var BotUser = User{
	Prompt: CompPrompt,
	Style:  BotStyle,
}

var HumanUser = User{
	Prompt: UserPrompt,
	Style:  UserStyle,
}
