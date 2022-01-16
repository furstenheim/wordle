package main

import (
	"io/ioutil"
	"log"
	"strings"
)

const WORD_LENGTH = 5

type CombinationString string

type Combination byte

type CombinationArray [5]CombinationColor

type Word string

type ComparisonAggregate [253]int


func main () {
	precomputed := computePrecomputed()

}

func computePrecomputed () Precomputed {
	file, openErr := ioutil.ReadFile("words.en.txt")
	if openErr != nil {
		log.Fatal("Error opening", openErr)
	}
	split := strings.Split(string(file), "\n")
	dictionary := make([]Word, len(split))
	for i, v := range(split) {
		dictionary[i] = Word(v)
	}
	precomputed := Precomputed{
		Dictionary: dictionary,
		ComparisonAggregate: make([]ComparisonAggregate, len(dictionary)),
	}
	for i, v := range (dictionary) {
		ca := ComparisonAggregate{}
		for _, w := range(dictionary) {
			ca[computeCombination(w, v).toNumber()]++
		}
		precomputed.ComparisonAggregate[i] = ca
	}
	return precomputed
}

type Precomputed struct {
	Dictionary []Word
	ComparisonAggregate []ComparisonAggregate
}



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

func (c CombinationArray) toNumber () Combination {
	var power, res byte
	power = 1
	for _, v := range (c) {
		res += (byte(v)) *  power
		power *= 3
	}
	return Combination(res)
}

type CombinationColor byte

const (
	Grey CombinationColor = iota
	Yellow
	Green
)

