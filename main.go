package main

const WORD_LENGTH = 5

type CombinationString string

type Combination byte

type CombinationArray [5]CombinationColor

type Word string

func computeCombination (input, solution Word) CombinationArray {
	count := map[byte]int{}
	if len(input) != WORD_LENGTH || len(solution) != WORD_LENGTH {
		panic("Unexpected length for string")
	}
	res := CombinationArray{}

	for i := 0; i < WORD_LENGTH; i++ {
		if input[i] == solution[i] {
			res[i] = Green
		} else {
			count[solution[i]]++
		}
	}
	for i := 0; i < WORD_LENGTH; i++ {
		if input[i] != solution[i] && count[input[i]] > 0 {
			count[input[i]]--
			res[i] = Yellow
		}
	}

	return res
}

type CombinationColor byte

const (
	Grey CombinationColor = iota
	Yellow
	Green
)

