package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func fetchAction(ctx *cli.Context) {
	if !ctx.Args().Present() {
		log.Fatal(errors.New("no URLs provided"))
	}

	uris := ctx.Args().Slice()
	ch := make(chan int, len(uris))
	for _, uri := range uris {
		if _, err := validateUrl(uri); err != nil {
			fmt.Println(err)
			break
		}

		fileName, err := parseFilename(uri)
		if err != nil {
			fmt.Println(err)
		}

		body, err := fetchUrl(uri)
		if err != nil {
			fmt.Println(err)
			break
		}

		if err = writeHtml(fileName+".html", body); err != nil {
			fmt.Println(err)
		}

		if ctx.Bool("metadata") {
			metadata := gatherMetadata(uri, body)
			printMetadata(metadata)
		}
		ch <- 1
	}

	close(ch)
}

func main() {
	app := &cli.App{
		Name:        "fetch",
		HideVersion: true,
		Usage:       "fetch a list of URLs and save them to disk",
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
		UsageText: "fetch [options] [urls]",
		Action: func(ctx *cli.Context) error {
			fetchAction(ctx)
			enableMirror(ctx.Bool("mirror"))
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
