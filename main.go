package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	lunduke := feeds["BryanLunduke"]
	videosToGet, err := getVideoUpdatesForChannel(&lunduke)

	if err != nil {
		log.Fatal(err)
	}

	var youtubeDlList []string
	for _, entry := range *videosToGet {
		// fmt.Printf("Title: %s, Link: %s \n\n", entry.Title, entry.Link.Href)
		youtubeDlList = append(youtubeDlList, "\""+entry.Link.Href+"\"")
	}

	downloadString := strings.Join(youtubeDlList, ", ")

	fmt.Printf("youtube-dl -f best -i --write-description %s", downloadString)
}

func getVideoUpdatesForChannel(channel *Channel) (*[]Entry, error) {
	rssFeed := channel.RSSURL
	rss, err := getRSSFeed(rssFeed)

	if err != nil {
		return nil, fmt.Errorf("Loading RSS feed xml failed for %s channel. URL is %s. Error: %s", channel.Name, channel.RSSURL, err)
	}

	entries := rss.Entry

	videos, err := getLocalVideosByChannel(channel)
	if err != nil {
		return nil, fmt.Errorf("Could not get videos off disk for channel %s, error %s", channel.Name, err)
	}

	entriesToDownload := getEntriesNotInVideoList(&entries, videos)

	return entriesToDownload, nil
}
