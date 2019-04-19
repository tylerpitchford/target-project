package test

//importing here causes the default usage output to include the 'test' parameters unfortunately

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"target-project/indexers"
	"target-project/search"
	"time"
)

const LOOP_COUNT = 2000000

func GoBenchmarkTextSearchNonConcurrent() {
	ExecuteSearch(1, false, false, "data", "./test/term.list")
}

func GoBenchmarkTextSearchConcurrent() {
	ExecuteSearch(1, false, true, "data", "./test/term.list")
}

func GoBenchmarkRegExSearchNonConcurrent() {
	ExecuteSearch(2, false, false, "data", "./test/term.list")
}

func GoBenchmarkRegExSearchConcurrent() {
	ExecuteSearch(2, false, true, "data", "./test/term.list")
}

func GoBenchmarkIndexSearchNonConcurrent() {
	ExecuteSearch(3, false, false, "data", "./test/term.list")
}

func GoBenchmarkIndexSearchConcurrent() {
	ExecuteSearch(3, false, true, "data", "./test/term.list")
}

func GoBenchmarkPositionalIndexSearchNonConcurrent() {
	ExecuteSearch(3, true, false, "data", "./test/term.list")
}

func GoBenchmarkPositionalIndexSearchConcurrent() {
	ExecuteSearch(3, true, true, "data", "./test/term.list")
}

func ExecuteSearch(searchType int, usePositional bool, concurrent bool, dataPath string, termFile string) {

	files := search.LoadFiles(dataPath)
	tokens := LoadRandomSearchTerms(termFile)

	if searchType == 3 {
		indexers.BuildIndicies(dataPath, usePositional)
		search.LoadIndices(files,usePositional)
	}

	for n := 0; n < LOOP_COUNT; n++ {
		//ignore regex parsing errors (which shouldn't happen anyway)
		searchParams, _ := search.NewSearchParameters(tokens[n%len(tokens)],
			searchType,
			files,
			usePositional,
			false)

		searchParams.Search(concurrent)
	}
}

func LoadRandomSearchTerms(path string) []string {
	var lines []string

	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })

	return lines
}

type benchFunction func()

func RunBenchmark(fn benchFunction) {
	startTime := time.Now()
	fn()
	log.Println("Total time: ", time.Now().Sub(startTime))
}

func RunBenchmarks() {
	log.Println("Starting benchmarks.")
	log.Println()
	log.Println("Text Search (non-concurrent)")
	RunBenchmark(GoBenchmarkTextSearchNonConcurrent)
	log.Println("Text Search (concurrent)")
	RunBenchmark(GoBenchmarkTextSearchConcurrent)
	log.Println("RegEx Search (non-concurrent)")
	RunBenchmark(GoBenchmarkRegExSearchNonConcurrent)
	log.Println("RegEx Search (concurrent)")
	RunBenchmark(GoBenchmarkRegExSearchConcurrent)
	log.Println("Index Search (single-token, non-concurrent)")
	RunBenchmark(GoBenchmarkIndexSearchNonConcurrent)
	log.Println("Index Search (single-token, concurrent)")
	RunBenchmark(GoBenchmarkIndexSearchConcurrent)
	log.Println("Index Search (positional, non-concurrent)")
	RunBenchmark(GoBenchmarkPositionalIndexSearchNonConcurrent)
	log.Println("Index Search (positional, concurrent)")
	RunBenchmark(GoBenchmarkPositionalIndexSearchConcurrent)
	log.Println()
	log.Println("Benchmarks complete.")
}


