package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func SayGoodBye() tea.Cmd {
	return tea.Quit
}

func DoBotMessage() tea.Msg {
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C
	return botMsg("Oh, wow, Cool!")
}

type botMsg string
