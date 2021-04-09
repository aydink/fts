package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/english"
)

var bookIndex *BookIndex
var pageIndex *PageIndex

var payloadStore *CdbStore

//var payloadStore *PayloadStore

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

	var err error
	/*
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

	payloadStore, err = NewCdbStore()
	if err != nil {
		log.Println("Failed to create cdb file")
		return
	}

	payloadStore.BuildDatabase()
	payloadStore.Freeze()

	//CratePagePayloadDatabase()
	//fmt.Println(LoadPagePayload("a33a19469bfa738e5292140fea7cea6f-21"))
}

func cleanUpBeforeExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			fmt.Println(sig.String(), "Ctrl-C captured")
			//fmt.Println("Closing cdb database")
			//pg.Close()
			os.Exit(0)
		}
	}()
}

func main() {

	f, err := os.OpenFile("out.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.SetFlags(log.Llongfile)

	// capture Ctrl-C exit event
	cleanUpBeforeExit()

	// build fulltext index
	buildIndex()

	runtime.GC()

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
