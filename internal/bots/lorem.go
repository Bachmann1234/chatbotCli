package bots

import (
	"gopkg.in/loremipsum.v1"
	"math/rand"
	"strconv"
)

type LoremBot struct {
	Name           string
	loremGenerator *loremipsum.LoremIpsum
}

func (b *LoremBot) GetBotResponse(_ []string, _ []string, _ string) BotResponse {
	numSentences := rand.Intn(4) + 1
	return BotResponse{
		Content: b.loremGenerator.Sentences(numSentences),
		Metadata: map[string]string{
			"numSentences": strconv.Itoa(numSentences),
		},
	}
}

func getLoremBot() ChatBotI {
	return &LoremBot{
		Name:           "lorem",
		loremGenerator: loremipsum.New(),
	}
}
