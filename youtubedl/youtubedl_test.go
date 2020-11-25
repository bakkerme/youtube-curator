package youtubedl

import (
	"fmt"
	"strings"
	"testing"
)

var video1 = "https://www.youtube.com/watch?v=KQA9Na4aOa1"
var video2 = "https://www.youtube.com/watch?v=OGK8gnP4TfA"
var video3 = "https://www.youtube.com/watch?v=FazJqPQ6xSs"
var videoEntries = []VideoEntry{
	VideoEntry{
		"yt:video:KQA9Na4aOa1",
		"Test Video New",
		Link{
			video1,
		},
		"2020-11-06T19:00:01+00:00",
		"2020-11-06T23:12:15+00:00",
		MediaGroup{
			"Test Video 1",
			Thumbnail{
				"https://i2.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
				480,
				360,
			},
			"Test Description New",
		},
	},
	VideoEntry{
		"yt:video:OGK8gnP4TfA",
		"Test Video 1",
		Link{
			video2,
		},
		"2020-11-06T19:00:01+00:00",
		"2020-11-06T23:12:15+00:00",
		MediaGroup{
			"Test Video 1",
			Thumbnail{
				"https://i2.ytimg.com/vi/OGK8gnP4TfA/hqdefault.jpg",
				480,
				360,
			},
			"Test Description",
		},
	},
	VideoEntry{
		"yt:video:FazJqPQ6xSs",
		"Test Video 2",
		Link{
			video3,
		},
		"2020-11-06T19:00:01+00:00",
		"2020-11-06T23:12:15+00:00",
		MediaGroup{
			"Test Video 2",
			Thumbnail{
				"https://i2.ytimg.com/vi/FazJqPQ6xSs/hqdefault.jpg",
				480,
				360,
			},
			"Test Description 2",
		},
	},
}

func TestGetYoutubeDLCommandForVideoList(t *testing.T) {
	t.Run("returns correct comamnd from video list", func(t *testing.T) {
		channel := YTChannel{
			"TestChannel",
			"http://example.com/rss.xml",
			"http://example.com/channel",
			ArchivalModeCurated,
		}

		toFind := fmt.Sprintf("\"%s\" \"%s\" \"%s\"", video1, video2, video3)
		command := getYoutubeDLCommandForVideoList(&channel, &videoEntries)

		if !strings.Contains(command, toFind) {
			t.Errorf("getYoutubeDLCommandForVideoList resulted in incorrect command. Expected to find videos \n %s in command \n %s", toFind, command)
		}
	})
}

func TestCommandForArchivalType(t *testing.T) {
	t.Run("outputs channel URL for archival mode", func(t *testing.T) {
		channelURL := "http://example.com/channel"
		ytchannel := YTChannel{
			"TestChannel",
			"http://example.com/rss.xml",
			channelURL,
			ArchivalModeArchive,
		}

		result, err := GetCommandForArchivalType(&ytchannel, &videoEntries)
		if err != nil {
			t.Error(err)
		}

		if !strings.Contains(result, channelURL) {
			t.Errorf("Channel URL is incorrect: Expected %s, got %s", channelURL, result)
		}
	})

	t.Run("outputs video URLs for curated mode", func(t *testing.T) {
		channelURL := "http://example.com/channel"
		ytchannel := YTChannel{
			"TestChannel",
			"http://example.com/rss.xml",
			channelURL,
			ArchivalModeCurated,
		}

		result, err := GetCommandForArchivalType(&ytchannel, &videoEntries)
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
