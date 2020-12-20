package collection

import (
	"encoding/json"
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
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
	Path     string
	ID       string
	FileType string
	BasePath string
}

// ArchivalModeArchive specifies that all videos are to be archived
const ArchivalModeArchive = "archive"

// ArchivalModeCurated specifies that only selected videos are to be archived
const ArchivalModeCurated = "curated"

// GetAvailableYTChannels returns all configured Youtube Channels on disk in the configured video directory
func GetAvailableYTChannels(cf *config.Config) (*map[string]YTChannel, error) {
	return getAvailableYTChannels(cf, &utils.DirReader{})
}

func getAvailableYTChannels(cf *config.Config, dr utils.DirReaderProvider) (*map[string]YTChannel, error) {
	ytChannels := map[string]YTChannel{}

	dirEntries, err := dr.ReadDir(cf.VideoDirPath)
	if err != nil {
		return nil, err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			ytChannelConfig, err := getYTChannelConfigForDirPath(cf.VideoDirPath+dirEntry.Name()+"/config.json", dr)
			if err != nil {
				return nil, err
			}
			ytChannels[ytChannelConfig.Name] = *ytChannelConfig
		}
	}

	return &ytChannels, nil
}

func getYTChannelConfigForDirPath(path string, dr utils.DirReaderProvider) (*YTChannel, error) {
	file, err := dr.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Can't get YT Channel config file. Looking for %s, got error %s", path, err)
	}

	var resp YTChannel
	if err := json.Unmarshal([]byte(file), &resp); err != nil {
		return nil, fmt.Errorf("Can't unmarshal YT Channel config file. Looking for %s, got error %s", path, err)
	}

	return &resp, nil
}

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

// Given a filename of a video on disk (created with youtube-dl), grab
// the filename from it and extract the video ID
func getVideoIDFromFileName(filename string) (string, error) {
	parseError := fmt.Errorf("Could not parse video ID for video %s", filename)

	re := regexp.MustCompile(`(\..{3}$)`)
	withoutType := re.ReplaceAllString(filename, "")

	if len(withoutType) < 11 {
		return "", parseError
	}

	id := withoutType[len(withoutType)-11 : len(withoutType)] // Get last 11 chars

	if len(id) != len(strings.ReplaceAll(id, " ", "")) { // This is probably not an ID
		return "", parseError
	}

	return id, nil
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

// GetLocalVideosByYTChannel is given a YTChannel, return the Videos on disk that are under that YTChannel
func GetLocalVideosByYTChannel(channel *YTChannel, cf *config.Config) (*[]Video, error) {
	return getLocalVideosFromDisk(channel, &utils.DirReader{}, cf)
}

func getLocalVideosFromDisk(channel *YTChannel, dr utils.DirReaderProvider, cf *config.Config) (*[]Video, error) {
	path := cf.VideoDirPath + channel.Name
	dirlist, err := dr.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return getLocalVideosFromDirList(&dirlist, path)
}

func getLocalVideosFromDirList(dirlist *[]os.FileInfo, path string) (*[]Video, error) {
	var videos []Video
	for _, file := range *dirlist {
		isValidVideo, _ := isValidVideo(file.Name())

		if isValidVideo {
			videoPath := path + "/" + file.Name()
			id, err := getVideoIDFromFileName(file.Name())
			if err != nil {
				return &videos, err
			}

			extension, err := getFileType(file.Name())
			if err != nil {
				return nil, err
			}

			video := Video{
				videoPath,
				id,
				extension,
				path,
			}

			videos = append(videos, video)
		}
	}

	return &videos, nil
}

func isValidVideo(filename string) (bool, error) {
	is, err := isMP4(filename)
	if is {
		return is, err
	}

	is, err = isMKV(filename)
	if is {
		return is, err
	}

	return false, nil
}

func isMP4(filename string) (bool, error) {
	fileType, err := getFileType(filename)
	return fileType == "mp4", err
}

func isMKV(filename string) (bool, error) {
	fileType, err := getFileType(filename)
	return fileType == "mkv", err
}

func getFileType(filename string) (string, error) {
	split := strings.Split(filename, ".")
	final := split[len(split)-1]

	if len(split) <= 1 || final == "" {
		return "", errors.New("Invalid file type, must have extension")
	}

	return strings.ToLower(final), nil
}
