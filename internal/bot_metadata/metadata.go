package bot_metadata

import (
	"dev/mattbachmann/chatbotcli/internal/bots"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"sort"
	"strings"
)

type Model struct {
	metadata map[string]string
}

func New() tea.Model {
	return Model{
		metadata: map[string]string{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case bots.BotResponse:
		m.metadata = msg.Metadata
	}
	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder
	keys := make([]string, 0, len(m.metadata))
	for k := range m.metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) > 0 {
		sb.WriteString("Metadata - ")
	}
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s: %s ", k, m.metadata[k]))
	}
	return sb.String()
}
