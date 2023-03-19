package main

import (
	"dev/mattbachmann/chatbotcli/internal/bubbletea"
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	systemPrompt := flag.String(
		"prompt",
		"You are a helpful assistant. Respond in a way that renders nicely in Markdown.",
		"The initial prompted hinting at the personality of the bots",
	)
	modelName := flag.String(
		"model",
		"gpt3_5",
		"which internal to use, gpt3_5, gpt4, or lorem right now. Lorem is just for returning generic text",
	)
	flag.Parse()
	p := tea.NewProgram(bubbletea.InitialModel(*systemPrompt, *modelName))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
