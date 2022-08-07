package htmlLinkParser

import (
	"io"
	"log"

	"golang.org/x/net/html"
)

func Parse(reader io.Reader) []Link {

	var links []Link

	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, Link{
						Href: a.Val,
						Text: getText(n),
					})
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links
}

func getText(n *html.Node) string {

	var text = ""

	if n.Type == html.TextNode {
		text = text + n.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}

	return text
}
