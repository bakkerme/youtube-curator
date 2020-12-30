package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"hyperfocus.systems/youtube-curator-server/youtubedl"
)

type ytChannelEntryResult struct {
	Channel collection.YTChannel
	Entries *[]youtubeapi.RSSVideoEntry
	Error   error
}

func main() {
	cfg, err := config.GetConfig(&config.FileConfigProvider{})
	if err != nil {
		panic(err)
	}

	ytLoader := collection.YTChannelLoad{}
	ytChannels, err := ytLoader.GetAvailableYTChannels(cfg)
	if err != nil {
		panic(err)
	}

	ch := make(chan ytChannelEntryResult)
	for _, chann := range *ytChannels {
		go func(channelToGet collection.YTChannel) {
			videosToGet, err := getVideoUpdatesForYTChannel(channelToGet, cfg)
			ch <- ytChannelEntryResult{channelToGet, videosToGet, err}
		}(chann)
	}

	for range *ytChannels {
		result := <-ch
		videosToGet := result.Entries
		channelToGet := result.Channel
		err := result.Error

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(channelToGet.Name())
		fmt.Printf("%d new videos available\n", len(*videosToGet))

		if len(*videosToGet) > 0 {
			command, err := youtubedl.GetCommandForArchivalType(channelToGet, videosToGet, cfg)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(command)
		}

		fmt.Println("")
	}
}

func getVideoUpdatesForYTChannel(ytc collection.YTChannel, cfg *config.Config) (*[]youtubeapi.RSSVideoEntry, error) {
	rssFeed := ytc.RSSURL()
	rss, err := youtubeapi.GetRSSFeed(rssFeed, &utils.HTTPClient{ConnTimeout: 3000000})

	if err != nil {
		return nil, fmt.Errorf("Loading RSS feed xml failed for %s collection. URL is %s. Error: %s", ytc.Name(), ytc.RSSURL(), err)
	}

	entries := rss.VideoEntry

	videos, err := ytc.GetLocalVideos(cfg)
	if err != nil {
		return nil, fmt.Errorf("Could not get videos off disk for ytc %s, error %s", ytc.Name(), err)
	}

	entriesToDownload := collection.GetEntriesNotInVideoList(&entries, videos)

	return entriesToDownload, nil
}
