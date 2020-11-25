package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/channel"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"hyperfocus.systems/youtube-curator-server/youtubedl"
)

type ytChannelEntryResult struct {
	Channel *collection.YTChannel
	Entries *[]youtubeapi.VideoEntry
	Error   error
}

func main() {
	ch := make(chan ytChannelEntryResult)
	for _, chann := range collection.Feeds {
		go func(channelToGet collection.YTChannel) {
			videosToGet, err := getVideoUpdatesForYTChannel(&channelToGet)
			ch <- ytChannelEntryResult{&channelToGet, videosToGet, err}
		}(chann)
	}

	for range collection.Feeds {
		result := <-ch
		videosToGet := result.Entries
		channelToGet := result.Channel
		err := result.Error

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(channelToGet.Name)
		fmt.Printf("%d new videos available\n", len(*videosToGet))

		if len(*videosToGet) > 0 {
			command, err := youtubedl.GetCommandForArchivalType(channelToGet, videosToGet)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(command)
		}

		fmt.Println("")
	}
}

func getVideoUpdatesForYTChannel(ytc *YTChannel) (*[]youtubeapi.VideoEntry, error) {
	rssFeed := ytc.RSSURL
	rss, err := youtubeapi.getRSSFeed(rssFeed)

	if err != nil {
		return nil, fmt.Errorf("Loading RSS feed xml failed for %s channel. URL is %s. Error: %s", ytc.Name, ytc.RSSURL, err)
	}

	entries := rss.youtubeapi.VideoEntry

	videos, err := channel.GetLocalVideosByYTChannel(ytc)
	if err != nil {
		return nil, fmt.Errorf("Could not get videos off disk for ytc %s, error %s", ytc.Name, err)
	}

	entriesToDownload := collection.GetEntriesNotInVideoList(&entries, videos)

	return entriesToDownload, nil
}
