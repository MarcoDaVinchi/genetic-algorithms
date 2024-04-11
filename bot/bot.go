package bot

import (
	"fmt"
	"strings"
)

type Bot struct {
	BotGenome Genome
}

func (b Bot) String() string {
	var genomeStrings []string
	for _, gene := range b.BotGenome.values {
		genomeStrings = append(genomeStrings, fmt.Sprintf("%d", gene))
	}
	return fmt.Sprintf("Bot genome: [%s]", strings.Join(genomeStrings, ","))
}

func NewBotWithRandomGenome() Bot {
	return Bot{
		BotGenome: NewRandomGenome(),
	}
}

func NewBotWithGenome(genome Genome) Bot {
	return Bot{
		BotGenome: genome,
	}
}

func (b Bot) GetGenome() Genome {
	return b.BotGenome
}
