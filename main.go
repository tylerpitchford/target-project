package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"target-project/indexers"
	"target-project/search"
	"target-project/test"
)

const SEARCH_TERM_PROMPT = "Enter the search term: "
const SEARCH_METHOD_PROMPT = "Search Method: 1) String Match 2) Regular Expression 3) Indexed: "
const SEARCH_METHOD_ERROR = "You must supply a search type of either: 1, 2, or 3. Please try again."

type RuntimeFlags struct {
	PositionalIndex bool
	RunBenchmarks bool
	DataDirectory *os.File
	RunConcurrent bool
	SearchToken string
	SearchType int
}

func ReadString(prompt string) (string) {
	input, err := ReadUserInput(prompt)

	//loop until we get a good entry
	for err != nil {
		input, err = ReadUserInput(prompt)
	}

	return input
}

func ReadInteger(prompt string) int {
	input := ReadString(prompt)

	var parsedInt int
	var intErr error
	parsedInt, intErr = ParseAndValidateInput(input)
	for ; intErr != nil; parsedInt, intErr = ParseAndValidateInput(input) {
		fmt.Println(intErr)
		if intErr != nil {
			fmt.Println(intErr)
		}
		input = ReadString(prompt)
	}

	return parsedInt
}

func CheckSearchTypeBounds(searchType int) error {
	if searchType < 1 || searchType > 3 {
		return errors.New(SEARCH_METHOD_ERROR)
	}
	return nil
}

func ParseAndValidateInput(input string) (int, error) {
	parsedInt, err := strconv.Atoi(input)

	if err != nil {
		return -1, errors.New(SEARCH_METHOD_ERROR)
	}

	err = CheckSearchTypeBounds(parsedInt)
	if err != nil {
		return -1, errors.New(SEARCH_METHOD_ERROR)
	}

	return parsedInt, nil
}

func ReadUserInput(prompt string) (string, error) {

	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := scanner.Text()

	return input, nil
}

func (r *RuntimeFlags) parse() {
	flag.BoolVar(&r.PositionalIndex,"positional", false, "Use a positional search indicies.")
	flag.BoolVar(&r.RunBenchmarks,"benchmark", false, "Run the benchmarks.")
	flag.BoolVar(&r.RunConcurrent,"concurrent", false, "Run the search concurrently.")
	flag.StringVar(&r.SearchToken,"token", "", "Provide the search token non-interactively.")
	flag.IntVar(&r.SearchType,"type", -1, "Provide the search type non-interactively.")

	dir := *flag.String("directory", "data", "Provide a directory where files should be searched or indexed. Only files with the extension .txt are considered.")

	flag.Parse()

	file, err := os.Open(dir)
	if err != nil {
		log.Fatal("The data directory wasn't found. Please try again.")
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if !fileInfo.IsDir() {
		log.Fatal("The directory flag must point to a directory. Please try again.")
	}

	if r.SearchToken != "" && r.SearchType != -1 {
		err = CheckSearchTypeBounds(r.SearchType)
		if err != nil {
			log.Fatal(SEARCH_METHOD_ERROR)
		}
	}

	r.DataDirectory = file
}

func interactiveSearch(runtime RuntimeFlags) {

	//should we build a positional index or a single-token?
	if runtime.PositionalIndex {
		indexers.BuildIndicies("data", true)
	} else {
		indexers.BuildIndicies("data", false)
	}

	//load all search files
	files := search.LoadFiles(runtime.DataDirectory.Name())
	//load the indicies
	search.LoadIndices(files, runtime.PositionalIndex)

	var searchToken string
	var searchType int

	if runtime.SearchToken != "" && runtime.SearchType != -1 {
		searchToken = runtime.SearchToken
		searchType = runtime.SearchType
	} else {
		searchToken = ReadString(SEARCH_TERM_PROMPT)
		fmt.Println()
		searchType = ReadInteger(SEARCH_METHOD_PROMPT)
		fmt.Println()
	}

	searchParams, err := search.NewSearchParameters(searchToken, searchType, files, runtime.PositionalIndex, true)
	if err != nil {
		log.Fatal(err)
	}

	//execute the search
	searchParams.Search(runtime.RunConcurrent)
}

func main() {
	var runtime RuntimeFlags
	runtime.parse()

	//make this a flag
	if runtime.RunBenchmarks {
		test.RunBenchmarks()
	} else {
		interactiveSearch(runtime)
	}
}

func init() {
	//override the auotmatic printing because of the benchmark imports
	flag.Usage = func() {
		fmt.Println(`  -benchmark
    	Run the benchmarks.
  -concurrent
    	Run the search concurrently.
  -directory string
    	Provide a directory where files should be searched or indexed. Only files with the extension .txt are considered. (default "data")
  -positional
    	Use a positional search indicies.
  -token string
    	Provide the search token non-interactively.
  -type int
    	Provide the search type non-interactively. (default -1)`)
	}
}


