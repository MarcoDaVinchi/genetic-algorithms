package bot

import "math/rand"

type Genome struct {
	values [128]int
}

func NewRandomGenome() Genome {
	values := [128]int{}
	for i := 0; i < 128; i++ {
		values[i] = rand.Intn(128)
	}
	return Genome{values: values}
}
