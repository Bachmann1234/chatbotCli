package openai

import (
	"dev/mattbachmann/chatbotcli/internal/integrations/openai"
	"testing"
)

func TestConstructMessages(t *testing.T) {
	// Verify that the messages are constructed correctly
	var result = openai.ConstructMessages(
		[]string{"Hello", "How are you?", "That's Great"},
		[]string{"Howdy", "good!"},
		"System Prompt",
	)

	expected := []openai.ChatGPTMessage{
		{"System Prompt", "system"},
		{"Hello", "user"},
		{"Howdy", "assistant"},
		{"How are you?", "user"},
		{"good!", "assistant"},
		{"That's Great", "user"},
	}

	if len(result) != len(expected) {
		t.Errorf("Result Length (%d) not equal to expected length (%d)", len(result), len(expected))
	}
	for index, element := range result {
		if element != expected[index] {
			t.Errorf("Result (%s) not equal to expected (%s)", element, expected[index])
		}
	}

}
