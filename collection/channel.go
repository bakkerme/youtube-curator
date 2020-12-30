package collection

import (
	"encoding/json"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"os"
)

type ytChannelLoader interface {
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
		if dirEntry.IsDir() {
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

	return &resp, nil
}

func getLocalVideos(channel YTChannel, cf *config.Config, dr utils.DirReaderProvider) (*[]Video, error) {
	path := cf.VideoDirPath + channel.Name()
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
