package main

import (
	"io"
	"net/http"
)

func githubRequest(method, url string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, accessToken)
	req.Header.Set("Accept", acceptType)
	return req, nil
}
