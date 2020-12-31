package collection

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/mkvmetadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/mp4metadata"
	"regexp"
	"strings"
)

// GetVideoMetadata looks up the metadata for a given Video and returns a VideoWithMetadata
// object containing all known information on the video, both Metadata and Video
func GetVideoMetadata(video *Video) (*VideoWithMetadata, error) {
	return getVideoMetadata(video, &videometadata.VideoMetadata{})
}

func getVideoMetadata(video *Video, vm videometadata.Provider) (*VideoWithMetadata, error) {
	metadataProvider, err := getMetadataCommandProviderForFileType(video.Path)
	if err != nil {
		return nil, err
	}

	resp, err := vm.Get(video.Path, metadataProvider)
	if err != nil {
		return nil, err
	}

	mt := resp.Metadata
	return &VideoWithMetadata{
		*mt,
		*video,
	}, nil

}

func getMetadataCommandProviderForFileType(filetype string) (videometadata.CommandProvider, error) {
	mp4, err := isMP4(filetype)
	if err != nil {
		return nil, err
	}

	mkv, err := isMKV(filetype)
	if err != nil {
		return nil, err
	}

	if mp4 {
		return mp4metadata.MP4MetadataCommandProvider{}, nil
	}

	if mkv {
		return mkvmetadata.MKVMetadataCommandProvider{}, nil
	}

	return nil, nil
}

// GetVideoByID finds a local Video file on disk with a provided ID
func GetVideoByID(ID string, cf *config.Config) (*Video, error) {
	return getVideoByID(ID, cf, &YTChannelLoad{})
}

func getVideoByID(ID string, cf *config.Config, ytcl ytChannelLoader) (*Video, error) {
	videos, err := getAllLocalVideos(cf, ytcl)
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

	if channels == nil {
		return &[]Video{}, nil
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
	if filename == "mp4" {
		return true, nil
	}

	fileType, err := getFileType(filename)
	return fileType == "mp4", err
}

func isMKV(filename string) (bool, error) {
	if filename == "mkv" {
		return true, nil
	}

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
