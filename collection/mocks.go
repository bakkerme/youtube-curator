package collection

import (
	"errors"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"os"
)

// MockYTChannel is a struct that represents the configuration for each channel archived
type MockYTChannel struct {
	IName                     string
	IID                       string
	IRSSURL                   string
	IChannelURL               string
	IArchivalMode             string
	IChannelType              string
	ILocalVideos              *[]LocalVideo
	ShouldErrorGetLocalVideos bool
}

// GetLocalVideos is given a mockYTChannel, return the Videos on disk that are under that YTChannel
func (ytc MockYTChannel) GetLocalVideos(cf *config.Config) (*[]LocalVideo, error) {
	if ytc.ShouldErrorGetLocalVideos {
		return nil, errors.New("Bad somethingarather")
	}

	return ytc.ILocalVideos, nil
}

// Name returns the name
func (ytc MockYTChannel) Name() string {
	return ytc.IName
}

// ID returns the ID
func (ytc MockYTChannel) ID() string {
	return ytc.IID
}

// RSSURL returns the RSS URL
func (ytc MockYTChannel) RSSURL() string {
	return ytc.IRSSURL
}

// ChannelURL returns the Channel URL
func (ytc MockYTChannel) ChannelURL() string {
	return ytc.IChannelURL
}

// ArchivalMode returns the Archival Mode string
func (ytc MockYTChannel) ArchivalMode() string {
	return ytc.IArchivalMode
}

// ChannelType returns the Channel Title string
func (ytc MockYTChannel) ChannelType() string {
	return ytc.IChannelType
}

// MockYTChannelLoad mocks out the YTChannelLoad interface
type MockYTChannelLoad struct {
	ReturnValue *map[string]YTChannel
	ShouldError bool
}

// GetAvailableYTChannels mocks the GetAvailableYTChannels function
func (ytcl MockYTChannelLoad) GetAvailableYTChannels(cf *config.Config) (*map[string]YTChannel, error) {
	if ytcl.ShouldError {
		return nil, errors.New("Did the biggest error")
	}

	return ytcl.ReturnValue, nil
}

var mockVideoDirPath = "/a/path/"
var mockChannelName = "TestGuy"
var mockChannelName2 = "TestGuy2"

// GetVideoMockData returns some mock data to test with. This data corresponds with the YoutubeAPI package video mock data
func GetVideoMockData() *[]LocalVideo {
	return &[]LocalVideo{
		LocalVideo{
			Path:      mockVideoDirPath + mockChannelName + "/Test Video New-18-elPdai_1.mp4",
			ID:        "18-elPdai_1",
			FileType:  "mp4",
			BasePath:  mockVideoDirPath + mockChannelName,
			Thumbnail: mockVideoDirPath + mockChannelName + "/Test Video New-18-elPdai_1.png",
		},

		LocalVideo{
			Path:      mockVideoDirPath + mockChannelName + "/Test Video 1-OGK8gnP4TfA.mp4",
			ID:        "OGK8gnP4TfA",
			FileType:  "mp4",
			BasePath:  mockVideoDirPath + mockChannelName,
			Thumbnail: mockVideoDirPath + mockChannelName + "/Test Video 1-OGK8gnP4TfA.png",
		},
		LocalVideo{
			Path:      mockVideoDirPath + "TestGuy/Test Video 2-FazJqPQ6xSs.mkv",
			ID:        "FazJqPQ6xSs",
			FileType:  "mkv",
			BasePath:  mockVideoDirPath + mockChannelName,
			Thumbnail: mockVideoDirPath + mockChannelName + "/Test Video 2-FazJqPQ6xSs.png",
		},
	}
}

// GetFileInfoDirMockData returns test directory data
func GetFileInfoDirMockData() *[]os.FileInfo {
	return &[]os.FileInfo{
		testutils.MockFileInfo{
			IName:  mockChannelName,
			ISize:  0,
			IIsDir: true,
		},
		testutils.MockFileInfo{
			IName:  "TestGuy2",
			ISize:  0,
			IIsDir: true,
		},
		testutils.MockFileInfo{
			IName:  "somerandomvideo.mp4",
			ISize:  12312312312321,
			IIsDir: false,
		},
	}

}

// GetFileInfoMockData returns mock file data that corresponds
// to the results of GetVideoMockData
func GetFileInfoMockData() *[]os.FileInfo {
	return &[]os.FileInfo{
		testutils.MockFileInfo{
			"Test Video New-18-elPdai_1.mp4",
			84000000,
			false,
		},
		testutils.MockFileInfo{
			"Test Video 1-OGK8gnP4TfA.mp4",
			31000000,
			false,
		},
		testutils.MockFileInfo{
			"Test Video 2-FazJqPQ6xSs.mkv",
			32000000,
			false,
		},
	}
}

var mockPathTestGuy = mockVideoDirPath + "TestGuy/config.json"
var mockPathTestGuy2 = mockVideoDirPath + "TestGuy2/config.json"
var mockYTConfigJSON = map[string][]byte{
	mockPathTestGuy: []byte(`{
			  "Name": "TestGuy",
			  "ID": "UC8dJOqcjyiA9Zo9aOxxiCMw",
			  "RSSURL": "https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
			  "ChannelURL": "https://www.youtube.com/user/TestGuy",
			  "ArchivalMode": "archive",
			  "ChannelType": "channel"
			}`),
	mockPathTestGuy2: []byte(`{
			  "Name": "TestGuy2",
			  "RSSURL": "https://www.youtube.com/feeds/videos.xml?playlist_id=PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
			  "ID": "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
			  "ChannelURL": "https://www.youtube.com/user/moviedan",
			  "ArchivalMode": "archive",
			  "ChannelType": "playlist"
			}`),
}

// MockYTChannelData contains mock YTChannel for testing purposes
var MockYTChannelData = map[string]YTChannelData{
	mockChannelName: YTChannelData{
		IName:         mockChannelName,
		IID:           "UC8dJOqcjyiA9Zo9aOxxiCMw",
		IRSSURL:       "https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
		IChannelURL:   "https://www.youtube.com/user/TestGuy",
		IArchivalMode: ArchivalModeArchive,
		IChannelType:  ChannelTypeChannel,
	},
	mockChannelName2: YTChannelData{
		IName:         mockChannelName2,
		IID:           "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
		IRSSURL:       "https://www.youtube.com/feeds/videos.xml?playlist_id=PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
		IChannelURL:   "https://www.youtube.com/user/moviedan",
		IArchivalMode: ArchivalModeArchive,
		IChannelType:  ChannelTypePlaylist,
	},
}
