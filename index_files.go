package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func prepareBooks(csvFile string) ([]Book, error) {

	//file, err := os.Open("mehaz/" + csvFile)
	file, err := os.Open("mehaz/" + csvFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	r := csv.NewReader(file)
	r.Comma = ';'
	r.Comment = '#'

	books := make([]Book, 0)

	for {
		record, err := r.Read()
		log.Println(record)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		book := Book{}
		book.Title = record[1] + " " + record[2]
		year, err := strconv.ParseUint(record[3], 10, 32)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		book.Year = int(year)
		book.Genre = record[4]
		book.Category = append(book.Category, record[5])

		//hash, err := preparePdfFile("mehaz/" + record[0])
		hash, err := preparePdfFile(record[0])
		if err != nil {
			log.Println(err)
			continue
			//return nil, err
		}

		book.Hash = hash
		books = append(books, book)

		processPdfFile(book)

		// save book struct as json file
		saveBookMeta(book)
	}

	return books, nil
}

func saveBookMeta(book Book) error {

	bookJson, err := json.Marshal(book)
	if err != nil {
		return err
	}

	file, err := os.Create("books/" + book.Hash + ".json")
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.Write(bookJson)
	if err != nil {
		return err
	}

	return nil
}

func loadBookMeta(filename string) (Book, error) {

	book := Book{}

	file, err := os.Open("books/" + filename)
	defer file.Close()
	if err != nil {
		return book, err
	}

	bookJson, err := ioutil.ReadAll(file)
	if err != nil {
		return book, err
	}

	err = json.Unmarshal(bookJson, &book)
	if err != nil {
		return book, err
	}

	return book, err
}

func reindexAllFiles() {
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
			indexBook(book)

			//store payload data in cdb file
			//ProcessPayloadFile(book.Hash)
		}
	}
}

func indexBook(book Book) {
	pages, err := getPages(book)
	book.NumPages = len(pages)

	bookIndex.Add(&book)

	if err != nil {
		fmt.Println(err)
	}

	for _, page := range pages {
		page.BookId = book.Id
		pageIndex.Add(&page)
	}

	//processPayload(book.Hash)
	//syncPayload()
}

/*
func indexBookWithPayloads(book Book) error {

	var err error
	var payloads []map[string][][4]int

	pages, err := getPages(book)
	if err != nil {
		log.Println(err)
		return err
	}

	payloads, err = CreatePayload(book.Hash)

	if err != nil {
		log.Println(err)
		return err
	}
	book.NumPages = len(pages)

	bookIndex.Add(&book)

	if len(payloads) != book.NumPages {
		log.Println(fmt.Errorf("Book has %d pages but payloads has %d pages, they must be equal", len(payloads), book.NumPages))
		return err
	}

	for _, page := range pages {
		page.BookId = book.Id

		pageIndex.Add(&page)

		//payloadStore.Put(gobEncode(payloads[page]))
	}

	processPayload(book.Hash)
	syncPayload()
}
*/

func gobEncode(tokens map[string][][4]int) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tokens); err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
