package collection

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/mkvmetadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/mp4metadata"
)

// GetVideoMetadata looks up the metadata for a given Video and returns a VideoWithMetadata
// object containing all known information on the video, both Metadata and Video
func GetVideoMetadata(video *LocalVideo) (*LocalVideoWithMetadata, error) {
	return getVideoMetadata(video, &videometadata.VideoMetadata{})
}

func getVideoMetadata(video *LocalVideo, vm videometadata.Provider) (*LocalVideoWithMetadata, error) {
	metadataProvider, err := getMetadataCommandProviderForFileType(video.Path)
	if err != nil {
		return nil, err
	}

	resp, err := vm.Get(video.Path, metadataProvider)
	if err != nil {
		return nil, err
	}

	mt := resp.Metadata
	return &LocalVideoWithMetadata{
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
		return mp4metadata.CommandProvider{}, nil
	}

	if mkv {
		return mkvmetadata.CommandProvider{}, nil
	}

	return nil, fmt.Errorf("Cannot find metadata parser for %s", filetype)
}

// GetVideoByID finds a local Video file on disk with a provided ID
func GetVideoByID(ID string, cf *config.Config) (*LocalVideo, error) {
	return getVideoByID(ID, cf, &YTChannelLoad{})
}

func getVideoByID(ID string, cf *config.Config, ytcl YTChannelLoader) (*LocalVideo, error) {
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
func GetAllLocalVideos(cf *config.Config) (*[]LocalVideo, error) {
	return getAllLocalVideos(cf, &YTChannelLoad{})
}

func getAllLocalVideos(cf *config.Config, ytcl YTChannelLoader) (*[]LocalVideo, error) {
	channels, err := ytcl.GetAvailableYTChannels(cf)
	if err != nil {
		return nil, fmt.Errorf("Cannot get all YT Channels. Got error %s", err)
	}

	if channels == nil {
		return &[]LocalVideo{}, nil
	}

	var videoList []LocalVideo
	for chName, ch := range *channels {
		v, err := ch.GetLocalVideos(cf)
		if err != nil {
			return nil, fmt.Errorf("Cannot get local videos for %s Channels. Got error %s", chName, err)
		}

		videoList = append(videoList, *v...)
	}

	return &videoList, nil
}
