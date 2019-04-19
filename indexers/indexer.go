package indexers

type Indexer interface {
	SetPath(string)
	BuildIndex()
	SerializeIndex()
	DeserializeIndex()
	PrintIndex()
	Tokenize(string) []string
	Search([]string) int
}
