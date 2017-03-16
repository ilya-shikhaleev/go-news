package main

import (
	"github.com/ilya-shikhaleev/go-news/parser"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	URL   string = "http://4gophers.com/news"
	CLASS string = "news-article"
	ADDR  string = ":12322"
)

type Article struct {
	Text template.HTML
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	p := parser.NewParser(URL, CLASS)
	res, err := p.Parse()
	if err != nil {
		log.Fatalln(err)
	}
	articles := make([]*Article, 0)
	for _, s := range res {
		articles = append(articles, &Article{template.HTML(s)})
	}

	var homeTpl = template.Must(template.ParseFiles("web/main.html"))

	data := struct {
		Articles []*Article
	}{articles}
	homeTpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./web/img/"))))
	http.Handle("/js/", http.FileServer(http.Dir("./web/")))
	http.Handle("/css/", http.FileServer(http.Dir("./web/")))

	s := &http.Server{
		Addr:           ADDR,
		Handler:        nil,
		ReadTimeout:    90 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
