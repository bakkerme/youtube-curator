package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/tageditor"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type videoUpdateResult struct {
	Path  string
	Error error
}

func startFileRenamer() {
	cfg, err := config.GetConfig(&config.EnvarConfigProvider{})
	if err != nil {
		panic(err)
	}

	ytc := collection.Feeds["TechnologyConnections"]
	videos, err := collection.GetLocalVideosByYTChannel(&ytc, cfg)
	if err != nil {
		panic(err)
	}

	var idsToGet []string
	for _, video := range *videos {
		if video.FileType == "mp4" {
			idsToGet = append(idsToGet, video.ID)
		}
	}

	chuckSize := 50
	chuckedIDs := chunkIntSlice(&idsToGet, chuckSize)
	for chunkID, ids := range *chuckedIDs {
		ytAPIVideos := getVideoFromYoutubeAPI(&ids)

		ch := make(chan videoUpdateResult)
		for id, ytVideo := range *ytAPIVideos {
			collectionVideo := (*videos)[(chunkID*chuckSize)+id]
			go func(yv youtubeapi.Video, cv collection.Video) {
				updateVideoFileToChannel(ch, yv, cv)
			}(ytVideo, collectionVideo)
		}

		for range *ytAPIVideos {
			result := <-ch
			if result.Error != nil {
				fmt.Printf("Failed to write %s. Error %s\n", result.Path, result.Error)
			} else {
				fmt.Printf("Wrote %s\n", result.Path)
			}
		}
	}
}

func updateVideoFileToChannel(ch chan videoUpdateResult, ytVideo youtubeapi.Video, collectionVideo collection.Video) {
	updated, err := updateVideoFile(ytVideo, collectionVideo)
	ch <- videoUpdateResult{
		updated,
		err,
	}
}

func updateVideoFile(ytVideo youtubeapi.Video, collectionVideo collection.Video) (string, error) {
	snippet := &ytVideo.Snippet

	publishedAt, err := time.Parse(time.RFC3339, snippet.PublishedAt)
	if err != nil {
		panic(err)
	}

	metadataProvider := &tageditor.MP4MetadataCommandProvider{}

	if err != nil {
		return "", err
	}

	metadataToWrite := &videometadata.Metadata{
		snippet.Title,
		snippet.Description,
		snippet.ChannelTitle,
		&publishedAt,
		nil,
	}

	videometadata.SetVideoMetadata(collectionVideo.Path, metadataToWrite, metadataProvider)

	newName := fmt.Sprintf("%s - %s-%s.%s",
		timeToSimpleDateString(publishedAt),
		sanitiseCMDInput(snippet.Title),
		sanitiseCMDInput(collectionVideo.ID),
		sanitiseCMDInput(collectionVideo.FileType),
	)

	err = os.Rename(collectionVideo.Path, fmt.Sprintf("%s/%s", collectionVideo.BasePath, newName))
	if err != nil {
		return "", err
	}

	return newName, nil
}

func getVideoFromYoutubeAPI(ids *[]string) *[]youtubeapi.Video {
	config, err := config.GetConfig(&config.EnvarConfigProvider{})
	if err != nil {
		log.Panicf("Config Loader threw an error %s", err)
	}

	videoList, err := youtubeapi.GetVideoInfo(ids, config)
	if err != nil {
		panic(fmt.Sprintf("Could not get video of ids %s, %s", *ids, err))
	}

	if videoList.NextPageToken != "" {
		panic(fmt.Sprintf("Next Page!!!! %d of %d per page", videoList.PageInfo.TotalResults, videoList.PageInfo.ResultsPerPage))
	}

	return &videoList.Items
}

func timeToSimpleDateString(t time.Time) string {
	return t.Format("20060102")
}

func sanitiseCMDInput(str string) string {
	return strings.ReplaceAll(str, "/", "_")
}

func chunkIntSlice(stringSlice *[]string, itemsPerChunk int) *[][]string {
	var stringChunkSlice [][]string
	chunks := math.Ceil(float64(len(*stringSlice)) / float64(itemsPerChunk))

	for i := 0; i < int(chunks); i++ {
		start := itemsPerChunk * i
		end := start + itemsPerChunk
		if end > len(*stringSlice) {
			end = len(*stringSlice)
		}
		stringChunkSlice = append(stringChunkSlice, (*stringSlice)[start:end])
	}

	return &stringChunkSlice
}
