package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"log"
	"strings"
)

func convertISODateToSimpleDateString(isoDateString string) string {
	split := strings.Split(isoDateString, "T")
	return strings.ReplaceAll(split[0], "-", "")
}

func tryYoutubeAPI() {
	config, err := config.GetConfig(&config.EnvarConfigProvider{})
	if err != nil {
		log.Panicf("Config Loader threw an error %s", err)
	}

	v65scribe := collection.Feeds["65scribe"]
	videos, err := collection.GetLocalVideosByYTChannel(&v65scribe)
	if err != nil {
		log.Panicf("Could not get local videos, error %s", err)
	}

	var ids []string
	for _, video := range *videos {
		ids = append(ids, video.ID)
	}

	videoList, err := youtubeapi.GetVideoInfo(&ids, config)
	if err != nil {
		panic(fmt.Sprintf("Could not get video of id %s, %s", ids, err))
	}

	if videoList.NextPageToken != "" {
		panic(fmt.Sprintf("Next Page!!!! %d of %d per page", videoList.PageInfo.TotalResults, videoList.PageInfo.ResultsPerPage))
	}

	for _, video := range videoList.Items {
		fmt.Println(video.Snippet.Title)
		fmt.Println(convertISODateToSimpleDateString(video.Snippet.PublishedAt))
		fmt.Println("")
	}

	// fmt.Println(videoList)
}
