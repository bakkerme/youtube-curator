package youtubedl

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"strings"
	"testing"
)

var video1 = "https://www.youtube.com/watch?v=KQA9Na4aOa1"
var video2 = "https://www.youtube.com/watch?v=OGK8gnP4TfA"
var video3 = "https://www.youtube.com/watch?v=FazJqPQ6xSs"
var videoEntries = []youtubeapi.RSSVideoEntry{
	youtubeapi.RSSVideoEntry{
		ID:    "yt:video:KQA9Na4aOa1",
		Title: "Test Video New",
		Link: youtubeapi.RSSLink{
			Href: video1,
		},
		Published: "2020-11-06T19:00:01+00:00",
		Updated:   "2020-11-06T23:12:15+00:00",
		MediaGroup: youtubeapi.RSSMediaGroup{
			Title: "Test Video 1",
			Thumbnail: youtubeapi.RSSThumbnail{
				URL:    "https://i2.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
				Width:  480,
				Height: 360,
			},
			Description: "Test Description New",
		},
	},
	youtubeapi.RSSVideoEntry{
		ID:    "yt:video:OGK8gnP4TfA",
		Title: "Test Video 1",
		Link: youtubeapi.RSSLink{
			Href: video2,
		},
		Published: "2020-11-06T19:00:01+00:00",
		Updated:   "2020-11-06T23:12:15+00:00",
		MediaGroup: youtubeapi.RSSMediaGroup{
			Title: "Test Video 1",
			Thumbnail: youtubeapi.RSSThumbnail{
				URL:    "https://i2.ytimg.com/vi/OGK8gnP4TfA/hqdefault.jpg",
				Width:  480,
				Height: 360,
			},
			Description: "Test Description",
		},
	},
	youtubeapi.RSSVideoEntry{
		ID:    "yt:video:FazJqPQ6xSs",
		Title: "Test Video 2",
		Link: youtubeapi.RSSLink{
			Href: video3,
		},
		Published: "2020-11-06T19:00:01+00:00",
		Updated:   "2020-11-06T23:12:15+00:00",
		MediaGroup: youtubeapi.RSSMediaGroup{
			Title: "Test Video 2",
			Thumbnail: youtubeapi.RSSThumbnail{
				URL:    "https://i2.ytimg.com/vi/FazJqPQ6xSs/hqdefault.jpg",
				Width:  480,
				Height: 360,
			},
			Description: "Test Description 2",
		},
	},
}

var mockConfig = &config.Config{
	VideoDirPath: "/base/path/",
}

func TestGetYoutubeDLCommandForVideoList(t *testing.T) {
	t.Run("returns correct comamnd from video list", func(t *testing.T) {
		channel := collection.YTChannelData{
			IName:         "TestChannel",
			IID:           "asdfasdf",
			IRSSURL:       "http://example.com/rss.xml",
			IChannelURL:   "http://example.com/channel",
			IArchivalMode: collection.ArchivalModeCurated,
		}

		toFind := fmt.Sprintf("\"%s\" \"%s\" \"%s\"", video1, video2, video3)
		command := getYoutubeDLCommandForVideoList(&channel, &videoEntries, "/base/path")

		if !strings.Contains(command, toFind) {
			t.Errorf("getYoutubeDLCommandForVideoList resulted in incorrect command. Expected to find videos \n %s in command \n %s", toFind, command)
		}
	})
}

func TestCommandForArchivalType(t *testing.T) {

	t.Run("outputs channel URL for archival mode", func(t *testing.T) {
		channelURL := "http://example.com/channel"
		ytchannel := collection.YTChannelData{
			IName:         "TestChannel",
			IID:           "asdfasdf",
			IRSSURL:       "http://example.com/rss.xml",
			IChannelURL:   channelURL,
			IArchivalMode: collection.ArchivalModeArchive,
		}

		result, err := GetCommandForArchivalType(&ytchannel, &videoEntries, mockConfig)
		if err != nil {
			t.Error(err)
		}

		if !strings.Contains(result, channelURL) {
			t.Errorf("Channel URL is incorrect: Expected %s, got %s", channelURL, result)
		}
	})

	t.Run("outputs video URLs for curated mode", func(t *testing.T) {
		channelURL := "http://example.com/channel"
		ytchannel := collection.YTChannelData{
			IName:         "TestChannel",
			IID:           "asdfasdf",
			IRSSURL:       "http://example.com/rss.xml",
			IChannelURL:   channelURL,
			IArchivalMode: collection.ArchivalModeCurated,
		}

		result, err := GetCommandForArchivalType(&ytchannel, &videoEntries, mockConfig)
		if err != nil {
			t.Error(err)
		}

		doesContain := []bool{
			strings.Contains(result, video1),
			strings.Contains(result, video2),
			strings.Contains(result, video3),
		}

		if !doesContain[0] || !doesContain[1] || !doesContain[2] {
			videoString := video1 + " " + video2 + " " + video3
			t.Errorf("Command did not result in expected url: Expected %s, got %s", videoString, result)
		}
	})
}
