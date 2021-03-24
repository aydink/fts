package main

type Posting struct {
	docId     uint32
	frequency uint32
	boost     float32
	positions []uint32
}

type FieldMeta struct {
	fieldId   uint32
	name      string
	fieldType uint32
}

type Term struct {
	Value       string  // string representaion of the Term
	Idf         float32 // Inverse Document Frequency of the Term
	numPostings uint32  // total number of positions for a term, index wide.
	Postings    []Posting
}
