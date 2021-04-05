package main

import (
	"fmt"
	"sort"

	"github.com/RoaringBitmap/roaring"
)

type Book struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Type     string   `json:"type"`
	Genre    string   `json:"genre"`
	Category []string `json:"category"`
	Year     int      `json:"year"`
	NumPages int      `json:"num_pages"`
	Hash     string   `json:"hash"`
}

type BookIndex struct {
	docId   int
	NumDocs int
	index   map[string][]Posting

	// docCategories will store docIds for every document that blongs to a category
	bookCategory map[string][]uint32
	bookType     map[string][]uint32
	bookGenre    map[string][]uint32

	// roaring bitmaps to store bookCategory bitmaps
	categoryBitmaps map[string]*roaring.Bitmap

	// Store book metadata
	bookStore []Book

	// store field length in number of tokens
	fieldLen []int

	// avarage field length
	avgFieldLen float64

	// Analyzer to use for text analysis and tokenization
	analyzer Analyzer
}

func NewBookIndex(analyzer Analyzer) *BookIndex {
	idx := &BookIndex{}
	idx.docId = 0

	idx.index = make(map[string][]Posting)
	idx.bookCategory = make(map[string][]uint32)
	idx.bookType = make(map[string][]uint32)
	idx.bookGenre = make(map[string][]uint32)

	idx.categoryBitmaps = make(map[string]*roaring.Bitmap)

	idx.bookStore = make([]Book, 0)

	// store field length in number of tokens
	idx.fieldLen = make([]int, 0)

	idx.analyzer = analyzer
	return idx
}

func (idx *BookIndex) Add(doc *Book) {
	// Set the docId of the document. It will be used as primary key for almost everyhing
	doc.Id = idx.docId

	// Start the Analysis process
	//idx.analyze(doc)
	tokens := idx.analyzer.Analyze(doc.Title)

	// add document categories to index
	for _, category := range doc.Category {
		idx.bookCategory[category] = append(idx.bookCategory[category], uint32(doc.Id))
	}

	idx.bookType[doc.Type] = append(idx.bookType[doc.Type], uint32(doc.Id))
	idx.bookGenre[doc.Genre] = append(idx.bookGenre[doc.Genre], uint32(doc.Id))

	for key, val := range tokenPositions(tokens) {
		//fmt.Println(key, val)
		posting := Posting{uint32(doc.Id), uint32(len(val)), 1.0, val}
		idx.index[key] = append(idx.index[key], posting)
	}

	idx.bookStore = append(idx.bookStore, *doc)

	idx.fieldLen = append(idx.fieldLen, len(tokens))

	// increment docId after ever document
	idx.docId++
}

func (idx *BookIndex) GetBook(hash string) Book {
	for _, book := range idx.bookStore {
		if book.Hash == hash {
			return book
		}
	}
	return Book{}
}

func (idx *BookIndex) Search(q string) []Posting {
	var result []Posting
	var temp []Posting
	var resultPhrase []Posting

	tokens := idx.analyzer.Analyze(q)

	for i, token := range tokens {
		if i == 0 {
			result = make([]Posting, len(idx.index[token.value]))
			copy(result, idx.index[token.value])
			//fmt.Println(result)
			idx.scorePosting(result)
			//fmt.Println(result)
		} else {
			//temp := idx.index[token.value]
			temp = make([]Posting, len(idx.index[token.value]))
			copy(temp, idx.index[token.value])
			idx.scorePosting(temp)

			// boolean AND query
			result = Intersection(temp, result)
			// boolean OR query
			//result = Union(temp, result)
			// Phrase Query
			//result = PhraseQuery_FullMatch(result, temp)
		}
	}

	for i, token := range tokens {
		if i == 0 {
			resultPhrase = make([]Posting, len(idx.index[token.value]))
			copy(resultPhrase, idx.index[token.value])
			//fmt.Println(result)
			idx.scorePosting(result)
			//fmt.Println(result)
		} else {
			//temp := idx.index[token.value]
			temp = make([]Posting, len(idx.index[token.value]))
			copy(temp, idx.index[token.value])
			idx.scorePosting(temp)

			// boolean AND query
			// result = Intersection(temp, result)
			// boolean OR query
			//result = Union(temp, result)
			// Phrase Query
			resultPhrase = PhraseQuery_FullMatch(resultPhrase, temp)
		}
	}

	result = Union(result, resultPhrase)

	//fmt.Println(result)
	sort.Sort(ByBoost(result))
	//fmt.Println("-------------------------------------------------")
	//fmt.Println(result)

	return result
}

func (idx *BookIndex) buildCategoryBitmap() {

	for k, v := range idx.bookCategory {
		rb := roaring.NewBitmap()
		rb.AddMany(v)
		idx.categoryBitmaps[k] = rb
	}
}

func (idx *BookIndex) HighlightTitle(text string, q string) string {
	// do actual highlighting <b>term</b>
	//text := idx.fragmentStore[v.docId].Text

	tokens := idx.analyzer.Analyze(q)
	// just search for unique tokens, don't search for twice
	tokens = getUniqueTokens(tokens)

	textTokens := bookIndex.analyzer.Analyze(text)

	starts := make([]uint32, 0)
	ends := make([]uint32, 0)

	for _, token := range tokens {
		for _, tt := range textTokens {
			if token.value == tt.value {
				starts = append(starts, tt.start)
				ends = append(ends, tt.end)
			}
		}
	}

	sort.Sort(byValue(starts))
	sort.Sort(byValue(ends))

	hlText := ""

	//starts: [39 67 97]
	//ends: [46 74 104]

	var cursor uint32 = 0
	for i := 0; i < len(starts); i++ {
		hlText += fmt.Sprint(text[cursor:(starts[i])])
		hlText += fmt.Sprint("<b>" + text[starts[i]:ends[i]] + "</b>")
		cursor = ends[i]
	}
	hlText += fmt.Sprint(text[cursor:])
	//fmt.Println(hlText)

	return hlText
}

func (idx *BookIndex) scorePosting(postings []Posting) {
	for i := range postings {
		postings[i].boost = float32(idf(float64(len(postings)), float64(idx.NumDocs)) * tf(float64(postings[i].frequency), float64(idx.fieldLen[postings[i].docId]), idx.avgFieldLen))
		//fmt.Println(postings[i].boost)
	}
}

func (idx *BookIndex) updateAvgFieldLen() {
	total := 0

	for _, v := range idx.fieldLen {
		total += v
	}

	idx.avgFieldLen = float64(total) / float64(idx.NumDocs)
}

// tokenPositions calculate position data for each token
func tokenPositions(tokens []Token) map[string][]uint32 {
	tp := make(map[string][]uint32)

	for i := range tokens {
		tp[tokens[i].value] = append(tp[tokens[i].value], tokens[i].position)
	}

	return tp
}
