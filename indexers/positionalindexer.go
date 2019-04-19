package indexers

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"strings"
	"sync"
	"unicode"
	"log"
)

var Empty struct{}

type PositionalIndexer struct {
	GenericIndexer
	index map[string]map[int]struct{}
}

func (i *PositionalIndexer) Tokenize(str string) []string {
	return i.tokenize([]byte(str))
}

func (i *PositionalIndexer) tokenize(byteSlice []byte) []string {
	var tokens []string

	runes := bytes.Runes(byteSlice)

	var buffer bytes.Buffer
	//inQuotes := false

	for _, rune := range runes {

		//general case, character, number, or special character
		if unicode.IsLetter(rune) || unicode.IsNumber(rune) || i.isSpecialCorpusRune(rune) {
			buffer.WriteRune(rune)
		} else if unicode.IsPunct(rune) { //punctuation
			tokens = i.closeOutToken(tokens, &buffer)
			tokens = append(tokens, string(rune))
		} else if unicode.IsSpace(rune) { //whitespace
			tokens = i.closeOutToken(tokens, &buffer)
		}
	}

	if buffer.Len() > 0 {
		tokens = i.closeOutToken(tokens, &buffer)
	}

	return tokens
}

func (i *PositionalIndexer)  BuildIndex() {
	bytes, err := ioutil.ReadFile(i.path)

	if err != nil {
		log.Fatal(err)
	}

	tokens := i.tokenize(bytes)
	tokenIndex := make(map[string]map[int]struct{})

	//this is already sorted
	for i, token := range tokens {
		if _, ok := tokenIndex[token]; !ok {
			tokenIndex[token] = make(map[int]struct{})
		}
		tokenIndex[token][i] = Empty
	}

	i.index = tokenIndex
}

func (i *PositionalIndexer)  GetIdxFilename() string {
	return strings.Replace(i.path, ".txt", ".idx", -1)
}

func (i *PositionalIndexer)  SerializeIndex() {
	buffer := new(bytes.Buffer)

	encoder := gob.NewEncoder(buffer)

	// Encoding the map
	err := encoder.Encode(i.index)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(i.GetIdxFilename(), buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *PositionalIndexer) DeserializeIndex() {
	byteData, err := ioutil.ReadFile(i.GetIdxFilename())
	if err != nil {
		log.Fatal(err)
	}

	decoder := gob.NewDecoder(bytes.NewReader(byteData))
	err = decoder.Decode(&i.index)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *PositionalIndexer) checkForToken(index int, tokens []string, results chan bool, wait *sync.WaitGroup) {
	defer wait.Done()
	results <- i.checkForNextToken(index+1, 1, tokens)
}

func (i *PositionalIndexer) checkForNextToken(index int, currentToken int, tokens []string, ) bool {

	if currentToken == len(tokens) {
		return true
	}

	token := tokens[currentToken]
	indices := i.index[token]

	if _, ok := indices[index]; ok {
		return i.checkForNextToken(index+1, currentToken+1, tokens)
	}
	return false
}

func (i *PositionalIndexer) Search(tokens []string) (count int) {
	results := make(chan bool)
	var wait sync.WaitGroup

	if positionalValues, ok := i.index[tokens[0]]; ok {

		wait.Add(1)
		//spawn off a thread to process results
		go i.processCountResults(results, len(i.index[tokens[0]]), &wait)

		for key, _ := range positionalValues {
			wait.Add(1)
			go i.checkForToken(key, tokens, results, &wait)
		}
	}

	//let all the results come back
	wait.Wait()

	return i.count
}

func (i *PositionalIndexer) processCountResults(results chan bool, resultNumber int, wait *sync.WaitGroup) {
	defer wait.Done()

	count := 0
	//we have results, sort and print
	for result := range results {
		if result {
			count++
		}

		resultNumber--

		if resultNumber == 0 {
			close(results)
		}
	}

	i.count = count
}

func (i *PositionalIndexer) PrintIndex() {
	for key, value := range i.index {
		log.Println(key, ":", value)
	}
}