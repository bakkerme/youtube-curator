package utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

// DefaultHTTPTimeout provides the default HTTP request timeout value
var DefaultHTTPTimeout = 10 * time.Second

// HTTPGet is a simplified Get with a 10s timeout
func HTTPGet(url string, timeout int) (*Response, []byte, error) {
	tr := &http.Transport{
		IdleConnTimeout: timeout,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		// return nil, err
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return resp, body, err
}
