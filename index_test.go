package fts

import (
	"fmt"
	"testing"
)

func Fields(t *testing.T) {

	analyzer := NewSimpleAnalyzer(NewSimpleTokenizer())

	turkishFilter := NewTurkishLowercaseFilter()
	turkishAccentFilter := NewTurkishAccentFilter()

	analyzer.AddTokenFilter(turkishFilter)
	analyzer.AddTokenFilter(turkishAccentFilter)

	index := New("index", analyzer)

	textTitle := "Merhaba dünya, bu bir title dünya"
	textContent := "TOKEN MOKEN Babası (sıfır-bir) otomobil hediye etti! Ekspertize götürünce hayatının şokunu yaşadı. ĞÜŞİÖÇ ÂÎÛ âîû"

	title := NewTextField("title", textTitle, INDEXED, STORED, ANALYZED)
	content := NewTextField("content", textContent, INDEXED, STORED, ANALYZED)

	doc := NewDocument()
	doc.Add(title)
	doc.Add(content)
	doc.AddCatergory("cat1")
	doc.AddCatergory("cat2")

	index.Add(doc)

	textTitle = "İkinci  hamlesi kalkınma, dünya dünya yalan dünya"
	textContent = "sana bir gül hediye etsem nasıl olur, nasıl desem, nasıl anlatsam otomobil mi alsam acaba"

	title = NewTextField("title", textTitle, INDEXED, STORED, ANALYZED)
	content = NewTextField("content", textContent, INDEXED, STORED, ANALYZED)

	doc = NewDocument()
	doc.Add(title)
	doc.Add(content)
	doc.AddCatergory("cat1")
	doc.AddCatergory("cat2")

	index.Add(doc)

	textTitle = "Üçüncü kalkınma hamlesi"
	textContent = "Otomobil ile bir tur atalım seninle otomobil ile"

	title = NewTextField("title", textTitle, INDEXED, STORED, ANALYZED)
	content = NewTextField("content", textContent, INDEXED, STORED, ANALYZED)

	doc = NewDocument()
	doc.Add(title)
	doc.Add(content)
	doc.AddCatergory("cat1")
	doc.AddCatergory("cat2")
	index.Add(doc)

	index.Commit()

	fmt.Printf("%+v\n\n", index)
	posting1, _ := index.ReadPostings("title:kalkınma")
	posting2, _ := index.ReadPostings("title:hamlesi")
	posting := Intersection(posting1.Postings, posting2.Postings)
	fmt.Printf("%+v\n", posting)

}
