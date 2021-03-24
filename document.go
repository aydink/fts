package main

//"fmt"

type Document struct {
	docId      uint32
	fields     []TextField
	categories map[string]bool
	//docValues  map[string]interface{}
	docValues map[string]DocValue
}

func NewDocument() *Document {
	//return &Document{0, nil, make(map[string]bool), make(map[string]interface{})}
	return &Document{0, nil, make(map[string]bool), make(map[string]DocValue)}

}

// add fields to the document
func (doc *Document) Add(field TextField) {
	doc.fields = append(doc.fields, field)
}

// add a category to the document. it is intended for facet filtering and counting
func (doc *Document) AddCatergory(category string) {
	doc.categories[category] = true
}

// add a primitive value to the document. such as int32, int64, float32, float64
func (doc *Document) AddDocValue(name string, docValue DocValue) {
	doc.docValues[name] = docValue
}

// return fields fo the document
func (doc *Document) Fields() []TextField {
	return doc.fields
}

// return categories of the document
func (doc *Document) Categories() map[string]bool {
	return doc.categories
}

// return docValues of the document
//func (doc *Document) DocValues() map[string]interface{} {
func (doc *Document) DocValues() map[string]DocValue {
	return doc.docValues
}

// tokenPositions calculate position data for each token
func (doc *Document) tokenPositions(tokens []Token) map[string][]uint32 {
	tp := make(map[string][]uint32)

	for i := range tokens {
		tp[tokens[i].value] = append(tp[tokens[i].value], tokens[i].position)
	}

	return tp
}
