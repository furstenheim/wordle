package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"bytes"
	"encoding/gob"
)

const WORD_LENGTH = 5

// Uncomment for wordle in Spanish
// const PRECOMPUTED_FILE = "encoded-precomputed-es"
// const WORDS_FILE = "words.es.txt"
// const POSSIBLE_WORDS = "words.es.txt"
// const INPUT_FILE = "input.es.txt"
// const REGEX = "Wordle \\(ES\\) #\\d{2} (\\d|X)/\\d\n\n((\U0001f7e9|\U0001f7e8|\u2b1b|\u2b1c){5}\n)+"
const PRECOMPUTED_FILE = "encoded-precomputed"
const WORDS_FILE = "words.en.2.txt"
const POSSIBLE_WORDS = "possible_answers.en.txt"
const INPUT_FILE = "input.txt"
const REGEX = "Wordle \\d{3} \\d/\\d\n\n((\U0001f7e9|\U0001f7e8|\u2b1b|\u2b1c){5}\n)+"
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
	sharedInputs := parseInput()
	// log.Println("computed precomputed")
	playedHints := []FullHint{}
	reader := bufio.NewReader(os.Stdin)
	for true {
		possibleWordsIndexes := getPossibleListOfWords(playedHints, precomputed, sharedInputs)

		if len(possibleWordsIndexes) == 1 {
			fmt.Printf("That was fast!! The solution is %s\n", precomputed.DictionaryPossible[possibleWordsIndexes[0]])
			return
		}

		min := len(possibleWordsIndexes)
		var maxWord Word

		for _, v := range(precomputed.Dictionary) {
			ca := ComparisonAggregate{}
			for _, i := range possibleWordsIndexes {
				w := precomputed.DictionaryPossible[i]
				combination := computeCombination(v, w).toNumber()
				ca[combination]++
			}
			minOfCombinations := 0
			for _, count := range ca {
				minOfCombinations = max(count, minOfCombinations)
			}
			if minOfCombinations < min {
				min = minOfCombinations
				maxWord = v
			}
		}
		fmt.Printf("There are currently %d possible words. Best move is to try \033[1m%s\033[0m\n", len(possibleWordsIndexes), maxWord)
		fmt.Printf("Enter the response to the clue:\n - \033[1mW\033[0m: White \n - \033[1mY\033[0m: Yellow \n - \033[1mG\033[0m: Green \n")

		input, _ := reader.ReadString('\n')
		split := strings.Split(input, "")
		playedHint := FullHint{}
		for i, v := range (split[0:5]) {
			switch v {
			case "W":
				playedHint.combinationArray[i] = White
			case "Y":
				playedHint.combinationArray[i] = Yellow
			default:
				playedHint.combinationArray[i] = Green
			}
		}
		playedHint.word = maxWord
		playedHints = append(playedHints, playedHint)
	}
}

func minOf (a, b int) int {
	if a < b {
		return a
	}
	return b
}


type FullHint struct {
	combinationArray CombinationArray
	word Word
}
func getPossibleListOfWords(playedHints []FullHint, precomputed Precomputed, sharedInputs []SharedInput) []int {
	caInput := sharedInputs[0].toComparisonAggregate()

	for _, v := range sharedInputs[1:] {
		caInput = mergeComparisonAggregate(caInput, v.toComparisonAggregate())
	}
	possibleWords := []int{}

	possibleWordsLabel: for i, ca := range(precomputed.ComparisonAggregate) {
		for _, hint := range playedHints {
			combination := computeCombination(hint.word, precomputed.DictionaryPossible[i])
			if combination != hint.combinationArray {
				continue possibleWordsLabel
			}
		}
		if caInput.isCompatibleWith(ca) {
			possibleWords = append(possibleWords, i)
		}
	}
	return possibleWords
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
	inputFile, readErr := ioutil.ReadFile(INPUT_FILE)
	if readErr != nil {
		log.Fatal("Error on input", readErr)
	}
	inputsRegex := regexp.MustCompile(REGEX)
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
			ca[i] = White
		} else if v == 11036 {
			ca[i] = White
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
		// log.Println("Loading precomputed from file")
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
	file, openErr := ioutil.ReadFile(WORDS_FILE)
	if openErr != nil {
		log.Fatal("Error opening", openErr)
	}

	filePossible, openPossibleErr := ioutil.ReadFile(POSSIBLE_WORDS)
	if openPossibleErr != nil {
		log.Fatal("Error opening possible", openPossibleErr)
	}
	split := strings.Split(string(file), "\n")
	splitPossible := strings.Split(string(filePossible), "\n")
	dictionary := make([]Word, len(split))
	dictionaryPossible := make([]Word, len(splitPossible))
	for i, v := range(split) {
		dictionary[i] = Word(v)
	}
	for i, v := range(splitPossible) {
		dictionaryPossible[i] = Word(v)
	}
	precomputed := Precomputed{
		DictionaryPossible: dictionaryPossible,
		Dictionary: dictionary,
		ComparisonAggregate: make([]ComparisonAggregate, len(dictionaryPossible)),
	}
	for i, v := range (dictionaryPossible) {
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
	DictionaryPossible []Word
	ComparisonAggregate []ComparisonAggregate
}



func computeCombination (input, solution Word) CombinationArray {
	count := map[rune]int{}
	inputArray := []rune(input)
	solutionArray := []rune(solution)
	if len(inputArray) != WORD_LENGTH || len(solutionArray) != WORD_LENGTH {
		log.Fatal("Unexpected length for string", inputArray, len(inputArray), solutionArray, len(solutionArray), fmt.Sprintf(" '%x' ", inputArray))
	}
	res := CombinationArray{}

	for i := 0; i < WORD_LENGTH; i++ {
		if inputArray[i] == solutionArray[i] {
			res[i] = Green
		} else {
			count[solutionArray[i]]++
		}
	}
	for i := 0; i < WORD_LENGTH; i++ {
		if inputArray[i] != solutionArray[i] && count[inputArray[i]] > 0 {
			count[inputArray[i]]--
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
	White CombinationColor = iota
	Yellow
	Green
)

