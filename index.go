package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/bkaradzic/go-lz4"
)

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

type DocValueInt32 struct {
	docId    uint32
	docValue int32
}

type DocValueInt64 struct {
	docId    uint32
	docValue int64
}

type DocValueFloat32 struct {
	docId    uint32
	docValue float32
}

type DocValueFloat64 struct {
	docId    uint32
	docValue float64
}

type InvertedIndex struct {
	docId            uint32
	NumDocs          uint32
	index            map[string][]Posting
	postingMap       map[string]uint32
	documentMap      map[uint32]uint32
	fieldIdCounter   uint32
	fieldsDictionary map[string]FieldMeta
	fieldIdMap       map[uint32]string

	indexDir    string
	fpPostings  *os.File
	fpDocuments *os.File

	documentPointer uint32

	// docCategories will store docIds for every document that blongs to a category
	docCategories map[string][]uint32

	docValuesInt32   map[string][]DocValueInt32
	docValuesInt64   map[string][]DocValueInt64
	docValuesFloat32 map[string][]DocValueFloat32
	docValuesFloat64 map[string][]DocValueFloat64

	// Analyzer to use for text analysis and tokenization
	analyzer Analyzer
}

type TermsDictionary interface {
	GetTermVector(string)
}

func New(dirName string, analyzer Analyzer) *InvertedIndex {
	idx := &InvertedIndex{}
	idx.indexDir = dirName
	idx.docId = 0

	if ok, err := exists(dirName); !ok && err == nil {
		os.Mkdir(dirName, 0777)
	} else {
		fmt.Println(err)
	}

	fpPostings, err := os.Create(dirName + "/" + "postings.bin")
	if err != nil {
		fmt.Println(err)
	}
	fpDocuments, err := os.Create(dirName + "/" + "documents.bin")
	if err != nil {
		fmt.Println(err)
	}

	idx.fpPostings = fpPostings
	idx.fpDocuments = fpDocuments

	idx.index = make(map[string][]Posting)
	idx.postingMap = make(map[string]uint32)
	idx.documentMap = make(map[uint32]uint32)

	idx.documentPointer = 0
	idx.fieldIdCounter = 0
	idx.fieldsDictionary = make(map[string]FieldMeta)
	idx.fieldIdMap = make(map[uint32]string)
	idx.docCategories = make(map[string][]uint32)

	idx.docValuesInt32 = make(map[string][]DocValueInt32)
	idx.docValuesInt64 = make(map[string][]DocValueInt64)
	idx.docValuesFloat32 = make(map[string][]DocValueFloat32)
	idx.docValuesFloat64 = make(map[string][]DocValueFloat64)

	idx.analyzer = analyzer

	return idx
}

func (idx *InvertedIndex) Debug() {
	fmt.Println("docValuesInt32")
	fmt.Println("--------------------------------------")
	fmt.Println(idx.docValuesFloat32)
	fmt.Println("--------------------------------------")

	fmt.Println("docValuesInt64")
	fmt.Println("--------------------------------------")
	fmt.Println(idx.docValuesInt64)
	fmt.Println("--------------------------------------")

	fmt.Println("docValuesFloat32")
	fmt.Println("--------------------------------------")
	fmt.Println(idx.docValuesFloat32)
	fmt.Println("--------------------------------------")

	fmt.Println("docValuesFloat64")
	fmt.Println("--------------------------------------")
	fmt.Println(idx.docValuesFloat64)
	fmt.Println("--------------------------------------")

	fmt.Println("Categories")
	fmt.Println("--------------------------------------")
	fmt.Println(idx.docCategories)

	fmt.Println("Fields")
	fmt.Println("--------------------------------------")
	fmt.Println(idx.fieldIdMap)
	fmt.Println(idx.fieldsDictionary)

}

func Open(dirName string) *InvertedIndex {
	idx := &InvertedIndex{}
	idx.indexDir = dirName

	ok, err := exists(dirName)
	if !ok {
		fmt.Println("index does not exist")
		os.Exit(1)
	} else {
		fmt.Println(err)
	}

	fpPostings, err := os.Open(dirName + "/" + "postings.bin")
	if err != nil {
		fmt.Println(err)
	}
	fpDocuments, err := os.Open(dirName + "/" + "documents.bin")
	if err != nil {
		fmt.Println(err)
	}

	idx.fpDocuments = fpDocuments
	idx.fpPostings = fpPostings

	idx.index = make(map[string][]Posting)
	idx.postingMap = make(map[string]uint32)
	idx.documentMap = make(map[uint32]uint32)

	idx.documentPointer = 0

	idx.loadDocumentMap()
	idx.loadPostingMap()

	return idx
}

func (idx *InvertedIndex) numPositions(term string) int {
	counter := 0
	for _, posting := range idx.index[term] {
		counter += len(posting.positions)
	}
	return counter
}

func (idx *InvertedIndex) idf(term string, docFreq int) float32 {
	numDocs := float64(idx.docId)
	//docFreq := float64(d)
	//fmt.Println(term)
	//fmt.Println(numDocs)
	//fmt.Println(docFreq)
	//return float32(math.Log(numDocs/(docFreq+1)) + 1)
	return float32(math.Log(numDocs/(float64(docFreq+1))) + 1)
}

func (idx *InvertedIndex) GetDocument(docId uint32) string {
	b := make([]byte, 4, 4)
	idx.fpDocuments.Seek(int64(idx.documentMap[docId]), 0)
	idx.fpDocuments.Read(b)
	l := bytesToUint32le(b)
	//fmt.Println(l)
	b = make([]byte, l, l)
	idx.fpDocuments.Read(b)
	data, _ := lz4.Decode(nil, b)
	return string(data)
}

func (idx *InvertedIndex) ReadPostings(term string) (Term, error) {
	t := Term{}

	fp, err := os.Open(idx.indexDir + "/" + "postings.bin")
	defer fp.Close()

	if err != nil {
		fmt.Errorf("hata:%s", err)
		return t, err
	}

	b := make([]byte, 4, 4)

	// locate the record
	fp.Seek(int64(idx.postingMap[term]), 0)
	// read the record length in bytes
	fp.Read(b)
	totalLengthInBytes := bytesToUint32le(b)

	// create a buffer to hold the record and fill it by reading from the file
	s := make([]byte, totalLengthInBytes-4, totalLengthInBytes-4)
	fp.Read(s)
	buffer := bytes.NewBuffer(s)

	// read Idf of the term
	buffer.Read(b)
	idf := math.Float32frombits(bytesToUint32le(b))
	fmt.Println("idf:", idf)
	t.Idf = idf

	// read of number of postings
	buffer.Read(b)
	//fmt.Println(buffer.Bytes())

	numPostings := bytesToUint32le(b)

	postings := make([]Posting, 0, numPostings)

	for i := 0; i < int(numPostings); i++ {
		p := Posting{}
		buffer.Read(b)
		p.docId = bytesToUint32le(b)
		buffer.Read(b)
		p.frequency = bytesToUint32le(b)
		buffer.Read(b)
		p.boost = math.Float32frombits(bytesToUint32le(b))

		// read term positions
		for j := 0; j < int(p.frequency); j++ {
			buffer.Read(b)
			p.positions = append(p.positions, bytesToUint32le(b))
		}
		postings = append(postings, p)
	}

	t.Value = term
	t.Postings = postings
	t.numPostings = numPostings
	//fmt.Println(postings)
	return t, nil
}

func (idx *InvertedIndex) savePostingMap() error {

	fp, err := os.Create(idx.indexDir + "/" + "postings.map")
	defer fp.Close()

	if err != nil {
		fmt.Errorf("hata:%s", err)
		return err
	}

	for key, position := range idx.postingMap {
		fp.WriteString(key + "\t" + strconv.Itoa(int(position)) + "\n")
	}
	fmt.Println("posting map saved")
	return nil
}

func (idx *InvertedIndex) saveDocumentMap() error {

	fp, err := os.Create(idx.indexDir + "/" + "documents.map")
	defer fp.Close()

	if err != nil {
		fmt.Errorf("hata:%s", err)
		return err
	}

	for key, position := range idx.documentMap {
		fp.WriteString(strconv.Itoa(int(key)) + "\t" + strconv.Itoa(int(position)) + "\n")
	}
	fmt.Println("document map saved")
	return nil
}

func (idx *InvertedIndex) loadPostingMap() error {

	fp, err := os.Open(idx.indexDir + "/" + "postings.map")
	defer fp.Close()

	if err != nil {
		fmt.Errorf("hata:%s", err)
		return err
	}

	r := bufio.NewReader(fp)
	for line, _, err := r.ReadLine(); err != io.EOF; line, _, err = r.ReadLine() {
		s := strings.Split(string(line), "\t")
		p, _ := strconv.Atoi(s[1])
		idx.postingMap[s[0]] = uint32(p)
	}

	fmt.Println("posting map loaded")
	return nil
}

func (idx *InvertedIndex) loadDocumentMap() error {

	fp, err := os.Open(idx.indexDir + "/" + "documents.map")
	defer fp.Close()

	if err != nil {
		fmt.Errorf("hata:%s", err)
		return err
	}

	r := bufio.NewReader(fp)
	for line, _, err := r.ReadLine(); err != io.EOF; line, _, err = r.ReadLine() {
		s := strings.Split(string(line), "\t")
		k, _ := strconv.Atoi(s[0])
		p, _ := strconv.Atoi(s[1])
		idx.documentMap[uint32(k)] = uint32(p)
	}

	fmt.Println("document map loaded")
	return nil
}

func (idx *InvertedIndex) Add(doc *Document) {
	// Set the docId of the document. It will be used as primary key for almost everyhing
	doc.docId = idx.docId

	// Start the Analysis process
	idx.analyze(doc)

	for i, _ := range doc.fields {
		fname := doc.fields[i].Name()
		if val, ok := idx.fieldsDictionary[fname]; ok {
			// you cannot change the fieldtyepe of a field once you start using it
			if doc.fields[i].fieldType != val.fieldType {
				panic("Cannot have different types of fields with the same name. Changing field type after first use is not supported.")
			}

			doc.fields[i].fid = val.fieldId

		} else {
			fm := FieldMeta{}
			fm.name = doc.fields[i].name
			fm.fieldId = idx.fieldIdCounter
			fm.fieldType = doc.fields[i].fieldType
			idx.fieldsDictionary[fname] = fm
			// set field ids of fields. we will need it when we are storing fields
			doc.fields[i].fid = fm.fieldId
			idx.fieldIdCounter++

			// fieldId to fieldname resolution
			idx.fieldIdMap[fm.fieldId] = fname
		}
	}

	// add document categories to index
	for k, _ := range doc.categories {
		idx.docCategories[k] = append(idx.docCategories[k], doc.docId)
	}

	// Once again changing docValue type associated with the field name is not allowed.
	// lets say you started "age" field with Int32, you cannot change it to Int64 or any other type once you enter Int32
	// add docValues to the index
	for key, val := range doc.docValues {
		if v, ok := idx.fieldsDictionary[key]; ok {
			if idx.fieldsDictionary[key].fieldType != v.fieldType {
				panic("Cannot change the type of docValue field associated with field name.")
			}
			//doc.fields[i].fid = val.fieldId
		} else {
			fm := FieldMeta{}
			fm.name = key
			fm.fieldId = idx.fieldIdCounter
			fm.fieldType = val.Type()
			idx.fieldsDictionary[key] = fm
			idx.fieldIdCounter++

			// fieldId to fieldname resolution
			idx.fieldIdMap[fm.fieldId] = key
		}

		switch v := val.(type) {
		case Int32:
			//fmt.Println(v)
			idx.docValuesInt32[key] = append(idx.docValuesInt32[key], DocValueInt32{doc.docId, int32(v)})
		case Int64:
			//fmt.Println(v)
			idx.docValuesInt64[key] = append(idx.docValuesInt64[key], DocValueInt64{doc.docId, int64(v)})
		case Float32:
			//fmt.Println(v)
			idx.docValuesFloat32[key] = append(idx.docValuesFloat32[key], DocValueFloat32{doc.docId, float32(v)})
		case Float64:
			//fmt.Println(v)
			idx.docValuesFloat64[key] = append(idx.docValuesFloat64[key], DocValueFloat64{doc.docId, float64(v)})
		default:
			panic("Unsupported document value type. Only int32, int64, float32 and float64 are supported.")
		}
	}

	for i, _ := range doc.fields {
		for key, val := range doc.tokenPositions(doc.fields[i].Tokens()) {
			//fmt.Println(key, val)
			posting := Posting{doc.docId, uint32(len(val)), 1.0, val}
			//idx.index[doc.fields[i].name+":"+key] = append(idx.index[key], posting)
			termKey := doc.fields[i].name + ":" + key

			if _, ok := idx.index[doc.fields[i].name+":"+key]; ok {
				idx.index[termKey] = append(idx.index[termKey], posting)
			} else {
				idx.index[termKey] = []Posting{posting}
			}
		}
	}

	// Store the field values
	idx.storeDocument(doc)

	// increment docId after ever document
	idx.docId++
}

// Analyze document and make it ready for inclusion into the inverted index
func (idx *InvertedIndex) analyze(doc *Document) {
	for i, _ := range doc.fields {
		if doc.fields[i].Analyzed() {
			doc.fields[i].SetTokens(idx.analyzer.Analyze(doc.fields[i].value))
		} else {
			doc.fields[i].SetTokens([]Token{{0, uint32(len(doc.fields[i].value)), 0, doc.fields[i].value}})
		}
	}
}

/*
StoreDocument saves content of the "stored" fields to disk

Storage Format
number of fileds (uint32) | fieldId (uint32) | field length (uint32) | field value (string)

*/
func (idx *InvertedIndex) storeDocument(doc *Document) {
	var buf bytes.Buffer
	var numFields uint32 = 0

	// calculate the number of stored fields
	for i, _ := range doc.fields {
		if doc.fields[i].stored {
			numFields++
		}
	}

	buf.Write(uint32ToBytes(numFields))

	// store pointer to the document
	idx.documentMap[doc.docId] = idx.documentPointer
	for i, _ := range doc.fields {
		if doc.fields[i].stored {
			// write fieldId
			buf.Write(uint32ToBytes(doc.fields[i].fid))
			// write field length
			buf.Write(uint32ToBytes(uint32(len(doc.fields[i].value))))
			// write the actual field data
			buf.WriteString(doc.fields[i].value)
		}
	}

	// write the stored field content to the file
	idx.fpDocuments.Write(buf.Bytes())

	// Store the document location for future use. You will use it when you want to fetch the stored field data
	idx.documentMap[doc.docId] = idx.documentPointer

	// update the current file pointer to indicate the starting position of the next document
	idx.documentPointer += uint32(buf.Len())

	// // store docValues
	// for key, val := range doc.docValues {
	// 	switch idx.fieldsDictionary[key].fieldType {
	// 	case FIELD_TYPE_INT32:
	// 		// write the fieldId of the docValue field
	// 		buf.Write(uint32ToBytes(idx.fieldsDictionary[key].fieldId))
	// 		buf.Write(int32ToBytes(int32(val.(type))))
	// 	case FIELD_TYPE_INT64:
	// 		// write the fieldId of the docValue field
	// 		buf.Write(uint32ToBytes(idx.fieldsDictionary[key].fieldId))
	// 		buf.Write(int32ToBytes(int32(val.(type))))
	// 	}
	// }

	// store document categories
}

// Commit saves inmemory inveted index to disk
func (idx *InvertedIndex) Commit() error {
	p := 0

	defer idx.fpPostings.Close()

	for key, postings := range idx.index {
		//fmt.Println(key, p)
		//buffer := new(bytes.Buffer)
		// 4 bytes uint32 for total length of byte slice
		// 4 bytes float32 Idf of the term
		// 4 bytes uint42 number of postings

		totalLengthInBytes := (4 + 4 + 4) + (len(postings) * 12) + (idx.numPositions(key) * 4)

		buffer := bytes.NewBuffer(make([]byte, 0, totalLengthInBytes))

		// write total length of byte slice that represent both postings and positions
		buffer.Write(uint32ToBytes(uint32(totalLengthInBytes)))
		// write idf of the term
		buffer.Write(uint32ToBytes(math.Float32bits(idx.idf(key, len(idx.index[key])))))
		// write number of postings
		buffer.Write(uint32ToBytes(uint32(len(idx.index[key]))))

		for _, posting := range postings {
			buffer.Write(uint32ToBytes(uint32(posting.docId)))
			buffer.Write(uint32ToBytes(uint32(posting.frequency)))
			buffer.Write(uint32ToBytes(math.Float32bits(posting.boost)))

			for _, position := range posting.positions {
				buffer.Write(uint32ToBytes(position))
			}
		}

		idx.postingMap[key] = uint32(p)
		idx.fpPostings.Write(buffer.Bytes())
		p += totalLengthInBytes
		//fmt.Println(buffer.Bytes())
	}
	//fmt.Println(postingMap)
	fmt.Println("index serialization finished")

	idx.savePostingMap()
	idx.saveDocumentMap()

	return nil
}
