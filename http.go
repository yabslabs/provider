package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func createGETRequest(url string, reader io.Reader) (*http.Request, error) {
	return defaultRequest(http.MethodGet, url, reader)
}

func createPOSTRequest(url string, body interface{}) (*http.Request, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return defaultRequest(http.MethodPost, url, bytes.NewBuffer(bodyJSON))
}

func sendRequest(client *http.Client, req *http.Request, responseObject interface{}) error {
	response, err := client.Do(req)
	defer func() {
		if err = response.Body.Close(); err != nil {
			log.Println("YABS--zdBN: unable to close body: ", err)
		}
	}()
	if response.Status != "200 OK" && !strings.Contains(response.Status, "201 Created") {
		return fmt.Errorf("request failed: %v", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseObject)
}

func defaultRequest(method, url string, reader io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, accessToken)
	req.Header.Set(acceptKey, githubAcceptType)
	return req, nil
}
