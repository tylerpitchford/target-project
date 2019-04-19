package indexers

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"unicode"
)

type SingleTokenIndexer struct {
	GenericIndexer
	index map[string]int
}

func (i *SingleTokenIndexer) isSurroundedByNumbers(index int, runes []rune) bool {
	if len(runes)-1 > index+1 {
		return unicode.IsNumber(runes[index-1]) && unicode.IsNumber(runes[index+1])
	}

	return false
}

func (i *SingleTokenIndexer) Tokenize(str string) []string {
	return []string{str}
}

//simple tokenizer tuned to the provided corpus
//does not address all cases in the English language
//nor does it provide any consideration for foreign languages
func (i *SingleTokenIndexer) tokenize(byteSlice []byte) []string {
	var tokens []string

	runes := bytes.Runes(byteSlice)

	var buffer bytes.Buffer
	inQuotes := false

	for index, rune := range runes {
		//general case, character, number, or special character
		if unicode.IsLetter(rune) || unicode.IsNumber(rune) || i.isSpecialCorpusRune(rune) {
			buffer.WriteRune(rune)
		} else if rune == '"' { //handle quoted tokens
			inQuotes = !inQuotes
			if !inQuotes { //when closing out quotes handling trailing punctuation.
				tokens = i.handleForTrailingPuncuationInAQuote(runes[index-1], tokens, &buffer)
				tokens = i.closeOutTokenPrePost(tokens, &buffer, `"`, `"`)
			}
		} else if inQuotes { //skip all other processing if in a quoted block
			buffer.WriteRune(rune)
		} else if rune == ',' { // , requires some special processing
			if i.isSurroundedByNumbers(index, runes) {
				buffer.WriteRune(rune)
			} else {
				tokens = i.closeOutToken(tokens, &buffer)
				tokens = append(tokens, string(rune))
			}
		} else if unicode.IsPunct(rune) { //punctuation
			tokens = i.closeOutToken(tokens, &buffer)
			tokens = append(tokens, string(rune))
		} else if unicode.IsSpace(rune) { //whitespace
			tokens = i.closeOutToken(tokens, &buffer)
		}
	}

	//handle the scenario where a quote isn't ended
	if buffer.Len() > 0 {
		tokens = i.closeOutToken(tokens, &buffer)
	}

	return tokens
}

func (i *SingleTokenIndexer) BuildIndex() {
	bytes, err := ioutil.ReadFile(i.path)

	if err != nil {
		log.Fatal(err)
	}

	tokens := i.tokenize(bytes)
	tokenIndex := make(map[string]int)

	//this is already sorted
	for _, token := range tokens {
		tokenIndex[token] = tokenIndex[token]+1
	}

	i.index = tokenIndex
}

func (i *SingleTokenIndexer) Search(token []string) (count int){
	return i.index[token[0]]
}

func (i *SingleTokenIndexer) SerializeIndex() {
	buffer := new(bytes.Buffer)

	encoder := gob.NewEncoder(buffer)

	// Encoding the map
	err := encoder.Encode(i.index)
	if err != nil {
		panic(err)
	}

	//file, err := os.Create()
	err = ioutil.WriteFile(i.GetIdxFilename(), buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *SingleTokenIndexer) DeserializeIndex() {
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

func (i *SingleTokenIndexer) PrintIndex() {
	for key, value := range i.index {
		log.Println(key, ":", value)
	}
}