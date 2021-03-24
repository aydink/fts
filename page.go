package main

import (
	"fmt"
	"sort"
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

func (idx *PageIndex) Search(q string) []Posting {
	tokens := idx.analyzer.Analyze(q)

	result := make([]Posting, 0)

	for i, token := range tokens {
		if i == 0 {
			result = idx.index[token.value]
			fmt.Println(result)
			idx.scorePosting(result)
			fmt.Println(result)
		} else {
			temp := idx.index[token.value]
			idx.scorePosting(temp)
			result = Intersection(temp, result)
		}
	}

	/*
		if len(result) > 10 {
			result = result[0:10]
		}
	*/
	fmt.Println(result)
	sort.Sort(ByBoost(result))
	fmt.Println("-------------------------------------------------")
	fmt.Println(result)

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
		fmt.Println(postings[i].boost)
	}
}
