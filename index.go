package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"sort"

	"github.com/colinmarc/cdb"
)

type Posting struct {
	docId     uint32
	frequency uint32
	boost     float32
	positions []uint32
}

type Term struct {
	Value       string  // string representaion of the Term
	Idf         float32 // Inverse Document Frequency of the Term
	numPostings uint32  // total number of positions for a term, index wide.
	Postings    []Posting
}

func (p *Posting) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(p.docId)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(p.frequency)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(p.boost)

	if err != nil {
		return nil, err
	}
	err = encoder.Encode(p.positions)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (p *Posting) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&p.docId)
	if err != nil {
		return err
	}
	err = decoder.Decode(&p.frequency)
	if err != nil {
		return err
	}
	err = decoder.Decode(&p.boost)
	if err != nil {
		return err
	}
	return decoder.Decode(&p.positions)
}

func (idx *PageIndex) MarshalIndex() error {

	writer, err := cdb.Create("data/page_index.cdb")
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range idx.index {

		buf := new(bytes.Buffer)
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(v)
		if err != nil {
			log.Println(err)
			return err
		}

		writer.Put([]byte(k), buf.Bytes())
	}

	writer.Freeze()
	writer.Close()

	return nil
}

func (idx *BookIndex) MarshalIndex() error {

	writer, err := cdb.Create("data/book_index.cdb")
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range idx.index {

		buf := new(bytes.Buffer)
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(v)
		if err != nil {
			log.Println(err)
			return err
		}

		writer.Put([]byte(k), buf.Bytes())
	}

	writer.Freeze()
	writer.Close()

	return nil
}

func ReadPosting(term string) []Posting {

	reader, err := cdb.Open("data/page_index.cdb")
	if err != nil {
		log.Println(err)
	}

	defer reader.Close()

	b, err := reader.Get([]byte(term))
	if err != nil {
		log.Println(err)
	}

	buf := bytes.NewBuffer(b)

	var postings []Posting
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&postings)
	if err != nil {
		log.Println(err)
	}

	return postings
}

func ReadBookPosting(term string) []Posting {

	reader, err := cdb.Open("data/book_index.cdb")
	if err != nil {
		log.Println(err)
	}

	defer reader.Close()

	b, err := reader.Get([]byte(term))
	if err != nil {
		log.Println(err)
	}

	buf := bytes.NewBuffer(b)

	var postings []Posting
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(&postings)
	if err != nil {
		log.Println(err)
	}

	return postings
}

// TODO
func (idx *PageIndex) Search_Cdb(q string) []Posting {
	tokens := idx.analyzer.Analyze(q)

	var result []Posting
	var temp []Posting
	var resultPhrase []Posting

	for i, token := range tokens {
		if i == 0 {
			result = ReadPosting(token.value)
			idx.scorePosting(result)
			//fmt.Println(result)
		} else {
			//temp := idx.index[token.value]
			temp = ReadPosting(token.value)
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
			resultPhrase = ReadPosting(token.value)
			idx.scorePosting(result)
			//fmt.Println(result)
		} else {
			//temp := idx.index[token.value]
			temp = ReadPosting(token.value)
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

// TODO
func (idx *BookIndex) Search_Cdb(q string) []Posting {
	tokens := idx.analyzer.Analyze(q)

	var result []Posting
	var temp []Posting
	var resultPhrase []Posting

	for i, token := range tokens {
		if i == 0 {
			result = ReadBookPosting(token.value)
			idx.scorePosting(result)
			//fmt.Println(result)
		} else {
			//temp := idx.index[token.value]
			temp = ReadBookPosting(token.value)
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
			resultPhrase = ReadBookPosting(token.value)
			idx.scorePosting(result)
			//fmt.Println(result)
		} else {
			//temp := idx.index[token.value]
			temp = ReadBookPosting(token.value)
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
