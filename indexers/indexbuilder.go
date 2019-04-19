package indexers

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func BuildIndicies(path string, positional bool) {
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
		var indexer Indexer
		if positional {
			indexer = &PositionalIndexer{}
		} else {
			indexer = &SingleTokenIndexer{}
		}
		indexer.SetPath(path)
		indexer.BuildIndex()
		indexer.SerializeIndex()
	}
}
