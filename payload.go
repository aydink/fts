package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akrylysov/pogreb"
	"github.com/colinmarc/cdb"
	"golang.org/x/net/html"
)

type BBox struct {
	XMin int16
	YMin int16
	XMax int16
	YMax int16
}

type Payload struct {
	Key   string              `json:"key"`
	Value map[string][][4]int `json:"value"`
}

/*
func GetTokenPositions(page string, q string) string {

	tokens := pageIndex.analyzer.Analyze(q)

	filteredTokens := make(map[string][][4]int)
	pageTokens := loadPayload(page)

	for _, token := range tokens {
		filteredTokens[token.value] = pageTokens[token.value]
	}

	jsonString, err := json.Marshal(filteredTokens)
	if err != nil {
		log.Println(err)
	}

	//fmt.Println(string(jsonString))
	return string(jsonString)

}
*/

var db *Database

type Database struct {
	pg *pogreb.DB
}

func NewDatabase() *Database {
	d := &Database{}
	var err error
	d.pg, err = pogreb.Open("payload", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return d
}

func init() {
	//initialize payload database
	//which use pogreb
	db = NewDatabase()
}

func (db *Database) SavePayload(key string, buf []byte) error {
	ok, err := db.pg.Has([]byte(key))
	if err != nil {
		log.Println(err)
		return err
	}

	if !ok {
		err = db.pg.Put([]byte(key), buf)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (db *Database) Sync() error {
	err := db.pg.Sync()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *Database) LoadPayload(key string) map[string][][4]int {

	m := make(map[string][][4]int)

	v, err := db.pg.Get([]byte(key))
	if err != nil {
		log.Println(err)
		return m
	}

	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&m); err != nil {
		log.Println(err)
		return m
	}

	return m
}

func (db *Database) ProcessBook(hash string) error {

	file, err := os.Open("books/" + hash + ".bbox.gob")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	payloads := make(map[string]map[string][][4]int)

	dec := gob.NewDecoder(file)

	if err := dec.Decode(&payloads); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Book key:%s, number of pages:%d", hash, len(payloads))
	//log.Printf("%+v", payloads)

	for key, tokens := range payloads {

		var buf bytes.Buffer

		enc := gob.NewEncoder(&buf)
		err := enc.Encode(tokens)

		if err != nil {
			log.Println(err)
			return err
		}

		//log.Printf("Saving payloads key:%s, length:%d", key, buf.Len())

		err = db.SavePayload(key, buf.Bytes())
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

/*
func processPayload(hash string) {

	var pageNumber int
	file, err := os.Open("books/" + hash + ".bbox.txt")
	if err != nil {
		log.Println(err)
	}

	z := html.NewTokenizer(file)

	var bbox [4]int
	var tokens map[string][][4]int

	insideWord := false

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			//postToElasticsearch(buf.Bytes())
			return

		case html.StartTagToken:
			t := z.Token()

			if t.Data == "page" {
				pageNumber++
				//fmt.Println(pageNumber, "------------------------------")

				tokens = make(map[string][][4]int)
			}

			if t.Data == "word" {

				//bbox = BBox{}
				bbox = [4]int{}

				for _, w := range t.Attr {
					n, err := strconv.ParseFloat(w.Val, 64)
					if err != nil {
						log.Println(err)
					}
					n = math.Floor(n + 0.5)
					coor := int(n)

					switch w.Key {
					case "xmin":
						bbox[0] = coor
					case "ymin":
						bbox[1] = coor
					case "xmax":
						bbox[2] = coor
					case "ymax":
						bbox[3] = coor
					}
				}

				insideWord = true
			} else {
				insideWord = false
			}

		case html.TextToken:
			if insideWord {
				token := strings.TrimSpace(z.Token().Data)

				// use the same analyzer with pageIndex
				parts := pageIndex.analyzer.Analyze(token)
				for _, v := range parts {
					tokens[v.value] = append(tokens[v.value], bbox)
				}
			}

		case html.EndTagToken:
			t := z.Token()
			if t.Data == "page" {
				//fmt.Println("end page:", pageNumber)
				//fmt.Println(len(tokens))

				// insert Payloads into KV store. Use md5 hash and page number as key
				key := hash + "-" + strconv.Itoa(pageNumber)
				savePayload(key, tokens)

				//	fmt.Println(key)
				//	fmt.Println("----------------------------------------------------------------")
				//	fmt.Println(tokens)
				//	fmt.Println("----------------------------------------------------------------")
			}
		}
	}
}

*/

/*
ProcessPayloadFile read and stores token positions in Elasticsearch
"payload" index using "data" type. Id of the document is hash + "-" + page
sample document

{
	"key": "md5 hash of the book",
	"value" : {
		"token1": [[1,2,3,4], [4,5,6,7]],
		"token2": [[11,12,13,14], [14,15,16,17], [8,9,11,13]]
	}
}

*/

/*

func CratePagePayloadDatabase() {

	writer, err := cdb.Create("data/page_payload.cdb")
	if err != nil {
		log.Fatal(err)
	}

	fileInfos, err := ioutil.ReadDir("books")
	if err != nil {
		log.Printf("opening books directory failed.")
		return
	}

	for _, file := range fileInfos {
		if filepath.Ext(file.Name()) == ".json" {
			book, err := loadBookMeta(file.Name())
			if err != nil {
				log.Printf("loading file meta from json file:%s faied\n", err)
				continue
			}
			log.Println(book)
			ProcessPayloadFile(writer, book.Hash)
		}
	}

	_, err = writer.Freeze()
	if err != nil {
		log.Fatal(err)
	}
}

func LoadPagePayload(key string) map[string][][4]int {

	m := make(map[string][][4]int)

	db, err := cdb.Open("data/page_payload.cdb")
	if err != nil {
		log.Println(err)
		return m
	}

	v, err := db.Get([]byte(key))
	if err != nil {
		log.Println(err)
		return m
	}

	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&m); err != nil {
		log.Println(err)
		return m
	}

	return m
}

func GetPagePayloadJSON(tokens map[string][][4]int) string {

	payloadJson, err := json.Marshal(tokens)
	if err != nil {
		log.Fatalf("failed the marshall payloads:%s\n", err)
	}

	//fmt.Println(string(payloadJson))
	return string(payloadJson)

}

func ProcessPayloadFile(writer *cdb.Writer, hash string) {

	var pageNumber int
	file, err := os.Open("books/" + hash + ".bbox.txt")
	if err != nil {
		log.Println(err)
	}

	z := html.NewTokenizer(file)

	var bbox [4]int
	var tokens map[string][][4]int

	insideWord := false

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			//postToElasticsearch(buf.Bytes())
			return

		case html.StartTagToken:
			t := z.Token()

			if t.Data == "page" {
				pageNumber++
				//fmt.Println(pageNumber, "------------------------------")

				tokens = make(map[string][][4]int)
			}

			if t.Data == "word" {

				//bbox = BBox{}
				bbox = [4]int{}

				for _, w := range t.Attr {
					n, err := strconv.ParseFloat(w.Val, 64)
					if err != nil {
						log.Println(err)
					}
					n = math.Floor(n + 0.5)
					coor := int(n)

					switch w.Key {
					case "xmin":
						bbox[0] = coor
					case "ymin":
						bbox[1] = coor
					case "xmax":
						bbox[2] = coor
					case "ymax":
						bbox[3] = coor
					}
				}

				insideWord = true
			} else {
				insideWord = false
			}

		case html.TextToken:
			if insideWord {
				token := strings.TrimSpace(z.Token().Data)

				// use the same analyzer with pageIndex
				parts := pageIndex.analyzer.Analyze(token)
				for _, v := range parts {
					tokens[v.value] = append(tokens[v.value], bbox)
				}
			}

		case html.EndTagToken:
			t := z.Token()
			if t.Data == "page" {
				//fmt.Println("end page:", pageNumber)
				//fmt.Println(len(tokens))

				// insert Payloads into KV store. Use md5 hash and page number as key
				key := hash + "-" + strconv.Itoa(pageNumber)
				StorePagePayload(writer, key, tokens)

				//	fmt.Println(key)
				//	fmt.Println("----------------------------------------------------------------")
				//	fmt.Println(tokens)
				//	fmt.Println("----------------------------------------------------------------")

			}
		}
	}
}


func StorePagePayload(writer *cdb.Writer, key string, tokens map[string][][4]int) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tokens); err != nil {
		log.Fatal(err)
	}

	writer.Put([]byte(key), buf.Bytes())
}
*/
//---------------------------------------------------------
// Plain binary implementation for storing payloads
//---------------------------------------------------------

type PayloadStore struct {
	file          *os.File
	currentOffset int64
}

func NewPayloadStore() *PayloadStore {
	file, err := os.Create("payload/payload.bin")
	if err != nil {
		log.Fatal(err)
	}
	return &PayloadStore{file: file, currentOffset: 0}
}

func (ps *PayloadStore) Put(payload []byte) int64 {
	n, err := ps.file.Write(payload)
	if err != nil {
		log.Println(err)
	}

	ps.currentOffset += int64(n)
	return ps.currentOffset

}
func (ps *PayloadStore) Get(offset int64, length int64) ([]byte, error) {
	if offset < ps.currentOffset {
		return nil, fmt.Errorf("requested offset=%d is not valid", offset)
	}
	ps.file.Seek(offset, 0)

	buf := make([]byte, length)

	_, err := ps.file.Read(buf)
	if err != nil {
		log.Println(err)
	}
	return buf, nil
}

func (ps *PayloadStore) Close() {
	err := ps.file.Close()
	if err != nil {
		log.Println(err)
	}
}

func CreatePayload(hash string) ([]byte, error) {

	var pageNumber int
	file, err := os.Open("books/" + hash + ".bbox.txt")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	z := html.NewTokenizer(file)

	var bbox [4]int
	var tokens map[string][][4]int

	payloads := make(map[string]map[string][][4]int)

	insideWord := false

	processing := true

	for processing {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			//postToElasticsearch(buf.Bytes())
			//return payloads, nil
			processing = false

		case html.StartTagToken:
			t := z.Token()

			if t.Data == "page" {
				pageNumber++
				//fmt.Println(pageNumber, "------------------------------")

				tokens = make(map[string][][4]int)
			}

			if t.Data == "word" {

				//bbox = BBox{}
				bbox = [4]int{}

				for _, w := range t.Attr {
					n, err := strconv.ParseFloat(w.Val, 64)
					if err != nil {
						log.Println(err)
					}
					n = math.Floor(n + 0.5)
					coor := int(n)

					switch w.Key {
					case "xmin":
						bbox[0] = coor
					case "ymin":
						bbox[1] = coor
					case "xmax":
						bbox[2] = coor
					case "ymax":
						bbox[3] = coor
					}
				}

				insideWord = true
			} else {
				insideWord = false
			}

		case html.TextToken:
			if insideWord {
				token := strings.TrimSpace(z.Token().Data)

				// use the same analyzer with pageIndex
				parts := pageIndex.analyzer.Analyze(token)
				for _, v := range parts {
					tokens[v.value] = append(tokens[v.value], bbox)
				}
			}

		case html.EndTagToken:
			t := z.Token()
			if t.Data == "page" {
				//fmt.Println("end page:", pageNumber)
				//fmt.Println(len(tokens))

				// insert Payloads into KV store. Use md5 hash and page number as key
				key := hash + "-" + strconv.Itoa(pageNumber)
				payloads[key] = tokens
				//payloads = append(payloads, tokens)
				//savePayload(key, tokens)
				/*
					fmt.Println(key)
					fmt.Println("----------------------------------------------------------------")
					fmt.Println(tokens)
					fmt.Println("----------------------------------------------------------------")
				*/
			}
		}
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	//log.Println("------------------------------------------------------")
	//log.Println(payloads)

	if err := enc.Encode(payloads); err != nil {
		log.Println(err)
		return nil, err
	}

	return buf.Bytes(), nil
}

//-------------------------------------------------------------
// Using CDB database to storo payloads
//-------------------------------------------------------------
type CdbStore struct {
	writer *cdb.Writer
	reader *cdb.CDB
}

func NewCdbStore() (*CdbStore, error) {

	cdbStore := &CdbStore{}
	writer, err := cdb.Create("data/page_payload.cdb")
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	cdbStore.writer = writer

	return cdbStore, nil
}
func OpenCdbStore() (*CdbStore, error) {

	cdbStore := &CdbStore{}

	reader, err := cdb.Open("data/page_payload.cdb")
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	cdbStore.reader = reader
	return cdbStore, nil
}

func (db *CdbStore) BuildDatabase() {

	fileInfos, err := ioutil.ReadDir("books")
	if err != nil {
		log.Println("opening books directory failed.")
		return
	}

	for _, file := range fileInfos {
		if filepath.Ext(file.Name()) == ".json" {
			book, err := loadBookMeta(file.Name())
			if err != nil {
				log.Printf("loading file meta from json file:%s faied\n", err)
				continue
			}
			log.Println(book)
			db.ProcessBook2(book.Hash)
		}
	}
}

func (db *CdbStore) Freeze() {
	r, err := db.writer.Freeze()
	if err != nil {
		log.Println(err)
	}

	db.reader = r
}

func (db *CdbStore) ProcessBook(hash string) error {

	var pageNumber int
	file, err := os.Open("books/" + hash + ".bbox.txt")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	z := html.NewTokenizer(file)

	var bbox [4]int
	var tokens map[string][][4]int

	insideWord := false

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			//postToElasticsearch(buf.Bytes())
			return nil

		case html.StartTagToken:
			t := z.Token()

			if t.Data == "page" {
				pageNumber++
				//fmt.Println(pageNumber, "------------------------------")
				tokens = make(map[string][][4]int)
			}

			if t.Data == "word" {

				bbox = [4]int{}

				for _, w := range t.Attr {
					n, err := strconv.ParseFloat(w.Val, 64)
					if err != nil {
						log.Println(err)
					}
					n = math.Floor(n + 0.5)
					coor := int(n)

					switch w.Key {
					case "xmin":
						bbox[0] = coor
					case "ymin":
						bbox[1] = coor
					case "xmax":
						bbox[2] = coor
					case "ymax":
						bbox[3] = coor
					}
				}

				insideWord = true
			} else {
				insideWord = false
			}

		case html.TextToken:
			if insideWord {
				token := strings.TrimSpace(z.Token().Data)

				// use the same analyzer with pageIndex
				parts := pageIndex.analyzer.Analyze(token)
				for _, v := range parts {
					tokens[v.value] = append(tokens[v.value], bbox)
				}
			}

		case html.EndTagToken:
			t := z.Token()
			if t.Data == "page" {
				//fmt.Println("end page:", pageNumber)
				//fmt.Println(len(tokens))

				// insert Payloads into KV store. Use md5 hash and page number as key
				key := hash + "-" + strconv.Itoa(pageNumber)
				db.SavePayload(key, tokens)

				//	fmt.Println(key)
				//	fmt.Println("----------------------------------------------------------------")
				//  fmt.Println(tokens)
				//  fmt.Println("----------------------------------------------------------------")
			}
		}
	}
}
func (db *CdbStore) ProcessBook2(hash string) error {

	file, err := os.Open("books/" + hash + ".bbox.gob")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	payloads := make(map[string]map[string][][4]int)

	dec := gob.NewDecoder(file)

	if err := dec.Decode(&payloads); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Book key:%s, number of pages:%d", hash, len(payloads))
	//log.Printf("%+v", payloads)

	for key, tokens := range payloads {

		var buf bytes.Buffer

		enc := gob.NewEncoder(&buf)
		err := enc.Encode(tokens)

		if err != nil {
			log.Println(err)
			return err
		}

		//log.Printf("Saving payloads key:%s, length:%d", key, buf.Len())

		err = db.writer.Put([]byte(key), buf.Bytes())
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (db *CdbStore) LoadPayload(key string) map[string][][4]int {

	m := make(map[string][][4]int)

	v, err := db.reader.Get([]byte(key))
	if err != nil {
		log.Println(err)
		return m
	}

	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&m); err != nil {
		log.Println(err)
		return m
	}

	return m
}

func (db *CdbStore) SavePayload(key string, tokens map[string][][4]int) error {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tokens); err != nil {
		log.Println(err)
		return err
	}

	err := db.writer.Put([]byte(key), buf.Bytes())
	if err != nil {
		log.Println(err)
	}
	return err
}

func (db *CdbStore) GetTokenPositions(page string, q string) string {

	tokens := pageIndex.analyzer.Analyze(q)

	filteredTokens := make(map[string][][4]int)
	pageTokens := db.LoadPayload(page)

	for _, token := range tokens {
		filteredTokens[token.value] = pageTokens[token.value]
	}

	jsonString, err := json.Marshal(filteredTokens)
	if err != nil {
		log.Println(err)
	}

	//fmt.Println(string(jsonString))
	return string(jsonString)

}
