package main

import (
	"fmt"
	"testing"
)

func TestBook(t *testing.T) {
	/*
		books := []Book{
			{Id: 1, Title: "Aşk tesadüfleri sever", Author: "Aşk Meşk", Type: "Roman", Genre: "Romantik", Year: 1900, NumPages: 120, Category: []string{"Yeni", "İndirim"}},
			{Id: 2, Title: "Aşk bunun neresinde, ey aşk", Author: "Ali Veli", Type: "Hikaye", Genre: "Deneme", Year: 2001, NumPages: 20, Category: []string{"Yeni", "Bedava"}},
			{Id: 3, Title: "Milletlerarası Özel Hukukta Kişilik Haklarının İnternet Yoluyla İhlalinde Sorumluluk", Author: "Esra Tekin", Type: "Kitap", Genre: "Akademik", Year: 2021, NumPages: 258, Category: []string{"Hukuk", "İnternet", "Özel Hukuk"}},
		}
	*/

	books, err := prepareBooks("xliste.csv")
	if err != nil {
		t.Log("Failed to load book list csv file")
		t.Fail()
	}

	analyzer := NewSimpleAnalyzer(NewSimpleTokenizer())

	turkishFilter := NewTurkishLowercaseFilter()
	turkishAccentFilter := NewTurkishAccentFilter()

	analyzer.AddTokenFilter(turkishFilter)
	analyzer.AddTokenFilter(turkishAccentFilter)

	index := NewBookIndex(analyzer)
	pageIndex := NewPageIndex(analyzer)

	allPages := make([]Page, 0)

	pageId := 0

	for k, book := range books {
		book.Id = k
		index.Add(&book)
		pages, err := getPages(book)

		if err != nil {
			fmt.Println(err)
		}

		for _, page := range pages {
			page.Id = pageId
			page.BookId = book.Id
			pageIndex.Add(&page)
			allPages = append(allPages, page)
			pageId++
		}
	}

	fmt.Printf("%+v\n", index)

	q := "doktrini"

	fmt.Printf("\nsearching for: %s\n", q)
	fmt.Println("----------------------------")
	fmt.Printf("%+v\n", index.Search(q))

	q = "ilerleme mihveri"

	fmt.Printf("\nsearching for: %s\n", q)
	fmt.Println("----------------------------")
	fmt.Printf("%+v\n", pageIndex.Search(q))

	hits := pageIndex.Search(q)
	for _, hit := range hits {

		fmt.Println(books[allPages[hit.docId].BookId].Title)
		fmt.Println("=============================================================================")
		fmt.Println(allPages[hit.docId].Content)
		fmt.Println("------------------------------------------------------------------------------")
	}
}
