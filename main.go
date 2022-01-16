package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

const WORD_LENGTH = 5

type CombinationString string

type Combination byte

type CombinationArray [5]CombinationColor

type Word string

type ComparisonAggregate [253]int

type SharedInput struct {
	combinations []Combination
}

func main () {
	// precomputed := computePrecomputed()
	inputFile, readErr := ioutil.ReadFile("input.txt")
	if readErr != nil {
		log.Fatal("Error on input", readErr)
	}
	inputsRegex := regexp.MustCompile("Wordle \\d{3} \\d/\\d\n\n((\U0001f7e9|\U0001f7e8|\u2b1b|\u2b1c){5}\n)+")
	inputs := inputsRegex.FindAllString(string(inputFile), -1)
	sharedInputs := []SharedInput{}
	for _, v := range(inputs) {
		log.Println(string(v))
		sa := SharedInput{combinations: []Combination{}}
		split := strings.Split(string(v), "\n")

		for _, l := range(split[2:len(split)-1]) {
			sa.combinations = append(sa.combinations, CombinationString(l).toCombination())
		}
		sharedInputs = append(sharedInputs, sa)
	}
	log.Println(sharedInputs)
}

func (c CombinationString) toCombination () Combination {
	ca := CombinationArray{}
	var i = 0
	var v rune
	for _, v = range(c) {
		if v == 11035 {
			ca[i] = Grey
		} else if v == 11036 {
			ca[i] = Grey
		} else if v == 129000 {
			ca[i] = Yellow
		} else if v == 129001 {
			ca[i] = Green
		} else {
			log.Fatal("Unknown color", string(v), v)
		}
		i++
	}
	if i != 5 {
		log.Fatal("Wrong length for combination", i)
	}
	return ca.toNumber()
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

