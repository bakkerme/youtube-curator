package main

import (
	"fmt"
	"testing"
)

func testGetYoutubeDLCommandForVideList(t *testing.T) {
	video1 := "https://www.youtube.com/watch?v=KQA9Na4aOa1"
	video2 := "https://www.youtube.com/watch?v=OGK8gnP4TfA"
	video3 := "https://www.youtube.com/watch?v=FazJqPQ6xSs"

	entries := []VideoEntry{
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

	channel := YTChannel{
		"TestChannel",
		"http://example.com/rss.xml",
		"http://example.com/channel",
		ArchivalModeArchive,
	}

	command := getYoutubeDLCommandForVideoList(&channel, &entries)

	expected := getYoutubeDLCommandForYTChannel(&channel, fmt.Sprintf("\"%s\" \"%s\" \"%s\"", video1, video2, video3))
	if command != expected {
		t.Errorf("getYoutubeDLCommandForVideoList resulted in incorrect command. Expected %s got %s", expected, command)
	}
}
