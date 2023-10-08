package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Keeping the slices flat for sake of time/demo
type HtmlMetadata struct {
	Anchors   []string
	Images    []string
	Links     []string
	LastFetch string `example:"Mon Jan 2 15:04:05 UTC 2006"`
	Url       string
}

type NodeTarget struct {
	Attr     string `example:"href"`
	Query    string `example:"a"`
	Document *html.Node
}

func (t NodeTarget) QueryAll() []string {
	return goquery.NewDocumentFromNode(t.Document).Find(t.Query).Map(func(i int, s *goquery.Selection) string {
		return s.AttrOr(t.Attr, "")
	})
}

func gatherMetadata(uri string, body []byte) HtmlMetadata {
	Document, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}
	return HtmlMetadata{
		Anchors:   NodeTarget{Query: "a", Document: Document, Attr: "href"}.QueryAll(),
		Images:    NodeTarget{Query: "img", Document: Document, Attr: "src"}.QueryAll(),
		Links:     NodeTarget{Query: "link", Document: Document, Attr: "href"}.QueryAll(),
		LastFetch: time.Now().UTC().Format(time.UnixDate),
		Url:       uri,
	}
}

func printMetadata(htmlMeta HtmlMetadata) {
	fmt.Printf("site: %s\n", htmlMeta.Url)
	fmt.Printf("num_links: %d\n", len(htmlMeta.Anchors)+len(htmlMeta.Links))
	fmt.Printf("images: %d\n", len(htmlMeta.Images))
	fmt.Printf("last_fetch: %s\n", htmlMeta.LastFetch)
	fmt.Println()
}
