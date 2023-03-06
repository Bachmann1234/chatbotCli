package main

import "github.com/charmbracelet/lipgloss"

const CompPrompt = "◈"
const UserPrompt = ">"

var BotStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#15fd00"))

var UserStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fc0ff5"))

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
