package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// YTChannel is a struct that represents the configuration for each channel archived
type YTChannel struct {
	Name         string
	RSSURL       string
	ChannelURL   string
	ArchivalMode string
}

// Video is a struct that represents a single video on disk
type Video struct {
	Path string
	ID   string
}

// ArchivalModeArchive specifies that all videos are to be archived
const ArchivalModeArchive = "archive"

// ArchivalModeCurated specifies that only selected videos are to be archived
const ArchivalModeCurated = "curated"

var feeds = map[string]YTChannel{
	"65scribe": YTChannel{
		"65scribe",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
		"https://www.youtube.com/user/65scribe",
		ArchivalModeArchive,
	},
	"Ashens": YTChannel{
		"Ashens",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCxt9Pvye-9x_AIcb1UtmF1Q",
		"https://www.youtube.com/user/ashens",
		ArchivalModeArchive,
	},
	"BryanLunduke": YTChannel{
		"BryanLunduke",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkK9UDm_ZNrq_rIXCz3xCGA",
		"https://www.youtube.com/user/bryanlunduke",
		ArchivalModeArchive,
	},
	"DanBell": YTChannel{
		"DanBell",
		"https://www.youtube.com/feeds/videos.xml?playlist_id=PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
		"https://www.youtube.com/user/moviedan",
		ArchivalModeArchive,
	},
	"LinusTechTips": YTChannel{
		"LinusTechTips",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCXuqSBlHAE6Xw-yeJA0Tunw",
		"https://www.youtube.com/user/LinusTechTips",
		ArchivalModeCurated,
	},
	"LukeSmith": YTChannel{
		"LukeSmith",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC2eYFnH61tmytImy1mTYvhA",
		"https://www.youtube.com/channel/UC2eYFnH61tmytImy1mTYvhA",
		ArchivalModeArchive,
	},
	"Mario64BetaArchive": YTChannel{
		"Mario64BetaArchive",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSar92RCPocysNbvBG84Mxw",
		"https://www.youtube.com/channel/UCSar92RCPocysNbvBG84Mxw",
		ArchivalModeArchive,
	},
	"Memospore": YTChannel{
		"Memospore",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChbm-JCx_jii5xn-2f5nwIA",
		"https://www.youtube.com/user/memospore",
		ArchivalModeCurated,
	},
	"MichaelMJD": YTChannel{
		"MichaelMJD",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCS-WzPVpAAli-1IfEG2lN8A",
		"https://www.youtube.com/user/mjd7999",
		ArchivalModeArchive,
	},
	"RedLetterMedia": YTChannel{
		"RedLetterMedia",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCrTNhL_yO3tPTdQ5XgmmWjA",
		"https://www.youtube.com/user/RedLetterMedia",
		ArchivalModeCurated,
	},
	"SurviveTheJive": YTChannel{
		"SurviveTheJive",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCZAENaOaceQUMd84GDc26EA",
		"https://www.youtube.com/user/ThomasRowsell",
		ArchivalModeArchive,
	},
	"TechnologyConnections": YTChannel{
		"TechnologyConnections",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCy0tKL1T7wFoYcxCe0xjN6Q",
		"https://www.youtube.com/channel/UCy0tKL1T7wFoYcxCe0xjN6Q",
		ArchivalModeArchive,
	},
}

// Given a filename of a video on disk (created with youtube-dl), grab
// the filename from it and extract the video ID
func getVideoIDFromFileName(filename string) (string, error) {
	parseError := fmt.Errorf("Could not parse video ID for video %s", filename)

	re := regexp.MustCompile(`(\..{3}$)`)
	withoutType := re.ReplaceAllString(filename, "")

	id := withoutType[len(withoutType)-11 : len(withoutType)] // Get last 11 chars

	if id == "" {
		return "", parseError
	}

	if len(id) != len(strings.ReplaceAll(id, " ", "")) { // This is probably not an ID
		return "", parseError
	}

	return id, nil
}

// Given an Entry from the RSS feed, and a list of Videos on disk,
// return the Entrys that are not represented on disk
func isEntryInVideoList(entry *Entry, videos *[]Video) bool {
	match := false
	for _, video := range *videos {
		if video.ID == entry.ID {
			match = true
		}
	}

	return match
}

// Given a list of Entries from the RSS feed and Videos on disk, return
// the Entries that don't appear as a Video on disk
func getEntriesNotInVideoList(entries *[]Entry, videos *[]Video) *[]Entry {
	var notInVideoList []Entry
	for _, entry := range *entries {
		match := isEntryInVideoList(&entry, videos)
		if !match { // Entry isn't in our list of videos
			notInVideoList = append(notInVideoList, entry)
		}
	}

	return &notInVideoList
}

// Given a YTChannel, return the Videos on disk that are under that YTChannel
func getLocalVideosByYTChannel(channel *YTChannel) (*[]Video, error) {
	path := "/media/Drive/Videos/Youtube/" + channel.Name
	dirlist, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return getLocalVideosFromDirList(&dirlist, path)
}

func getLocalVideosFromDirList(dirlist *[]os.FileInfo, path string) (*[]Video, error) {
	var videos []Video
	for _, file := range *dirlist {
		isValidVideo, _ := isMP4(file.Name())

		if isValidVideo {
			videoPath := path + "/" + file.Name()
			id, err := getVideoIDFromFileName(videoPath)
			if err != nil {
				return &videos, err
			}

			video := Video{
				videoPath,
				id,
			}

			videos = append(videos, video)
		}
	}

	return &videos, nil
}

func isMP4(filename string) (bool, error) {
	fileType, err := getFileType(filename)
	return fileType == "mp4", err
}

func getFileType(filename string) (string, error) {
	split := strings.Split(filename, ".")
	final := split[len(split)-1]

	if len(split) <= 1 || final == "" {
		return "", errors.New("Invalid file type, must have extension")
	}

	return strings.ToLower(final), nil
}
