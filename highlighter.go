package main

import (
	"sort"
)

type Fragment struct {
	Id    int
	Text  string
	Score float32
}

type FragmentIndex struct {
	docId   int
	NumDocs int
	index   map[string][]Posting

	// store page content for future use
	fragmentStore []Fragment

	// store field length in number of tokens
	fieldLen []int

	// avarageFieldLen
	avgFieldLen float64

	// Analyzer to use for text analysis and tokenization
	analyzer Analyzer
}

func NewFragmentIndex(analyzer Analyzer) *FragmentIndex {
	idx := &FragmentIndex{}
	idx.docId = 0

	idx.index = make(map[string][]Posting)

	// store page content for future use
	idx.fragmentStore = make([]Fragment, 0)
	idx.fieldLen = make([]int, 0)

	idx.analyzer = analyzer
	return idx
}

func (idx *FragmentIndex) Add(doc *Fragment) {
	// Set the docId of the document. It will be used as primary key for almost everyhing
	doc.Id = idx.docId

	// Start the Analysis process
	tokens := idx.analyzer.Analyze(doc.Text)

	for key, val := range tokenPositions(tokens) {
		//fmt.Println(key, val)
		posting := Posting{uint32(doc.Id), uint32(len(val)), 1.0, val}
		idx.index[key] = append(idx.index[key], posting)
	}

	idx.fragmentStore = append(idx.fragmentStore, *doc)
	idx.fieldLen = append(idx.fieldLen, len(tokens))

	// increment docId after ever document
	idx.docId++
	// increment total number of documents in index
	idx.NumDocs++
}

func (idx *FragmentIndex) updateAvgFieldLen() {
	total := 0

	for _, v := range idx.fieldLen {
		total += v
	}

	idx.avgFieldLen = float64(total) / float64(idx.NumDocs)
}

func (idx *FragmentIndex) Score(q string) []Fragment {
	tokens := idx.analyzer.Analyze(q)

	result := make([]Posting, 0)

	for i, token := range tokens {
		if i == 0 {
			result = idx.index[token.value]
			idx.scorePosting(result)
		} else {
			temp := idx.index[token.value]
			idx.scorePosting(temp)
			result = Union(temp, result)
		}
	}

	sort.Sort(ByBoost(result))

	fragments := make([]Fragment, 0)
	for _, v := range result {
		idx.fragmentStore[v.docId].Score = v.boost
		fragments = append(fragments, idx.fragmentStore[v.docId])
	}

	return fragments
}

func (idx *FragmentIndex) _scorePosting(postings []Posting) {
	for _, v := range postings {
		v.boost = float32(idf(float64(len(postings)), float64(idx.NumDocs)) * tf(float64(v.frequency), float64(idx.fieldLen[v.docId]), idx.avgFieldLen))
	}
}

func (idx *FragmentIndex) scorePosting(postings []Posting) {
	for i := range postings {
		postings[i].boost = float32(idf(float64(len(postings)), float64(idx.NumDocs)) * tf(float64(postings[i].frequency), float64(idx.fieldLen[postings[i].docId]), idx.avgFieldLen))
		//fmt.Println(postings[i].boost)
	}
}

type HighlightSegment struct {
	Score    int
	Fragment string
}

type ByScore []HighlightSegment

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

func HighlightText(text string, q string) string {

	analyzer := NewSimpleAnalyzer(NewSimpleTokenizer())

	turkishFilter := NewTurkishLowercaseFilter()
	turkishAccentFilter := NewTurkishAccentFilter()

	analyzer.AddTokenFilter(turkishFilter)
	analyzer.AddTokenFilter(turkishAccentFilter)

	index := NewFragmentIndex(analyzer)

	sentences := sentenceTokenizer.Tokenize(text)
	for _, s := range sentences {
		fragment := Fragment{0, s.Text, 0}
		index.Add(&fragment)
	}

	index.updateAvgFieldLen()
	fragments := index.Score(q)

	hlText := ""

	if len(fragments) > 2 {
		fragments = fragments[0:2]
	}

	for _, v := range fragments {
		hlText += v.Text + "<br>---------------------------------<br>"
	}

	return hlText
}
