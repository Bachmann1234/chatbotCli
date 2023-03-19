package bot_metadata

import (
	"dev/mattbachmann/chatbotcli/internal/bots"
	"dev/mattbachmann/chatbotcli/internal/presentation"
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
	sb.WriteString(presentation.MetadataStyle.Render("Metadata - "))
	keys := make([]string, 0, len(m.metadata))
	for k := range m.metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(presentation.MetadataStyle.Render(fmt.Sprintf("%s: %s ", k, m.metadata[k])))
	}
	sb.WriteString("\n")
	return sb.String()
}
