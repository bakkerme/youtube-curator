package collection

import (
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"hyperfocus.systems/youtube-curator-server/videometadata"
)

// ArchivalModeArchive specifies that all videos are to be archived
const ArchivalModeArchive = "archive"

// ArchivalModeCurated specifies that only selected videos are to be archived
const ArchivalModeCurated = "curated"

// ChannelTypeChannel represents a YTChannel that is a channel
const ChannelTypeChannel = "channel"

// ChannelTypePlaylist represents a YTChannel that is a playlist
const ChannelTypePlaylist = "playlist"

// YTChannel provides an interface for interacting with a YT Channel
type YTChannel interface {
	GetLocalVideos(cf *config.Config) (*[]LocalVideo, error)
	ID() string
	Name() string
	RSSURL() string
	ChannelURL() string
	ArchivalMode() string
	ChannelType() string
}

// LocalVideo is a struct that represents a single video on disk
type LocalVideo struct {
	Path      string
	ID        string
	FileType  string
	BasePath  string
	Thumbnail string
}

// LocalVideoWithMetadata represents a video on the filesystem,
// along with the metadata from that video
type LocalVideoWithMetadata struct {
	videometadata.Metadata
	LocalVideo
}

// YTChannelData is a struct that represents the configuration for each channel archived
type YTChannelData struct {
	IName         string `json:"name"`
	IID           string `json:"id"`
	IRSSURL       string `json:"rssURL"`
	IChannelURL   string `json:"channelURL"`
	IArchivalMode string `json:"archivalMode"`
	IChannelType  string `json:"channelType"`
}

// GetLocalVideos is given a YTChannelData, return the Videos on disk that are under that YTChannel
func (ytc YTChannelData) GetLocalVideos(cf *config.Config) (*[]LocalVideo, error) {
	return getLocalVideos(&ytc, cf, &utils.DirReader{})
}

// Name returns the name
func (ytc YTChannelData) Name() string {
	return ytc.IName
}

// ID returns the ID
func (ytc YTChannelData) ID() string {
	return ytc.IID
}

// RSSURL returns the RSS URL
func (ytc YTChannelData) RSSURL() string {
	return ytc.IRSSURL
}

// ChannelURL returns the Channel URL
func (ytc YTChannelData) ChannelURL() string {
	return ytc.IChannelURL
}

// ArchivalMode returns the Archival Mode string
func (ytc YTChannelData) ArchivalMode() string {
	return ytc.IArchivalMode
}

// ChannelType returns the Channel Type string
func (ytc YTChannelData) ChannelType() string {
	return ytc.IChannelType
}
