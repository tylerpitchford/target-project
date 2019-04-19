package search

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"target-project/indexers"
	"log"
)

type SearchableFile struct {
	Path string
	StringData string
	SearchIndexer indexers.Indexer
}

func LoadFiles(path string) (results []*SearchableFile) {
	var paths []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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

	for _, path := range paths {

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &SearchableFile{path, string(bytes), nil})
	}

	return results
}


func LoadIndices(files []*SearchableFile, positional bool) {
	for _, file := range files {
		if positional {
			file.SearchIndexer = &indexers.PositionalIndexer{}
		} else {
			file.SearchIndexer = &indexers.SingleTokenIndexer{}
		}
		file.SearchIndexer.SetPath(file.Path)
		file.SearchIndexer.DeserializeIndex()
	}
}