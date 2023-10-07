package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
	"golang.org/x/net/html"
)

var testDocument = `<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<h1>Test Page</h1>
<div>
<img src="https://example.com/image.png" />
</div>
<div>
<a href="https://example.com/page1">Page 1</a>
<a href="https://example.com/page2">Page 2</a>
</div>
</body>
</html>`

func createTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testDocument))
	}))
}

func Test_fetchUrl(t *testing.T) {
	ts := createTestServer(t)
	defer ts.Close()
	body, err := fetchUrl(ts.URL)
	if err != nil {
		t.Fatalf("fetchUrl: %v", err)
	}
	if len(body) == 0 {
		t.Fatalf("fetchUrl: empty body")
	}
}

func Test_parseFilename(t *testing.T) {
	body, err := parseFilename("https://example.com")
	if err != nil {
		t.Fatalf("parseFilename: %v", err)
	}
	if body != "example.com" {
		t.Fatalf("parseFilename: expected 'example.com', got '%v'", body)
	}
}

func Test_writeHtml(t *testing.T) {
	defer t.Cleanup(func() {
		if err := os.Remove("test.html"); err != nil {
			t.Logf("failed to remove test.html: %v", err)
		}
	})
	ts := createTestServer(t)
	defer ts.Close()
	body, err := fetchUrl(ts.URL)
	if err != nil {
		t.Fatalf("fetchUrl: %v", err)
	}
	if err = writeHtml("test.html", body); err != nil {
		t.Fatalf("writeHtml: %v", err)
	}
}

func Test_queryAttributeValue(t *testing.T) {
	ts := createTestServer(t)
	defer ts.Close()
	body, err := fetchUrl(ts.URL)
	if err != nil {
		t.Fatalf("fetchUrl: %v", err)
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		t.Fatalf("html.Parse: %v", err)
	}
	count := len(NodeTarget{Query: "a", Document: doc, Attr: "href"}.QueryAll())
	if count != 2 {
		t.Fatalf("countNodes: expected 2, got %v", count)
	}
}

func Test_gatherMetadata(t *testing.T) {
	ts := createTestServer(t)
	defer ts.Close()
	body, err := fetchUrl(ts.URL)
	if err != nil {
		t.Fatalf("fetchUrl: %v", err)
	}
	metadata := gatherMetadata(ts.URL, body)

	if len(metadata.Images) != 1 && len(metadata.Links) != 2 {
		t.Fatalf("gatherMetadata: expected %v, got %v", len(metadata.Images)+len(metadata.Links), 3)
	}
}

// This test could be much cleaner, but for the
// initial green-light tests it suffices
func Test_fetchAction(t *testing.T) {
	defer t.Cleanup(func() {
		if err := os.Remove("google.com.html"); err != nil {
			t.Logf("failed to remove google.com.html: %v", err)
		}
	})
	ts := createTestServer(t)
	defer ts.Close()
	app := cli.NewApp()
	app.Commands = []*cli.Command{
		{
			Name: "fetch",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "metadata",
					Usage: "print metadata about each fetched URL",
				},
				&cli.BoolFlag{
					Name:  "mirror",
					Usage: "start a local mirror to serve retrieved files",
				},
			},
			Action: func(ctx *cli.Context) error {
				fetchAction(ctx)
				enableMirror(ctx.Bool("mirror"))
				return nil
			},
		},
	}
	err := app.Run([]string{"test", "fetch", "https://google.com"})
	if err != nil {
		t.Fatalf("fetchAction: %v", err)
	}
}
