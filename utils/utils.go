package utils

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// YTCHTTPClient a HTTP client with timeout
type YTCHTTPClient interface {
	Get(url string) (*http.Response, []byte, error)
	// Post(url string, body []string, timeout time.Duration)
}

// DefaultHTTPTimeout provides the default HTTP request timeout value
var DefaultHTTPTimeout = 10 * time.Second

// HTTPClient is an implementation of the YTCHTTPClient interface for http requets
// with modifications specific to this project
type HTTPClient struct {
	ConnTimeout time.Duration
}

// Get is a simplified Get with a configurable timeout value.
// You can provide the default package timeout by using DefaultHTTPTimeout as the timeout
// value
func (ht *HTTPClient) Get(url string) (*http.Response, []byte, error) {
	tr := &http.Transport{
		IdleConnTimeout: ht.ConnTimeout,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
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

// OSCommandProvider provides the ability to run commands on the OS level
type OSCommandProvider interface {
	Run(string, ...string) (*[]byte, error)
}

// OSCommand implements OSCommandProvider to provide the ability to run commands on the OS
type OSCommand struct{}

// Run runs a command on the OS
func (osc *OSCommand) Run(name string, arg ...string) (*[]byte, error) {
	out, err := exec.Command(name, arg...).Output()
	return &out, err
}

// DirReaderProvider provides the ability to read file directories on disk
type DirReaderProvider interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
}

// DirReader implements DirReaderProvider to provide the ability to read file directories on disk
type DirReader struct{}

// ReadDir reads the contents of a directory
func (dr *DirReader) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}
