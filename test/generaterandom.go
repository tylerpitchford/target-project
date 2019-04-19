package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"log"
)

func GenerateTokenList() {
	randomfile, _ := os.Create("random.terms")
	defer randomfile.Close()

	var paths []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		//only process .txt files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range(paths) {

		bytes, _ := ioutil.ReadFile(path)

		str := string(bytes)

		tokens := strings.Fields(str)

		for _, token := range tokens {
			_, err := regexp.Compile(`(\s|\b)` + token + `(\s|\b)`)

			if err == nil {
				randomfile.WriteString(token + "\n")
			}
		}
	}
}
