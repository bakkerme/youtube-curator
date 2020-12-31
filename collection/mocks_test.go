package collection

import (
	"errors"
	"hyperfocus.systems/youtube-curator-server/config"
	"os"
	"time"
)

type mockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (f mockFileInfo) Name() string {
	return f.name
}
func (f mockFileInfo) Size() int64 {
	return f.size
}
func (f mockFileInfo) Mode() os.FileMode {
	return 0
}
func (f mockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (f mockFileInfo) IsDir() bool {
	return f.isDir
}
func (f mockFileInfo) Sys() interface{} {
	return nil
}

// mockYTChannel is a struct that represents the configuration for each channel archived
type mockYTChannel struct {
	IName                     string
	IRSSURL                   string
	IChannelURL               string
	IArchivalMode             string
	ILocalVideos              *[]Video
	shouldErrorGetLocalVideos bool
}

// GetLocalVideos is given a mockYTChannel, return the Videos on disk that are under that YTChannel
func (ytc mockYTChannel) GetLocalVideos(cf *config.Config) (*[]Video, error) {
	if ytc.shouldErrorGetLocalVideos {
		return nil, errors.New("Bad somethingarather")
	}

	return ytc.ILocalVideos, nil
}

// Name returns the name
func (ytc mockYTChannel) Name() string {
	return ytc.IName
}

// RSSURL returns the RSS URL
func (ytc mockYTChannel) RSSURL() string {
	return ytc.IRSSURL
}

// ChannelURL returns the Channel URL
func (ytc mockYTChannel) ChannelURL() string {
	return ytc.IChannelURL
}

// ArchivalMode returns the Archival Mode string
func (ytc mockYTChannel) ArchivalMode() string {
	return ytc.IArchivalMode
}

type mockYTChannelLoad struct {
	returnValue *map[string]YTChannel
	shouldError bool
}

func (ytcl mockYTChannelLoad) GetAvailableYTChannels(cf *config.Config) (*map[string]YTChannel, error) {
	if ytcl.shouldError {
		return nil, errors.New("Did the biggest error")
	}

	return ytcl.returnValue, nil
}
