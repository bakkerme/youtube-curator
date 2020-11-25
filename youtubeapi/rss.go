package youtubeapi

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// RSSThumbnail represents the thumbnail image of a video
type RSSThumbnail struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// RSSMediaGroup contains Video metadata, like the Title, Description and RSSThumbnail
type RSSMediaGroup struct {
	Title       string       `xml:"title"`
	Thumbnail   RSSThumbnail `xml:"thumbnail"`
	Description string       `xml:"description"`
}

// RSSLink is the link to the video, in the youtube web interface
type RSSLink struct {
	Href string `xml:"href,attr"`
}

// RSSVideoEntry contains information about a Video from the RSS feed
type RSSVideoEntry struct {
	ID         string        `xml:"videoId"`
	Title      string        `xml:"title"`
	Link       RSSLink       `xml:"link"`
	Published  string        `xml:"published"`
	Updated    string        `xml:"updated"`
	MediaGroup RSSMediaGroup `xml:"group"`
}

// RSS is a struct designed to contain video data from a youtube channel RSS feed
type RSS struct {
	ID         string          `xml:"id"`
	Title      string          `xml:"title"`
	VideoEntry []RSSVideoEntry `xml:"entry"`
}

func getRSSFeed(url string) (*RSS, error) {
	tr := &http.Transport{
		IdleConnTimeout: 10 * time.Second,
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