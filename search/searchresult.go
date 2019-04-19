package search

type SearchResult struct {
	Filename string
	Count int
}

type ResultSorter []SearchResult

func (r ResultSorter) Len() int           { return len(r) }
func (r ResultSorter) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResultSorter) Less(i, j int) bool {
	if r[i].Count != r[j].Count {
		return r[i].Count > r[j].Count
	} else {
		return r[i].Filename > r[j].Filename
	}
}
