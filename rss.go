package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Thumbnail struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type MediaGroup struct {
	Title       string    `xml:"title"`
	Thumbnail   Thumbnail `xml:"thumbnail"`
	Description string    `xml:"description"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

type Entry struct {
	ID         string     `xml:"id"`
	Title      string     `xml:"title"`
	Link       Link       `xml:"link"`
	Published  string     `xml:"published"`
	Updated    string     `xml:"updated"`
	MediaGroup MediaGroup `xml:"group"`
}

type RSS struct {
	ID    string  `xml:"id"`
	Title string  `xml:"title"`
	Entry []Entry `xml:"entry"`
}

func getRSSFeed(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return convertRSSStringToRSS(string(body))
	}

	return nil, fmt.Errorf("Returned invalid response for address %s. Response was %d", url, resp.StatusCode)
}

func convertRSSStringToRSS(file string) (*RSS, error) {
	var rss RSS
	if err := xml.Unmarshal([]byte(file), &rss); err != nil {
		return nil, err
	}

	return &rss, nil
}
