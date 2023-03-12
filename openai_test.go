package main

import "testing"

func TestConstructMessages(t *testing.T) {
	// Verify that the messages are constructed correctly
	var result = constructMessages(
		[]string{"Hello", "How are you?", "That's Great"},
		[]string{"Howdy", "good!"},
		"System Prompt",
		0,
	)

	expected := []ChatGBTMessage{
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

func TestConstructMessagesDropOne(t *testing.T) {
	// Verify that the messages are constructed correctly
	var result = constructMessages(
		[]string{"Hello", "How are you?", "That's Great"},
		[]string{"Howdy", "good!"},
		"System Prompt",
		1,
	)

	expected := []ChatGBTMessage{
		{"System Prompt", "system"},
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

func TestConstructMessagesDropThree(t *testing.T) {
	// Verify that the messages are constructed correctly
	var result = constructMessages(
		[]string{"Hello", "How are you?", "That's Great"},
		[]string{"Howdy", "good!"},
		"System Prompt",
		3,
	)

	expected := []ChatGBTMessage{
		{"System Prompt", "system"},
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
