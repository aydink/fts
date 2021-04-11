package main

import (
	"bytes"
	"encoding/gob"
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
