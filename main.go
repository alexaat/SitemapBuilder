package main

import (
	"SitemapBuilder/htmlLinkParser"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

var domain string
var urlSet = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`

func main() {

	url := flag.String("url", "", "The url to website in format http(s)://...")
	flag.Parse()
	if url == nil || strings.TrimSpace(*url) == "" {
		return
	}
	domain = *url

	var linksBD = []Url{{Loc: domain, checked: false}}

	var currentDbSize = len(linksBD)

	for {

		for _, v := range linksBD {

			if !v.checked {
				v.checked = true // Mark url as checked to prevent double-checking
				fmt.Println("Analizing", v.Loc)
				var content = getPageContent(v.Loc)            // Page content as stringmk
				var links = getLinks(content)                  // Get htmlLinkParser links
				var covertedLinks = convertLinksToModel(links) // Convert to Xml links
				linksBD = append(linksBD, covertedLinks...)    // Add fresh links to db
				linksBD = filterAndAdjust(linksBD)             // Filter Links
			}
		}

		if len(linksBD) <= currentDbSize {
			break
		}
		currentDbSize = len(linksBD)

	}

	var xmlString = string(toXml(linksBD))
	fmt.Println("Number of links: ", len(linksBD))
	fmt.Println(xmlString)

}

func toXml(links []Url) []byte {
	if xmlstring, err := xml.MarshalIndent(links, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + urlSet + string(xmlstring))
		return xmlstring
	}
	return []byte{}
}

func convertLinksToModel(links []htmlLinkParser.Link) []Url {
	var result = []Url{}
	for _, l := range links {
		var link = Url{
			Loc:     l.Href,
			checked: false,
		}
		result = append(result, link)
	}
	return result
}

func getPageContent(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("unable to open %v\n", url))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("unable to read response")
	}
	return string(body)
}

func getLinks(content string) []htmlLinkParser.Link {
	var reader = strings.NewReader(content)
	return htmlLinkParser.Parse(reader)
}

func filterAndAdjust(links []Url) []Url {
	var result = []Url{}
	for _, link := range links {
		var text = link.Loc
		//Check for emails
		res1, e := regexp.MatchString("mailto:[a-zA-Z0-9]*@.*", text)
		if res1 || e != nil {
			fmt.Println("email found: ", text)
			continue
		}

		//Check for links without domain and adjust trailing /
		if !strings.HasPrefix(text, "http://") && !strings.HasPrefix(text, "https://") {
			if !strings.HasSuffix(domain, "/") && !strings.HasPrefix(text, "/") {
				text = domain + "/" + text
			} else {
				text = domain + text
			}
		}

		//Remove trailing /
		text = strings.TrimSuffix(text, "/")

		// filter for dublicates and forien links
		var newUrl = Url{Loc: text}
		if !slices.Contains(result, newUrl) && strings.Contains(text, domain) {
			result = append(result, newUrl)
		}
	}

	return result
}
