package indexers

import (
	"bytes"
	"strings"
	"unicode"
)

type GenericIndexer struct {
	path string
	count int
}

func (i *GenericIndexer) SetPath(path string) {
	i.path = path
}

func (i *GenericIndexer) closeOutToken(tokens []string, buffer *bytes.Buffer) []string {
	//close out the previous buffer
	return i.closeOutTokenPrePost(tokens, buffer, "","")
}

func (i *GenericIndexer) closeOutTokenPrePost(tokens []string, buffer *bytes.Buffer, pre string, post string) []string {
	//close out the previous buffer
	if buffer.Len() > 0 {
		tokens = append(tokens, pre + buffer.String() + post)
		buffer.Reset()
	}
	return tokens
}

func (i *GenericIndexer) handleForTrailingPuncuationInAQuote(prevRune rune, tokens []string, buffer *bytes.Buffer) []string {
	if unicode.IsPunct(prevRune) {
		// we're in a quoted block with trailing punc, remove it, split it into a new token
		buffer.Truncate(buffer.Len()-1)
		tokens = i.closeOutToken(tokens, buffer)
		tokens = append(tokens, ",")
	}
	return tokens
}

func (i *GenericIndexer) isSpecialCorpusRune(r rune) bool {
	return r == '-' || r == '\'' || r == '’'
}

func (i *GenericIndexer) GetIdxFilename() string {
	return strings.Replace(i.path, ".txt", ".idx", -1)
}