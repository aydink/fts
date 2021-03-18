package fts

import (
	"fmt"
	"testing"
)

func TestBook(t *testing.T) {
	books := []Book{
		{Id: 1, Title: "Aşk tesadüfleri sever", Author: "Aşk Meşk", Type: "Roman", Genre: "Romantik", Year: 1900, NumPages: 120, Category: []string{"Yeni", "İndirim"}},
		{Id: 2, Title: "Aşk bunun neresinde, ey aşk", Author: "Ali Veli", Type: "Hikaye", Genre: "Deneme", Year: 2001, NumPages: 20, Category: []string{"Yeni", "Bedava"}},
		{Id: 3, Title: "Milletlerarası Özel Hukukta Kişilik Haklarının İnternet Yoluyla İhlalinde Sorumluluk", Author: "Esra Tekin", Type: "Kitap", Genre: "Akademik", Year: 2021, NumPages: 258, Category: []string{"Hukuk", "İnternet", "Özel Hukuk"}},
	}

	analyzer := NewSimpleAnalyzer(NewSimpleTokenizer())

	turkishFilter := NewTurkishLowercaseFilter()
	turkishAccentFilter := NewTurkishAccentFilter()

	analyzer.AddTokenFilter(turkishFilter)
	analyzer.AddTokenFilter(turkishAccentFilter)

	index := NewBookIndex(analyzer)

	for _, book := range books {
		index.Add(&book)
	}

	fmt.Printf("%+v\n", index)

	q := " Aşk"

	fmt.Printf("\nsearching for: %s\n", q)
	fmt.Println("----------------------------")
	fmt.Printf("%+v\n", index.Search(q))

}
