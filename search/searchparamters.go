package search

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"target-project/indexers"
	"time"
)

type SearchParameters struct {
	SearchToken string
	SearchTokenRegex *regexp.Regexp
	SearchTokenIndex []string
	SearchType int
	SearchFiles []*SearchableFile
	UsePositionalIndex bool
	EnableOutput bool
}

type searchFunction func() int

func NewSearchParameters(token string, searchType int, files []*SearchableFile, usePositional bool, enableOutput bool) (SearchParameters, error) {

	s := SearchParameters{}

	s.SearchToken = token
	s.SearchType = searchType
	s.SearchFiles = files
	s.EnableOutput = enableOutput
	s.UsePositionalIndex = usePositional

	//precompile the regex
	if searchType == 2 {
		regex, err := regexp.Compile(s.SearchToken)
		if err != nil {
			return s, err
		}
		s.SearchTokenRegex = regex
	}

	if searchType == 3 {
		if s.UsePositionalIndex {
			indexer := &indexers.PositionalIndexer{}
			s.SearchTokenIndex = indexer.Tokenize(s.SearchToken)
		} else {
			indexer := &indexers.SingleTokenIndexer{}
			s.SearchTokenIndex = indexer.Tokenize(s.SearchToken)
		}
	}

	return s, nil
}

func (s *SearchParameters) Search(concurrent bool) []SearchResult {
	var searchResults []SearchResult
	currentTime := time.Now()

	switch s.SearchType {

	case 1:
		if concurrent {
			searchResults = s.StringMatchConcurrent()
		} else {
			searchResults = s.StringMatchNonConcurrent()
		}
	case 2:
		if concurrent {
			searchResults = s.RegexMatchConcurrent()
		} else {
			searchResults = s.RegexMatchNonConcurrent()
		}
	case 3:
		if concurrent {
			searchResults = s.IndexSearchConcurrent()
		} else {
			searchResults = s.IndexSearchNonConcurrent()
		}
	}

	//sort
	sort.Sort(ResultSorter(searchResults))

	//print
	if s.EnableOutput {
		for _, result := range searchResults  {
			fmt.Println("\t", result.Filename, "-", result.Count, "matches")
			fmt.Println()
		}
	}

	if s.EnableOutput {
		fmt.Println("Elapsed time:", time.Now().Sub(currentTime))
	}

	return searchResults
}

//NON-CONCURRENT SEARCHES
func (s *SearchParameters) StringMatchNonConcurrent() []SearchResult {
	var results []SearchResult

	for _, file := range s.SearchFiles {
		_, filename := filepath.Split(file.Path)
		result := SearchResult{filename, strings.Count(file.StringData, s.SearchToken)}
		results = append(results, result)
	}

	return results
}

func (s *SearchParameters) RegexMatchNonConcurrent() []SearchResult {
	var results []SearchResult

	for _, file := range s.SearchFiles {
		_, filename := filepath.Split(file.Path)
		result := SearchResult{filename, len(s.SearchTokenRegex.FindAllStringIndex(file.StringData, -1))}
		results = append(results, result)
	}

	return results
}

func (s *SearchParameters) IndexSearchNonConcurrent() []SearchResult {
	var results []SearchResult

	for _, file := range s.SearchFiles {
		_, filename := filepath.Split(file.Path)
		result := SearchResult{filename, file.SearchIndexer.Search(s.SearchTokenIndex)}
		results = append(results, result)
	}

	return results
}

//CONCURRENT SEARCHES
func (s *SearchParameters) CountInstances(file SearchableFile, results chan SearchResult, fn searchFunction) {
	_, filename := filepath.Split(file.Path)
	results <- SearchResult{filename, fn()}
}

func (s *SearchParameters) CountInstancesTextSearch(file SearchableFile, results chan SearchResult) {
	fn := func() int {
		return strings.Count(file.StringData, s.SearchToken)
	}
	s.CountInstances(file, results, fn)
}

func (s *SearchParameters) StringMatchConcurrent() []SearchResult {
	results := make(chan SearchResult)
	resultNumber := len(s.SearchFiles)

	for _, file := range s.SearchFiles {
		go s.CountInstancesTextSearch(*file, results)
	}

	var searchResults []SearchResult


	for result := range results {
		searchResults = append(searchResults, result)
		resultNumber--

		if resultNumber == 0 {
			close(results)
		}
	}

	return searchResults
}

func (s *SearchParameters) CountInstancesRegEx(file SearchableFile, regex *regexp.Regexp, results chan SearchResult) {
	fn := func() int {
		return len(regex.FindAllStringIndex(file.StringData, -1))
	}
	s.CountInstances(file, results, fn)
}

func (s *SearchParameters) RegexMatchConcurrent() []SearchResult {
	results := make(chan SearchResult)
	resultNumber := len(s.SearchFiles)

	for _, file := range s.SearchFiles {
		go s.CountInstancesRegEx(*file, s.SearchTokenRegex, results)
	}

	var searchResults []SearchResult
	for result := range results {
		searchResults = append(searchResults, result)
		resultNumber--

		if resultNumber == 0 {
			close(results)
		}
	}

	return searchResults
}

func (s *SearchParameters) CountInstancesIndexSearch(file SearchableFile, results chan SearchResult) {
	fn := func() int {
		return file.SearchIndexer.Search(s.SearchTokenIndex)
	}
	s.CountInstances(file, results, fn)
}

func (s *SearchParameters) IndexSearchConcurrent() []SearchResult {
	results := make(chan SearchResult)
	resultNumber := len(s.SearchFiles)

	for _, file := range s.SearchFiles {
		go s.CountInstancesIndexSearch(*file, results)
	}

	var searchResults []SearchResult
	for result := range results {
		searchResults = append(searchResults, result)
		resultNumber--

		if resultNumber == 0 {
			close(results)
		}
	}

	return searchResults
}