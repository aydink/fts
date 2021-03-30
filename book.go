package main

import "github.com/RoaringBitmap/roaring"

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
	tokens := idx.analyzer.Analyze(q)

	result := make([]Posting, 0)

	for i, token := range tokens {
		if i == 0 {
			result = idx.index[token.value]
		} else {
			temp := idx.index[token.value]
			result = Intersection(temp, result)
		}
	}

	return result
}

func (idx *BookIndex) buildCategoryBitmap() {

	for k, v := range idx.bookCategory {
		rb := roaring.NewBitmap()
		rb.AddMany(v)
		idx.categoryBitmaps[k] = rb
	}
}

// tokenPositions calculate position data for each token
func tokenPositions(tokens []Token) map[string][]uint32 {
	tp := make(map[string][]uint32)

	for i := range tokens {
		tp[tokens[i].value] = append(tp[tokens[i].value], tokens[i].position)
	}

	return tp
}
