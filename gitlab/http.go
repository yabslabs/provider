package main

import (
	"io"
	"net/http"
)

func gitlabRequest(method, url string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Private-Token", accessToken)
	return req, nil
}
