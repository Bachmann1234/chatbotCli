package main

import "github.com/charmbracelet/lipgloss"

const maxWidth = 80
const CompPrompt = "◈ "
const UserPrompt = "> "

var BotStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#15fd00")).PaddingTop(1)

var UserStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fc0ff5")).PaddingTop(1)

type User struct {
	prompt string
	style  lipgloss.Style
}

var botUser = User{
	prompt: CompPrompt,
	style:  BotStyle,
}

var humanUser = User{
	prompt: UserPrompt,
	style:  UserStyle,
}

var SystemMsg = ChatGBTMessage{"You are a helpful assistant", "system"}
