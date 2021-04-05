package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

var bookIndex *BookIndex
var pageIndex *PageIndex

var sentenceTokenizer *sentences.DefaultSentenceTokenizer

func buildIndex() {
	analyzer := NewSimpleAnalyzer(NewSimpleTokenizer())

	turkishFilter := NewTurkishLowercaseFilter()
	turkishAccentFilter := NewTurkishAccentFilter()
	turkishStemFilter := NewTurkishStemFilter()

	analyzer.AddTokenFilter(turkishFilter)
	analyzer.AddTokenFilter(turkishAccentFilter)
	analyzer.AddTokenFilter(turkishStemFilter)

	bookIndex = NewBookIndex(analyzer)
	pageIndex = NewPageIndex(analyzer)

	/*
		books, err := prepareBooks("li.csv")
		if err != nil {
			fmt.Println("Failed to load book list csv file")
			return
		}

		for _, book := range books {
			indexBook(book)
		}

		books, err := prepareBooks("xliste.csv")
		if err != nil {
			fmt.Println("Failed to load book list csv file", err)
			return
		}

		for _, book := range books {
			indexBook(book)
		}
	*/

	reindexAllFiles()

	pageIndex.updateAvgFieldLen()
	bookIndex.updateAvgFieldLen()

	bookIndex.buildCategoryBitmap()
	pageIndex.buildCategoryBitmap(bookIndex)

	//CratePagePayloadDatabase()
	//fmt.Println(LoadPagePayload("a33a19469bfa738e5292140fea7cea6f-21"))
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
	http.HandleFunc("/stats", tokenStatHandler)
	http.HandleFunc("/api/addbook", indexFileHandler)
	http.HandleFunc("/api/payloads", payloadHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
