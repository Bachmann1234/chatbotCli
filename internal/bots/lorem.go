package bots

import (
	"gopkg.in/loremipsum.v1"
	"math/rand"
)

type LoremBot struct {
	Name           string
	loremGenerator *loremipsum.LoremIpsum
}

func (b *LoremBot) GetBotResponse(_ []string, _ []string, _ string) string {
	return b.loremGenerator.Sentences(rand.Intn(4) + 1)
}

func getLoremBot() ChatBotI {
	return &LoremBot{
		Name:           "lorem",
		loremGenerator: loremipsum.New(),
	}
}
