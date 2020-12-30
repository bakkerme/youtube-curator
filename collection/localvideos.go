package collection

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"regexp"
	"strings"
)

// GetVideoByID finds a local Video file on disk with a provided ID
func GetVideoByID(ID string, cf *config.Config) (*Video, error) {
	videos, err := GetAllLocalVideos(cf)
	if err != nil {
		return nil, fmt.Errorf("Could not get Video with ID %s. Error %s", ID, err)
	}

	for _, video := range *videos {
		if video.ID == ID {
			return &video, nil
		}
	}

	return nil, nil
}

// GetAllLocalVideos returns a full list of Videos found locally
func GetAllLocalVideos(cf *config.Config) (*[]Video, error) {
	return getAllLocalVideos(cf, &YTChannelLoad{})
}

func getAllLocalVideos(cf *config.Config, ytcl ytChannelLoader) (*[]Video, error) {
	channels, err := ytcl.GetAvailableYTChannels(cf)
	if err != nil {
		return nil, fmt.Errorf("Cannot get all YT Channels. Got error %s", err)
	}

	var videoList []Video
	for chName, ch := range *channels {
		v, err := ch.GetLocalVideos(cf)
		if err != nil {
			return nil, fmt.Errorf("Cannot get local videos for %s Channels. Got error %s", chName, err)
		}

		videoList = append(videoList, *v...)
	}

	return &videoList, nil
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
