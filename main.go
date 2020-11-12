package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	ch := make(chan YTChannelEntryResult)
	for _, chann := range feeds {
		// channelToGet := chann
		go func (channelToGet YTChannel) {
			// done := make(chan YTChannelEntryResult)
			// go getVideoUpdatesForYTChannelWithChannel(done, &channelToGet)
			videosToGet, err := getVideoUpdatesForYTChannel(&channelToGet)
			// channelResult := <-done

			// videosToGet := channelResult.Entries
			// err := channelResult.Error

			if err != nil {
				log.Fatal(err)
				ch <- YTChannelEntryResult{nil, err}
			}

			fmt.Println(channelToGet.Name)
			fmt.Printf("%d new videos available\n", len(*videosToGet))

			if len(*videosToGet) > 0 {
				if channelToGet.ArchivalMode == ArchivalModeCurated {
					command := getYoutubeDLCommandForVideoList(videosToGet)
					fmt.Println(command)
				} else if channelToGet.ArchivalMode == ArchivalModeArchive {
					fmt.Println(channelToGet.ChannelURL)
				}
			}

			fmt.Println("")

			ch <- YTChannelEntryResult{videosToGet, err}
		}(chann)
	}

	for _, chann := range feeds {
		fmt.Println(chann.Name)
		<-ch
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

type YTChannelEntryResult struct {
	Entries *[]Entry
	Error   error
}

func getVideoUpdatesForYTChannelWithChannel(c chan YTChannelEntryResult, ytc *YTChannel) {
	entries, err := getVideoUpdatesForYTChannel(ytc)
	c <- YTChannelEntryResult{entries, err}
}

func getVideoUpdatesForYTChannel(ytc *YTChannel) (*[]Entry, error) {
	rssFeed := ytc.RSSURL
	rss, err := getRSSFeed(rssFeed)

	if err != nil {
		return nil, fmt.Errorf("Loading RSS feed xml failed for %s channel. URL is %s. Error: %s", ytc.Name, ytc.RSSURL, err)
	}

	entries := rss.Entry

	videos, err := getLocalVideosByYTChannel(ytc)
	if err != nil {
		return nil, fmt.Errorf("Could not get videos off disk for ytc %s, error %s", ytc.Name, err)
	}

	entriesToDownload := getEntriesNotInVideoList(&entries, videos)

	return entriesToDownload, nil
}
