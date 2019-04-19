package search

import (
	"fmt"
	"sort"
	"target-project/indexers"
	"testing"
)

const DATA_DIR = "../data"

func generateSearchResultSlice(frenchCount int, hitchikerCount int, warpCount int) []SearchResult {
	searchResults := []SearchResult{
		{"french_armed_forces.txt", frenchCount},
		{"hitchhikers.txt", hitchikerCount},
		{"warp_drive.txt", warpCount},
	}

	sort.Sort(ResultSorter(searchResults))

	return searchResults
}

type TestSearchResult struct {
	searchToken string
	searchType int
	useConcurrent bool
	usePositional bool
	dataPath string
	results []SearchResult
	expectError bool
}

func TestSearch(t *testing.T) {
	searchTests := []TestSearchResult{
		//The
		//non-current
		{"The", 1, false, false, DATA_DIR, generateSearchResultSlice(7,9,0), false},
		{"The", 2, false, false, DATA_DIR, generateSearchResultSlice(7,9,0), false},
		{"The", 3, false, false, DATA_DIR, generateSearchResultSlice(7,6,0), false},
		{"The", 3, false, true, DATA_DIR, generateSearchResultSlice(7,8,0), false},

		//concurrent
		{"The", 1, true, false, DATA_DIR, generateSearchResultSlice(7,9,0), false},
		{"The", 2, true, false, DATA_DIR, generateSearchResultSlice(7,9,0), false},
		{"The", 3, true, false, DATA_DIR, generateSearchResultSlice(7,6,0), false},
		{"The", 3, true, true, DATA_DIR, generateSearchResultSlice(7,8,0), false},

		//Bir Hakeim (1942).
		//non-current
		{"Bir Hakeim (1942).", 1, false, false, DATA_DIR, generateSearchResultSlice(1,0,0), false},
		{"Bir Hakeim (1942).", 2, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"Bir Hakeim (1942).", 3, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"Bir Hakeim (1942).", 3, false, true, DATA_DIR, generateSearchResultSlice(1,0,0), false},

		//concurrent
		{"Bir Hakeim (1942).", 1, true, false, DATA_DIR, generateSearchResultSlice(1,0,0), false},
		{"Bir Hakeim (1942).", 2, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"Bir Hakeim (1942).", 3, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"Bir Hakeim (1942).", 3, true, true, DATA_DIR, generateSearchResultSlice(1,0,0), false},

		//unicode - film’s
		//non-current
		{"film’s", 1, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"film’s", 2, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"film’s", 3, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"film’s", 3, false, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		//concurrent
		{"film’s", 1, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"film’s", 2, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"film’s", 3, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"film’s", 3, true, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		//online guide).
		//good text, malformed regex, wrong token, good positional
		//non-concurrent
		{"online guide).", 1, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"online guide).", 2, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), true},
		{"online guide).", 3, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"online guide).", 3, false, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		//concurrent
		{"online guide).", 1, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"online guide).", 2, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), true},
		{"online guide).", 3, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"online guide).", 3, true, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		//[The] Guide
		//good text, bad regex, wrong token, good positional
		{"[The] Guide", 1, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"[The] Guide", 2, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"[The] Guide", 3, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"[The] Guide", 3, false, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		{"[The] Guide", 1, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"[The] Guide", 2, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"[The] Guide", 3, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"[The] Guide", 3, true, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		//"[The] Guide"
		//good text, bad regex, good token, good positional
		{"\"[The] Guide\"", 1, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"\"[The] Guide\"", 2, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"\"[The] Guide\"", 3, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"\"[The] Guide\"", 3, false, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		{"\"[The] Guide\"", 1, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"\"[The] Guide\"", 2, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{"\"[The] Guide\"", 3, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{"\"[The] Guide\"", 3, true, true, DATA_DIR, generateSearchResultSlice(0,1,0), false},

		//\"\[The\] Guide\"
		//bad text, good regex, bad token, bad positional
		{`\"\[The\] Guide\"`, 1, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{`\"\[The\] Guide\"`, 2, false, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{`\"\[The\] Guide\"`, 3, false, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{`\"\[The\] Guide\"`, 3, false, true, DATA_DIR, generateSearchResultSlice(0,0,0), false},

		{`\"\[The\] Guide\"`, 1, true, false, DATA_DIR, generateSearchResultSlice(0,0,0), false},
		{`\"\[The\] Guide\"`, 2, true, false, DATA_DIR, generateSearchResultSlice(0,1,0), false},
		{`\"\[The\] Guide\"`, 3, true, false, DATA_DIR, generateSearchResultSlice(0,0, 0), false},
		{`\"\[The\] Guide\"`, 3, true, true, DATA_DIR, generateSearchResultSlice(0,0,0), false},

		//.
		{".", 1, false, false, DATA_DIR, generateSearchResultSlice(26,10,5), false},
		{".", 2, false, false, DATA_DIR, generateSearchResultSlice(4068,1867,1049), false},
		{".", 3, false, false, DATA_DIR, generateSearchResultSlice(26,10,5), false},
		{".", 3, false, true, DATA_DIR, generateSearchResultSlice(26,10,5), false},

		{".", 1, true, false, DATA_DIR, generateSearchResultSlice(26,10,5), false},
		{".", 2, true, false, DATA_DIR, generateSearchResultSlice(4068,1867,1049), false},
		{".", 3, true, false, DATA_DIR, generateSearchResultSlice(26,10,5), false},
		{".", 3, true, true, DATA_DIR, generateSearchResultSlice(26,10,5), false},

	}

	for _, test := range searchTests {
		t.Run(fmt.Sprintf("%s %d %t %t %t ", test.searchToken, test.searchType, test.useConcurrent, test.usePositional, test.expectError), func(t *testing.T) {
			results, err := executeSearch(test.searchToken, test.searchType, test.useConcurrent, test.usePositional, test.dataPath)
			if err != nil {
				if test.expectError {
					return
				}
				t.Error("Unexpected error: ", err)
			} else if !Equal(test.results, results) {
				t.Error("Results don't match.")
			}
		})
	}
}

func Equal(a, b []SearchResult) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func executeSearch(searchToken string, searchType int, concurrent bool, usePositional bool, dataPath string) ([]SearchResult, error) {

	files := LoadFiles(dataPath)

	if searchType == 3 {
		indexers.BuildIndicies(dataPath, usePositional)
		LoadIndices(files,usePositional)
	}

	searchParams, err := NewSearchParameters(searchToken,
		searchType,
		files,
		usePositional,
		false)

	if err != nil {
		return nil, err
	}

	return searchParams.Search(concurrent), err

}