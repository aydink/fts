package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type HitResult struct {
	Book   Book
	Page   Page
	HlText string
}

// getFullFilterName return full name of the filter
// eg. "year": "Yıl", "genre": "Türü"
func getFullFilterName(key string) string {

	filterFullNames := map[string]string{
		"year":     "Basım yılı",
		"genre":    "Türü",
		"category": "Kategori",
	}

	if value, found := filterFullNames[key]; found {
		return value
	}
	return key
}

func getFilters(v url.Values) [][3]string {
	filterNames := []string{"genre", "type", "category"}

	filters := make([][3]string, 0)

	for _, name := range filterNames {
		if v.Get(name) != "" {
			filters = append(filters, [3]string{name, getFullFilterName(name), v.Get(name)})
		}
	}
	return filters
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		fmt.Fprintf(w, "Hata: %s!", err)
	}
	t.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	v := r.URL.Query()
	filters := getFilters(v)

	//t, err := template.ParseFiles("templates/search.html")
	//t := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/search.html", "templates/partial_facet.html", "templates/partial_pagination.html", "templates/partial_definition.html"))
	t := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	q := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	//searchType := r.URL.Query().Get("w")

	start := r.URL.Query().Get("start")
	startInt, err := strconv.Atoi(start)
	//	fmt.Println("start:", startInt)

	if err != nil {
		//fmt.Println("error parsing 'start' parameter")
		startInt = 0
	}

	templateName := "search"

	hits := pageIndex.Search(q)

	if len(category) > 0 {
		hits = pageIndex.facetFilterCategory(hits, category)
	}

	data := make(map[string]interface{})
	data["q"] = q
	data["categoryFacet"] = pageIndex.getFacetCounts(hits)

	totalHits := len(hits)
	data["TotalHits"] = totalHits

	data["pages"] = Paginate(startInt, 10, len(hits))
	data["filters"] = filters

	if startInt < totalHits {
		if (startInt + 10) < totalHits {
			hits = hits[startInt : startInt+10]
		} else {
			hits = hits[startInt:]
		}
	}
	/*
		if len(hits) > 10 {
			hits = hits[0:10]
		}
	*/
	hitResults := make([]HitResult, 0)

	for _, hit := range hits {

		result := HitResult{}
		result.Page = pageIndex.GetPage(int(hit.docId))
		result.Book = bookIndex.bookStore[result.Page.BookId]
		result.HlText = HighlightText(result.Page.Content, q)

		hitResults = append(hitResults, result)

		//fmt.Println(index.bookId[hit.docId], books[pages[hit.docId].BookId-1].Title)
		//fmt.Println("---", pages[hit.docId].Content)
	}

	data["hits"] = hitResults

	// do not search book titles if a filter applied
	if len(category) == 0 && startInt == 0 {

		books := make([]Book, 0)

		postings := bookIndex.Search(q)
		for _, posting := range postings {
			book := bookIndex.bookStore[posting.docId]
			book.Title = bookIndex.HighlightTitle(book.Title, q)
			books = append(books, book)
		}

		data["titles"] = books
	}
	/*
		if searchType == "title" {
			//data = titleQuery(keywords, startInt, getFilters(r.URL.Path))
			data = titleQuery(keywords, startInt, filters)
			templateName = "title"
		} else {
			//data = query(keywords, startInt, getFilters(r.URL.Path))
			data = query(keywords, startInt, filters)

			//here we will show book titles that mathes most of (%60) of keywords
			//it is smiliar to providing wikipedia definition for given keywords
			//only if we are displaying the first page of search result. Otherwise
			//it is meaningless to display best mathing book titles.
			if startInt == 0 {
				titles := titleQuerySimple(keywords)
				data["titles"] = titles["books"]
			}
		}

		// show dictionary definion on only first page
		if startInt == 0 {
			data["definition"], data["hasDefinition"] = queryDictionary(keywords)
		}
	*/
	err = t.ExecuteTemplate(w, templateName, data)
	if err != nil {
		fmt.Println(err)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	q := r.URL.Query().Get("q")
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		pageInt = 0
	}

	curPage := pageIndex.GetPage(pageInt)
	curBook := bookIndex.bookStore[curPage.BookId]
	hash := curBook.Hash

	image := hash + "-" + strconv.Itoa(curPage.PageNumber)

	createImage(image)

	data := make(map[string]interface{})
	data["q"] = q
	data["title"] = curBook.Title
	data["image"] = image
	data["content"] = curPage.Content
	data["hash"] = hash
	data["pageNumber"] = curPage.PageNumber
	data["numPages"] = curBook.NumPages
	data["numPages"] = pageIndex.NumDocs
	data["page"] = pageInt

	//fmt.Printf("%+v", data)
	t.ExecuteTemplate(w, "document", data)
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	q := r.URL.Query().Get("q")

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, GetTokenPositions(page, q))
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("page")
	parts := strings.Split(query, "-")
	hash := parts[0]
	page := parts[1]

	createImage(query)

	http.ServeFile(w, r, "static/images/"+hash+"-"+page+".png")
}

func reindexHandler(w http.ResponseWriter, r *http.Request) {
	go reindexAllFiles()
	fmt.Fprint(w, "Reindexing all pdf files")
}

func createImage(query string) {

	parts := strings.Split(query, "-")
	hash := parts[0]
	page := parts[1]

	//fmt.Println("hash:", hash, "page:", page, "file:", fileMap[hash])

	if _, err := os.Stat("static/images/" + hash + "-" + page + ".png"); os.IsNotExist(err) {
		_, err := exec.Command("poppler/pdftocairo", "-png", "-singlefile", "-f", page, "-l", page, "books/"+hash+".pdf", "static/images/"+hash+"-"+page).Output()
		if err != nil {
			//log.Fatalln(err)
			log.Println(err)
		}
	} else {
		//fmt.Println("-----------------------", "using cashed image")
	}
}

// PDF file handler
// send pdf file and sets a proper title
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("book")
	if len(hash) > 32 {
		hash = hash[:32]
	}

	if len(hash) < 32 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Geçersiz bir istekte bulundunuz.")
		log.Printf("download pdf, invalid hash value:%s", hash)
		return
	}

	// check if user wants to download file
	force := r.URL.Query().Get("force")

	file, err := os.Open("books/" + hash + ".pdf")
	defer file.Close()
	if err != nil {
		log.Printf("failed to serve pdf file:%s", "books/"+hash+".pdf")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Geçersiz bir istekte bulundunuz.")
		return
	}

	book := bookIndex.GetBook(hash)
	name := book.Title
	name = strings.TrimSpace(name)

	// if there is an explicit url prameter "force=true" then force browser to download not try to display the pdf file
	if force == "true" {
		w.Header().Set("Content-Disposition", "attachment; filename="+name+".pdf")
	}

	w.Header().Set("Content-Type", "application/pdf")
	io.Copy(w, file)
}

// send token statistics
func tokenStatHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Disposition", "attachment; filename=token_stats.txt")
	w.Header().Set("Content-Type", "text/plain")

	for _, v := range pageIndex.tokenStats() {
		fmt.Fprintf(w, "%s\t%d\n", v.Name, v.Count)
	}
}
