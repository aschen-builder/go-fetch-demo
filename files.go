package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

func writeHtml(file string, bytes []byte) error {
	f, err := os.Create(file)
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed writing to file: %v", file),
			err,
		)
	}

	defer f.Close()
	if _, err = f.Write(bytes); err != nil {
		return errors.Join(
			fmt.Errorf("failed writing to file: %v", file),
			err,
		)
	}

	f.Sync()
	return nil
}

func parseFilename(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "unknown", errors.Join(
			fmt.Errorf("failed to parse url: %v", uri),
			err,
		)
	}

	return u.Hostname(), nil
}
