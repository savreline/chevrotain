package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// https://stackoverflow.com/questions/34995071/how-to-efficiently-store-html-response-to-a-file-in-golang
// https://godoc.org/golang.org/x/net/html#Parse
// https://stackoverflow.com/questions/2818852/is-there-a-queue-implementation
// https://stackoverflow.com/questions/34018908/golang-why-dont-we-have-a-set-datastructure
// body, err := ioutil.ReadFile("body.html")
// err = ioutil.WriteFile(pgTitle+".html", body, 0644)
func main() {
	var queue []string
	var pgTitle string
	var strOfLinks string
	var setOfLinks map[string]bool
	var lastLink string
	startPage := os.Args[1]
	lastPage := startPage
	curPage := startPage
	maxPerPage, _ := strconv.Atoi(os.Args[2])
	maxDepth, _ := strconv.Atoi(os.Args[3])
	queue = append(queue, startPage)
	os.Mkdir(startPage, 0644)

	/* Parse Title */
	var getPageTitle func(*html.Node)
	getPageTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			pgTitle = n.FirstChild.Data
			pgTitle = pgTitle[:len(pgTitle)-len(" - Wikipedia")]
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			getPageTitle(c)
		}
	}

	/* Parse Links */
	var getLinks func(*html.Node)
	getLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if len(a.Val) < 6 {
						break
					}
					linkHead := a.Val[:6]
					linkTail := a.Val[6:]
					if linkHead == "/wiki/" && !strings.Contains(a.Val, ":") &&
						linkTail != "Main_Page" && linkTail != curPage &&
						!setOfLinks[linkTail] {
						strOfLinks = strOfLinks + linkTail + "\n"
						queue = append(queue, linkTail)
						setOfLinks[linkTail] = true
						if len(setOfLinks) >= maxPerPage {
							lastLink = linkTail
							return
						}
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if len(setOfLinks) >= maxPerPage {
				return
			}
			getLinks(c)
		}
	}

	/* BFS */
	i := -1
	for len(queue) > 0 && i < maxDepth {
		fmt.Println(queue[0], ":", lastPage, ":", i)
		time.Sleep(1 * time.Second)
		if queue[0] == lastPage {
			i++
		}
		curPage = queue[0]

		/* Download the page */
		resp, err := http.Get("https://en.wikipedia.org/wiki/" + curPage)
		body, err := ioutil.ReadAll(resp.Body)
		doc, err := html.Parse(strings.NewReader(string(body[:])))
		if err != nil {
			panic(err)
		}

		/* Flush str and set */
		strOfLinks = ""
		setOfLinks = make(map[string]bool)

		/* Pop off queue, process links (getLinks enqueues links) */
		queue = queue[1:]
		getPageTitle(doc)
		getLinks(doc)

		/* Save last page */
		if curPage == lastPage {
			lastPage = lastLink
		}

		/* Write current link file */
		err = ioutil.WriteFile(startPage+"/"+pgTitle+".link", []byte(strOfLinks), 0644)
		if err != nil {
			panic(err)
		}
	}
}
