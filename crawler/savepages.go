package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// https://stackoverflow.com/questions/34995071/how-to-efficiently-store-html-response-to-a-file-in-golang
// https://godoc.org/golang.org/x/net/html#Parse
// body, err := ioutil.ReadFile("body.html")
// err = ioutil.WriteFile(pgTitle+".html", body, 0644)
func main() {
	resp, err := http.Get("https://en.wikipedia.org/wiki/Java")
	body, err := ioutil.ReadAll(resp.Body)
	doc, err := html.Parse(strings.NewReader(string(body[:])))
	if err != nil {
		panic(err)
	}

	/* Get Page Title */
	var pgTitle string
	var getPageTitle func(*html.Node)
	getPageTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			pgTitle = n.FirstChild.Data
			pgTitle = pgTitle[:len(pgTitle)-len(" - Wikipedia")]
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getPageTitle(c)
		}
	}
	getPageTitle(doc)

	/* Parse Links */
	var str string
	var getLinks func(*html.Node)
	getLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if a.Val[:6] == "/wiki/" && !strings.Contains(a.Val, ":") {
						str = str + a.Val + "\n"
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getLinks(c)
		}
	}
	getLinks(doc)

	err = ioutil.WriteFile(pgTitle+".link", []byte(str), 0644)
	if err != nil {
		panic(err)
	}
}
