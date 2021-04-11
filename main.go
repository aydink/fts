package main

import (
	"encoding/gob"
	"flag"
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

	indexFiles()
	buildPayloadDatabase()

	pageIndex.updateAvgFieldLen()
	bookIndex.updateAvgFieldLen()

	bookIndex.buildCategoryBitmap()
	pageIndex.buildCategoryBitmap(bookIndex)

	SavePageIndex()
}

func buildPayloadDatabase() {

	var err error

	if *flagBuildPayload {

		payloadStore, err = NewCdbStore()
		if err != nil {
			log.Println("Failed to create cdb file")
			return
		}

		payloadStore.BuildDatabase()
		payloadStore.Freeze()
	} else {

		payloadStore, err = OpenCdbStore()
		if err != nil {
			log.Println("Failed to open cdb file")
			return
		}
	}
}

func indexFiles() {

	if *flagRebuild {
		books, err := prepareBooks("xliste.csv")
		if err != nil {
			fmt.Println("Failed to load book list csv file", err)
			return
		}

		for _, book := range books {
			indexBook(book)
		}
	} else {
		reindexAllFiles()
	}
}

func SavePageIndex() error {

	//var buf bytes.Buffer

	f, err := os.Create("index/page_index.gob")
	if err != nil {
		log.Println("Failed to create page_index.gob file")
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(pageIndex)

	if err != nil {
		log.Println(err)
		return err
	}

	f.Close()

	/*

		f, err = os.Create("index/page_pagestore.gob")
		if err != nil {
			log.Println("Failed to create page_pagestore.gob file")
			return err
		}
		defer f.Close()

		enc = gob.NewEncoder(f)
		err = enc.Encode(pageIndex.pageStore)

		if err != nil {
			log.Println(err)
			return err
		}


		f, err = os.Create("index/page_pagestore.gob")
		if err != nil {
			log.Println("Failed to create page_pagestore.gob file")
			return err
		}
		defer f.Close()

		enc = gob.NewEncoder(f)
		err = enc.Encode(pageIndex.pageStore)

		if err != nil {
			log.Println(err)
			return err
		}
	*/

	return nil
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
			payloadStore.reader.Close()
			os.Exit(0)
		}
	}()
}

var flagRebuild *bool
var flagBuildPayload *bool

func main() {

	flagRebuild = flag.Bool("rebuild", false, "rebuild index form scratch using csv file")
	flagBuildPayload = flag.Bool("payload", false, "rebuild payload cdb file form scratch")

	flag.Parse()

	fmt.Println(*flagRebuild)
	fmt.Println(*flagBuildPayload)

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
