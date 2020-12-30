package collection

import (
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
)

// GetEntriesNotInVideoList is given a list of Entries from the RSS feed and Videos on disk, return
// the Entries that don't appear as a Video on disk
func GetEntriesNotInVideoList(entries *[]youtubeapi.RSSVideoEntry, videos *[]Video) *[]youtubeapi.RSSVideoEntry {
	var notInVideoList []youtubeapi.RSSVideoEntry
	for _, entry := range *entries {
		match := isEntryInVideoList(&entry, videos)
		if !match { // RSSVideoEntry isn't in our list of videos
			notInVideoList = append(notInVideoList, entry)
		}
	}

	return &notInVideoList
}

// Given an RSSVideoEntry from the RSS feed, and a list of Videos on disk,
// return the Entrys that are not represented on disk
func isEntryInVideoList(entry *youtubeapi.RSSVideoEntry, videos *[]Video) bool {
	match := false
	for _, video := range *videos {
		if video.ID == entry.ID {
			match = true
		}
	}

	return match
}
