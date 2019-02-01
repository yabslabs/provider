package util

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

const (
	statusOK      = "200 OK"
	statusCreated = "201 Created"
	acceptKey     = "Accept"
)

type requestGenerator func(method, url string, reader io.Reader) (*http.Request, error)

func CreateGETRequest(url string, reader io.Reader, requestFunc requestGenerator) (*http.Request, error) {
	return requestFunc(http.MethodGet, url, reader)
}

func CreatePOSTRequest(url string, body interface{}, requestFunc requestGenerator) (*http.Request, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return requestFunc(http.MethodPost, url, bytes.NewBuffer(bodyJSON))
}

func DoRequest(client *http.Client, req *http.Request) ([]byte, error) {
	response, err := client.Do(req)
	defer func() {
		if err = response.Body.Close(); err != nil {
			log.Println("YABS--zdBN: unable to close body: ", err)
		}
	}()
	if response.Status != statusOK && !strings.Contains(response.Status, statusCreated) {
		return nil, fmt.Errorf("request failed: %v", response.Status)
	}
	return ioutil.ReadAll(response.Body)
}

func DoRequestWithUnmarshal(client *http.Client, req *http.Request, responseObject interface{}) error {
	body, err := DoRequest(client, req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseObject)
}
