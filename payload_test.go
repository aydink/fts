package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/colinmarc/cdb"
)

func TestCratePagePayloadDatabase(t *testing.T) {

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

	db, err := writer.Freeze()
	if err != nil {
		log.Fatal(err)
	}

	v, err := db.Get([]byte("cd9c20b669efc8942a46cee198f44a7a-20"))
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(v)

	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	m := make(map[string][][4]int)

	if err := dec.Decode(&m); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m["TAARRUZ"])
}

func TestLoadPagePayload(t *testing.T) {
	fmt.Println(GetPagePayloadJSON(LoadPagePayload("cd9c20b669efc8942a46cee198f44a7a-20")))
}
