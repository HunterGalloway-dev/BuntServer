package anagrams

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type sortByte []byte

func (s sortByte) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortByte) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByte) Len() int {
	return len(s)
}

func SortString(s string) string {
	b := []byte(s)
	sort.Sort(sortByte(b))

	return string(b)
}

const (
	minWordLength = 3
	maxWordLength = 6
)

type AnagramSolver struct {
	wordMap map[string][]string
}

func BuildAnagramSolver(wordFileDir string) *AnagramSolver {
	anagramSolver := new(AnagramSolver)
	anagramSolver.wordMap = make(map[string][]string)

	files, err := os.ReadDir(wordFileDir)
	if err != nil {
		log.Fatal(err)
	}

	wordCount := 0
	for _, file := range files {
		if !file.IsDir() {
			filePath := fmt.Sprintf("./%v/%v", wordFileDir, file.Name())

			file, err := os.Open(filePath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				word := scanner.Text()
				sortedWord := SortString(word)

				if length := len(word); length <= maxWordLength && length >= minWordLength {
					_, exists := anagramSolver.wordMap[sortedWord]
					wordCount++

					if exists {
						anagramSolver.wordMap[sortedWord] = append(anagramSolver.wordMap[sortedWord], word)
					} else {
						anagramSolver.wordMap[sortedWord] = []string{word}
					}
				}
			}
		}
	}
	log.Printf("Processed %v words\n", wordCount)

	return anagramSolver
}

func (a *AnagramSolver) GetWords(input string) []string {
	var words []string

	possCombos := genCombos(input)

	for _, combo := range possCombos {
		if _, ok := a.wordMap[combo]; ok {
			words = append(words, a.wordMap[combo]...)
		}
	}

	return words
}

func genCombo(input []string, goal int, curCombo string) []string {
	var combos []string

	if goal == 0 {
		return []string{curCombo}
	}
	if len(input) != 0 {
		for i, val := range input {
			combos = append(combos, genCombo(input[i+1:], goal-1, curCombo+val)...)
		}
	}

	return combos
}

func genCombos(input string) []string {
	inputSlice := strings.Split(SortString(input), "")
	var combos []string

	for goal := maxWordLength; goal >= minWordLength; goal-- {
		combos = append(combos, genCombo(inputSlice, goal, "")...)
	}

	return combos
}
