package main

import (
	"fmt"
	"sort"

	"github.com/RoaringBitmap/roaring"
)

type Page struct {
	Id         int    `json:"id"`
	BookId     int    `json:"book_id"`
	Content    string `json:"content"`
	PageNumber int    `json:"page_number"`
}

type PageIndex struct {
	docId   int
	NumDocs int
	index   map[string][]Posting

	// docCategories will store docIds for every document that blongs to a category
	bookId []int

	// roaring bitmaps to store bookCategory bitmaps
	categoryBitmaps map[string]*roaring.Bitmap

	// store page content for future use
	pageStore []Page

	// store field length in number of tokens
	fieldLen []int

	// avarage field length
	avgFieldLen float64

	// Analyzer to use for text analysis and tokenization
	analyzer Analyzer
}

func NewPageIndex(analyzer Analyzer) *PageIndex {
	idx := &PageIndex{}
	idx.docId = 0

	idx.index = make(map[string][]Posting)
	idx.bookId = make([]int, 0)

	idx.categoryBitmaps = make(map[string]*roaring.Bitmap)

	// store field length in number of tokens
	idx.fieldLen = make([]int, 0)

	// store page content for future use
	idx.pageStore = make([]Page, 0)

	idx.analyzer = analyzer
	return idx
}

func (idx *PageIndex) Add(doc *Page) {
	// Set the docId of the document. It will be used as primary key for almost everyhing
	doc.Id = idx.docId

	idx.bookId = append(idx.bookId, doc.BookId)
	// Start the Analysis process
	//idx.analyze(doc)
	tokens := idx.analyzer.Analyze(doc.Content)

	for key, val := range tokenPositions(tokens) {
		//fmt.Println(key, val)
		posting := Posting{uint32(doc.Id), uint32(len(val)), 1.0, val}
		idx.index[key] = append(idx.index[key], posting)
	}

	// increment docId after ever document
	idx.docId++

	idx.pageStore = append(idx.pageStore, *doc)

	idx.fieldLen = append(idx.fieldLen, len(tokens))

	// increment total number of documents in index
	idx.NumDocs++
}

func (idx *PageIndex) Original_Search(q string) []Posting {
	tokens := idx.analyzer.Analyze(q)

	var result []Posting
	var temp []Posting
	//result := make([]Posting, 0)

	for i, token := range tokens {
		if i == 0 {
			//result = idx.index[token.value]
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
			// result = Intersection(temp, result)
			// boolean OR query
			//result = Union(temp, result)
			// Phrase Query
			result = PhraseQuery_FullMatch(result, temp)
		}
	}

	//idx.getFacetCounts(result)

	/*
		if len(result) > 100 {
			result = result[0:100]
		}
	*/

	//fmt.Println(result)
	sort.Sort(ByBoost(result))
	//fmt.Println("-------------------------------------------------")
	//fmt.Println(result)

	return result
}

// TODO
func (idx *PageIndex) Search(q string) []Posting {
	tokens := idx.analyzer.Analyze(q)

	var result []Posting
	var temp []Posting
	var resultPhrase []Posting

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
			// result = Intersection(temp, result)
			// boolean OR query
			//result = Union(temp, result)
			// Phrase Query
			resultPhrase = PhraseQuery_FullMatch(result, temp)
		}
	}

	result = Union(result, resultPhrase)

	//fmt.Println(result)
	sort.Sort(ByBoost(result))
	//fmt.Println("-------------------------------------------------")
	//fmt.Println(result)

	return result
}

func (idx *PageIndex) updateAvgFieldLen() {
	total := 0

	for _, v := range idx.fieldLen {
		total += v
	}

	idx.avgFieldLen = float64(total) / float64(idx.NumDocs)
}

func (idx *PageIndex) GetPage(pageId int) Page {
	return idx.pageStore[pageId]
}

func (idx *PageIndex) scorePosting(postings []Posting) {
	for i := range postings {
		postings[i].boost = float32(idf(float64(len(postings)), float64(idx.NumDocs)) * tf(float64(postings[i].frequency), float64(idx.fieldLen[postings[i].docId]), idx.avgFieldLen))
		//fmt.Println(postings[i].boost)
	}
}

func (idx *PageIndex) buildCategoryBitmap(bookIndex *BookIndex) {

	for k, v := range bookIndex.categoryBitmaps {
		rb := roaring.NewBitmap()

		for _, page := range idx.pageStore {
			if v.Contains(uint32(page.BookId)) {
				rb.AddInt(page.Id)
			}
		}

		idx.categoryBitmaps[k] = rb
	}
}

func (idx *PageIndex) getFacetCounts(postings []Posting) []FacetCount {
	facetCounts := make([]FacetCount, 0)

	rb := roaring.NewBitmap()
	for _, posting := range postings {
		rb.Add(posting.docId)
	}

	for k, v := range idx.categoryBitmaps {
		fc := FacetCount{}
		fc.Name = k
		fc.Count = int(v.AndCardinality(rb))

		// add only if facet count is not zero
		if fc.Count > 0 {
			facetCounts = append(facetCounts, fc)
		}
	}

	sort.Sort(byFacetCount(facetCounts))
	fmt.Printf("%+v\n", facetCounts)

	return facetCounts
}

func (idx *PageIndex) facetFilterCategory(postings []Posting, category string) []Posting {

	result := make([]Posting, 0)
	rb := idx.categoryBitmaps[category]

	for _, posting := range postings {
		if rb.Contains(posting.docId) {
			result = append(result, posting)
		}
	}
	return result
}

func (idx *PageIndex) tokenStats() []FacetCount {

	stats := make([]FacetCount, 0)

	for k, v := range idx.index {
		fc := FacetCount{}
		fc.Name = k

		counter := 0
		for _, posting := range v {
			counter += int(posting.frequency)
		}

		fc.Count = counter

		stats = append(stats, fc)
	}

	sort.Sort(byFacetCount(stats))
	return stats
}
