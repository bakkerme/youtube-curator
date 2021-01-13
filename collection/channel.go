package collection

import (
	"encoding/json"
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"os"
	"regexp"
	"strings"
)

// YTChannelLoader provides an interface for loading YT Channels
type YTChannelLoader interface {
	GetAvailableYTChannels(cf *config.Config) (*map[string]YTChannel, error)
}

// YTChannelLoad allows YT Channels to be loaded
type YTChannelLoad struct{}

// GetAvailableYTChannels returns all configured Youtube Channels on disk in the configured video directory
func (ytcl YTChannelLoad) GetAvailableYTChannels(cf *config.Config) (*map[string]YTChannel, error) {
	return getAvailableYTChannels(cf, &utils.DirReader{})
}

func getAvailableYTChannels(cf *config.Config, dr utils.DirReaderProvider) (*map[string]YTChannel, error) {
	ytChannels := map[string]YTChannel{}

	dirEntries, err := dr.ReadDir(cf.VideoDirPath)
	if err != nil {
		return nil, err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() && dirEntry.Name()[0] != '.' {
			ytChannelConfig, err := getYTChannelConfigForDirPath(cf.VideoDirPath+dirEntry.Name()+"/config.json", dr)
			if err != nil {
				return nil, err
			}

			ytChannels[ytChannelConfig.Name()] = *ytChannelConfig
		}
	}

	return &ytChannels, nil
}

func getYTChannelConfigForDirPath(path string, dr utils.DirReaderProvider) (*YTChannelData, error) {
	file, err := dr.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Can't get YT Channel config file. Looking for %s, got error %s", path, err)
	}

	var resp YTChannelData
	if err := json.Unmarshal([]byte(file), &resp); err != nil {
		return nil, fmt.Errorf("Can't unmarshal YT Channel config file. Looking for %s, got error %s", path, err)
	}

	if invalid := checkYTChannelConfig(&resp); invalid != nil {
		return nil, fmt.Errorf("Channel config at %s is invalid. %s", path, invalid)
	}

	return &resp, nil
}

func checkYTChannelConfig(ytc *YTChannelData) error {
	invalidFields := []string{}

	invalidFields = *appendIfInvalidYTC(ytc.Name(), "name", "notequal", &invalidFields, "")
	invalidFields = *appendIfInvalidYTC(ytc.ID(), "id", "notequal", &invalidFields, "")
	invalidFields = *appendIfInvalidYTC(ytc.RSSURL(), "rssURL", "notequal", &invalidFields, "")
	invalidFields = *appendIfInvalidYTC(ytc.ChannelURL(), "channelURL", "notequal", &invalidFields, "")
	invalidFields = *appendIfInvalidYTC(ytc.ArchivalMode(), "archivalMode", "equal", &invalidFields, ArchivalModeCurated, ArchivalModeArchive)
	invalidFields = *appendIfInvalidYTC(ytc.ChannelType(), "channelType", "equal", &invalidFields, ChannelTypeChannel, ChannelTypePlaylist)

	if len(invalidFields) > 0 {
		return fmt.Errorf("%s fields invalid", strings.Join(invalidFields, ","))
	}

	return nil
}

func appendIfInvalidYTC(value string, field string, operation string, invalidFields *[]string, invalidCases ...string) *[]string {
	match := false
	for _, ic := range invalidCases {
		if operation == "notequal" && value == ic {
			result := append(*invalidFields, field)
			return &result
		}

		if operation == "equal" && value == ic {
			match = true
		}
	}

	if operation == "equal" && !match {
		result := append(*invalidFields, field)
		return &result
	}

	return invalidFields
}

func getLocalVideos(channel YTChannel, cf *config.Config, dr utils.DirReaderProvider) (*[]LocalVideo, error) {
	path := cf.VideoDirPath + channel.Name()
	dirlist, err := dr.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return getLocalVideosFromDirList(&dirlist, path)
}

func getLocalVideosFromDirList(dirlist *[]os.FileInfo, path string) (*[]LocalVideo, error) {
	var videos []LocalVideo
	for _, file := range *dirlist {
		valid, _ := isValidVideo(file.Name())

		if valid {
			videoPath := path + "/" + file.Name()

			re := regexp.MustCompile(`(?m)\.(.*)$`)
			thumbPath := path + "/" + string(re.ReplaceAll([]byte(file.Name()), []byte(".png")))

			id, err := getVideoIDFromFileName(file.Name())
			if err != nil {
				return &videos, err
			}

			extension, err := getFileType(file.Name())
			if err != nil {
				return nil, err
			}

			video := LocalVideo{
				Path:      videoPath,
				ID:        id,
				FileType:  extension,
				BasePath:  path,
				Thumbnail: thumbPath,
			}

			videos = append(videos, video)
		}
	}

	return &videos, nil
}

func isValidVideo(filename string) (bool, error) {
	is, err := isMP4(filename)
	if is && doesVideoHaveID(filename) {
		return is, err
	}

	is, err = isMKV(filename)
	if is && doesVideoHaveID(filename) {
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

func doesVideoHaveID(filename string) bool {
	_, err := getVideoIDFromFileName(filename)
	if err != nil {
		return false
	}

	return true
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
