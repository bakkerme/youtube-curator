package main

import (
	"encoding/xml"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
)

// Thumbnail represents the thumbnail image of a video
type Thumbnail struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// MediaGroup contains Video metadata, like the Title, Description and Thumbnail
type MediaGroup struct {
	Title       string    `xml:"title"`
	Thumbnail   Thumbnail `xml:"thumbnail"`
	Description string    `xml:"description"`
}

// Link is the link to the video, in the youtube web interface
type Link struct {
	Href string `xml:"href,attr"`
}

// VideoEntry contains information about a Video from the RSS feed
type VideoEntry struct {
	ID         string     `xml:"videoId"`
	Title      string     `xml:"title"`
	Link       Link       `xml:"link"`
	Published  string     `xml:"published"`
	Updated    string     `xml:"updated"`
	MediaGroup MediaGroup `xml:"group"`
}

// RSS is a struct designed to contain video data from a youtube channel RSS feed
type RSS struct {
	ID    string  `xml:"id"`
	Title string  `xml:"title"`
	VideoEntry []VideoEntry `xml:"entry"`
}

func getRSSFeed(url string) (*RSS, error) {
	tr := &http.Transport{
		IdleConnTimeout:    10 * time.Second,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
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
