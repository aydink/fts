package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

var bookStore []Book
var pageStore []Page
var bookIndex *BookIndex
var pageIndex *PageIndex

var sentenceTokenizer *sentences.DefaultSentenceTokenizer

func buildIndex() {
	analyzer := NewSimpleAnalyzer(NewSimpleTokenizer())

	turkishFilter := NewTurkishLowercaseFilter()
	turkishAccentFilter := NewTurkishAccentFilter()

	analyzer.AddTokenFilter(turkishFilter)
	analyzer.AddTokenFilter(turkishAccentFilter)

	bookIndex = NewBookIndex(analyzer)
	pageIndex = NewPageIndex(analyzer)

	books, err := prepareBooks("xliste.csv")
	if err != nil {
		fmt.Println("Failed to load book list csv file")
		return
	}

	for _, book := range books {
		indexBook(book)
	}

	pageIndex.updateAvgFieldLen()
}

func main() {

	log.SetFlags(log.Llongfile)

	buildIndex()

	var err error
	sentenceTokenizer, err = english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/reindex", reindexHandler)

	http.HandleFunc("/search/", searchHandler)
	http.HandleFunc("/page", pageHandler)
	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/download/", downloadHandler)
	//http.HandleFunc("/api/addbook", indexFileHandler)
	//http.HandleFunc("/api/payloads", payloadHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}