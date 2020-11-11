package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	for _, chann := range feeds {
		channelToGet := chann

		if channelToGet.ArchivalMode == ArchivalModeArchive {
			videosToGet, err := getVideoUpdatesForChannel(&channelToGet)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(channelToGet.Name)
			fmt.Printf("%d new videos available\n", len(*videosToGet))

			if len(*videosToGet) > 0 {
				command := getYoutubeDLCommandForVideoList(videosToGet)
				fmt.Println(command)
			}

			fmt.Println("")
		}
	}
}

const youtubeDLBaseCommand = "youtube-dl -f best -i --write-description --add-metadata --all-subs --sub-format \"srt\" --embed-subs "

func getYoutubeDLCommandForVideoList(list *[]Entry) string {
	var youtubeDlList []string
	for _, entry := range *list {
		// fmt.Printf("Title: %s, Link: %s \n\n", entry.Title, entry.Link.Href)
		youtubeDlList = append(youtubeDlList, "\""+entry.Link.Href+"\"")
	}

	downloadString := strings.Join(youtubeDlList, ", ")

	return youtubeDLBaseCommand + downloadString
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
