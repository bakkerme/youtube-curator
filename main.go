package main

import (
	"fmt"
)

type ytChannelEntryResult struct {
	Entries *[]VideoEntry
	Error   error
}

func main() {
	ch := make(chan ytChannelEntryResult)
	for _, chann := range feeds {
		go func(channelToGet YTChannel) {
			videosToGet, err := getVideoUpdatesForYTChannel(&channelToGet)
			ch <- ytChannelEntryResult{videosToGet, err}
		}(chann)
	}

	for _, channelToGet := range feeds {
		result := <-ch
		videosToGet := result.Entries
		err := result.Error

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(channelToGet.Name)
		fmt.Printf("%d new videos available\n", len(*videosToGet))

		if len(*videosToGet) > 0 {
			command, err := getCommandForArchivalType(&channelToGet, videosToGet)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(command)
		}

		fmt.Println("")
	}
}

func getVideoUpdatesForYTChannel(ytc *YTChannel) (*[]VideoEntry, error) {
	rssFeed := ytc.RSSURL
	rss, err := getRSSFeed(rssFeed)

	if err != nil {
		return nil, fmt.Errorf("Loading RSS feed xml failed for %s channel. URL is %s. Error: %s", ytc.Name, ytc.RSSURL, err)
	}

	entries := rss.VideoEntry

	videos, err := getLocalVideosByYTChannel(ytc)
	if err != nil {
		return nil, fmt.Errorf("Could not get videos off disk for ytc %s, error %s", ytc.Name, err)
	}

	entriesToDownload := getEntriesNotInVideoList(&entries, videos)

	return entriesToDownload, nil
}
