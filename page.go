package fts

type Page struct {
	Id         uint32 `json:"id"`
	BookId     uint32 `json:"book_id"`
	Content    string `json:"content"`
	PageNumber uint32 `json:"page_number"`
}

type PageIndex struct {
	docId   uint32
	NumDocs uint32
	index   map[string][]Posting

	// docCategories will store docIds for every document that blongs to a category
	bookId []uint32

	// Analyzer to use for text analysis and tokenization
	analyzer Analyzer
}

func NewPageIndex(analyzer Analyzer) *PageIndex {
	idx := &PageIndex{}
	idx.docId = 0

	idx.index = make(map[string][]Posting)
	idx.bookId = make([]uint32, 0)

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
		posting := Posting{doc.Id, uint32(len(val)), 1.0, val}
		idx.index[key] = append(idx.index[key], posting)
	}

	// increment docId after ever document
	idx.docId++

	// increment total number of documents in index
	idx.NumDocs++
}

func (idx *PageIndex) Search(q string) []Posting {
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
