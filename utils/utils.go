package utils

import (
	"io/ioutil"
	"net/http"
	"time"
)

// DefaultHTTPTimeout provides the default HTTP request timeout value
var DefaultHTTPTimeout = 10 * time.Second

// HTTPGet is a simplified Get with a configurable timeout value.
// You can provide the default package timeout by using DefaultHTTPTimeout as the timeout
// value
func HTTPGet(url string, timeout time.Duration) (*http.Response, []byte, error) {
	tr := &http.Transport{
		IdleConnTimeout: timeout,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, err
}
