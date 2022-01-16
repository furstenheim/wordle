package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"bytes"
	"encoding/gob"
)

const WORD_LENGTH = 5
const PRECOMPUTED_FILE = "encoded-precomputed"
type CombinationString string

type Combination byte

type CombinationArray [5]CombinationColor

type Word string

type ComparisonAggregate [253]int

type SharedInput struct {
	combinations []Combination
}

func main () {
	precomputed := computePrecomputed()
	log.Println("computed precomputed")
	sharedInputs := parseInput()
	caInput := sharedInputs[0].toComparisonAggregate()

	for _, v := range sharedInputs[1:] {
		caInput = mergeComparisonAggregate(caInput, v.toComparisonAggregate())
	}
 	possibleWords := []Word{}

	for i, ca := range(precomputed.ComparisonAggregate) {
		if caInput.isCompatibleWith(ca) {
			possibleWords = append(possibleWords, precomputed.Dictionary[i])
		}
	}
	log.Println(possibleWords)
	log.Println(len(possibleWords))
}

func (ca1 ComparisonAggregate) isCompatibleWith (ca2 ComparisonAggregate) bool {
	for i, v1 := range(ca1) {
		v2 := ca2[i]
		if v2 < v1 {
			return false
		}
	}
	return true
}

func mergeComparisonAggregate (ca1, ca2 ComparisonAggregate) ComparisonAggregate {
	res := ComparisonAggregate{}
	for i, v1 := range(ca1) {
		v2 := ca2[i]
		res[i] = max(v1, v2)
	}
	return res
}

func max (a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (si SharedInput) toComparisonAggregate () ComparisonAggregate {
	ca := ComparisonAggregate{}
	for _, v := range (si.combinations) {
		ca[v]++
	}
	return ca
}



func parseInput () []SharedInput {
	inputFile, readErr := ioutil.ReadFile("input.txt")
	if readErr != nil {
		log.Fatal("Error on input", readErr)
	}
	inputsRegex := regexp.MustCompile("Wordle \\d{3} \\d/\\d\n\n((\U0001f7e9|\U0001f7e8|\u2b1b|\u2b1c){5}\n)+")
	inputs := inputsRegex.FindAllString(string(inputFile), -1)
	sharedInputs := []SharedInput{}
	for _, v := range(inputs) {
		sa := SharedInput{combinations: []Combination{}}
		split := strings.Split(string(v), "\n")

		for _, l := range(split[2:len(split)-1]) {
			sa.combinations = append(sa.combinations, CombinationString(l).toCombination())
		}
		sharedInputs = append(sharedInputs, sa)
	}
	return sharedInputs
}

func fileExists (path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false

	} else {
		return false
	}
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
	if fileExists(PRECOMPUTED_FILE) {
		log.Println("Loading precomputed from file")
		precomputedFile, precomputedErr := os.Open(PRECOMPUTED_FILE)
		if precomputedErr != nil {
			log.Fatal("Error reading precomputed")
		}
		dec := gob.NewDecoder(precomputedFile)
		var precomputed Precomputed
		decodeErr := dec.Decode(&precomputed)
		if decodeErr != nil {
			log.Fatal("Decode error", decodeErr)
		}
		return precomputed
	}
	file, openErr := ioutil.ReadFile("words.en.2.txt")
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
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(precomputed)
	ioutil.WriteFile(PRECOMPUTED_FILE, buf.Bytes(), 0644)
	return precomputed
}

type Precomputed struct {
	Dictionary []Word
	ComparisonAggregate []ComparisonAggregate
}



func computeCombination (input, solution Word) CombinationArray {
	count := map[byte]int{}
	if len(input) != WORD_LENGTH || len(solution) != WORD_LENGTH {
		log.Fatal("Unexpected length for string", input, len(input))
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

