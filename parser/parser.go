package parser

import (
	"errors"
	"golang.org/x/net/html"
	"net/http"
)

type Parser struct {
	url   string
	class string
}

func (p *Parser) Parse() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, p.url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return p.parseResponse(resp)
}

func (p *Parser) parseResponse(resp *http.Response) ([]string, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Returned status " + resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	articles := p.findArticles(doc)
	texts := make([]string, 0)
	for _, article := range articles {
		text := ""
		addLinks(article, &text)
		texts = append(texts, text)
	}
	return texts, nil
}

func (p *Parser) findArticles(doc *html.Node) []*html.Node {
	articles := make([]*html.Node, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "article" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "news-article" {
					articles = append(articles, n)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return articles
}

func addLinks(article *html.Node, text *string) {
	for c := article.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "a" {
			href := ""
			for _, attr := range c.Attr {
				if attr.Key == "href" {
					href = attr.Val
				}
			}
			*text += "<a href='" + href + "'>"
		}
		if c.Type == html.TextNode {
			*text += c.Data
		}
		addLinks(c, text)
		if c.Type == html.ElementNode && c.Data == "a" {
			*text += "</a>"
		}
	}
}

func NewParser(url, class string) *Parser {
	return &Parser{url, class}
}
