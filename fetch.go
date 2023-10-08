package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func validateUrl(uri string) (*url.URL, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to parse url: %v", uri),
			err,
		)
	}
	return u, nil
}

func fetchUrl(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil || res.StatusCode >= 400 {
		return nil, errors.Join(
			fmt.Errorf("failed to fetch url: %v", uri),
			err,
		)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to read response body: %v", uri),
			err,
		)
	}
	return body, nil
}
