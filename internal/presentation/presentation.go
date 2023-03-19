package presentation

import "github.com/charmbracelet/lipgloss"

const MaxWidth = 80
const CompPrompt = "◎ "
const UserPrompt = "▶ "

var MetadataStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#b8aec8"))

var BotStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#15fd00")).PaddingTop(1)

var UserStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fc0ff5")).PaddingTop(1)

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
