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
	Text   template.HTML
	NewRow bool
}

func homeHandler(c http.ResponseWriter, r *http.Request) {
	p := parser.NewParser(URL, CLASS)
	res, err := p.Parse()
	if err != nil {
		log.Fatalln(err)
	}
	articles := make([]*Article, 0)
	for i, s := range res {
		articles = append(articles, &Article{template.HTML(s), i%6 == 0})
	}

	var homeTempl = template.Must(template.ParseFiles("web/main.html"))

	data := struct {
		Articles []*Article
	}{articles}
	homeTempl.Execute(c, data)
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./web/img/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css/"))))
	http.Handle("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir("./web/font/"))))

	s := &http.Server{
		Addr:           ADDR,
		Handler:        nil,
		ReadTimeout:    1000 * time.Second,
		WriteTimeout:   1000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
