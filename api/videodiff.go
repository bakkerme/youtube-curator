package api

import (
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
)

// getEntriesNotInVideoList is given a list of Entries from the RSS feed and Videos on disk, return
// the Entries that don't appear as a Video on disk
func getEntriesNotInVideoList(remoteVideos *[]youtubeapi.Video, localVideos *[]collection.LocalVideo) *[]youtubeapi.Video {
	var notInVideoList []youtubeapi.Video = []youtubeapi.Video{}
	for _, rmv := range *remoteVideos {
		match := isEntryInVideoList(&rmv, localVideos)
		if !match {
			notInVideoList = append(notInVideoList, rmv)
		}
	}

	return &notInVideoList
}

// Given an RSSVideoEntry from the RSS feed, and a list of Videos on disk,
// return the Entrys that are not represented on disk
func isEntryInVideoList(removeVideo *youtubeapi.Video, localVideos *[]collection.LocalVideo) bool {
	match := false
	for _, video := range *localVideos {
		if video.ID == removeVideo.ID {
			match = true
		}
	}

	return match
}
