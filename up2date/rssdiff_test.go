package main

import (
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"testing"
)

func TestGetEntriesNotInVideoList(t *testing.T) {
	t.Run("With a Match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{
			URL:    "",
			Width:  0,
			Height: 0,
		}

		mediaGroup := youtubeapi.RSSMediaGroup{
			Title:       "",
			Thumbnail:   thumbnail,
			Description: "",
		}

		outstandingEntry := youtubeapi.RSSVideoEntry{
			ID:         "wxyzabcdefg",
			Title:      "Video 4",
			Link:       youtubeapi.RSSLink{Href: "http://link4"},
			Published:  "",
			Updated:    "",
			MediaGroup: mediaGroup,
		}
		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{
				ID:         "12345678911",
				Title:      "Video 1",
				Link:       youtubeapi.RSSLink{Href: "http://link"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
			youtubeapi.RSSVideoEntry{
				ID:         "abcdefghijk",
				Title:      "Video 2",
				Link:       youtubeapi.RSSLink{Href: "http://link2"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
			youtubeapi.RSSVideoEntry{
				ID:         "lmnopqrxtuv",
				Title:      "Video 3",
				Link:       youtubeapi.RSSLink{Href: "http://link3"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
			outstandingEntry,
		}

		videos := []collection.Video{
			collection.Video{
				Path:     "Video 1-12345678911.mp4",
				ID:       "12345678911",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 2-abcdefghijk.mp4",
				ID:       "abcdefghijk",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 3-lmnopqrxtuv.mp4",
				ID:       "lmnopqrxtuv",
				FileType: "mp4",
				BasePath: "/",
			},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if (*entriesToGet)[0] != outstandingEntry {
			t.Errorf("The outstanding entry not in the video list is incorrect, got %s, expected %s", (*entriesToGet)[0].ID, outstandingEntry.ID)
		}
	})

	t.Run("With no match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{
			URL:    "",
			Width:  0,
			Height: 0,
		}

		mediaGroup := youtubeapi.RSSMediaGroup{
			Title:       "",
			Thumbnail:   thumbnail,
			Description: "",
		}

		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{
				ID:         "12345678911",
				Title:      "Video 1",
				Link:       youtubeapi.RSSLink{Href: "http://link"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
			youtubeapi.RSSVideoEntry{
				ID:         "abcdefghijk",
				Title:      "Video 2",
				Link:       youtubeapi.RSSLink{Href: "http://link2"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
			youtubeapi.RSSVideoEntry{
				ID:         "lmnopqrxtuv",
				Title:      "Video 3",
				Link:       youtubeapi.RSSLink{Href: "http://link3"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
		}

		videos := []collection.Video{
			collection.Video{
				Path:     "Video 1-12345678911.mp4",
				ID:       "12345678911",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 2-abcdefghijk.mp4",
				ID:       "abcdefghijk",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 3-lmnopqrxtuv.mp4",
				ID:       "lmnopqrxtuv",
				FileType: "mp4",
				BasePath: "/",
			},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) > 0 {
			t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
		}
	})

	t.Run("With no video list provided", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{
			URL:    "",
			Width:  0,
			Height: 0,
		}

		mediaGroup := youtubeapi.RSSMediaGroup{
			Title:       "",
			Thumbnail:   thumbnail,
			Description: "",
		}

		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{
				ID:         "12345678911",
				Title:      "Video 1",
				Link:       youtubeapi.RSSLink{Href: "http://link"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
			youtubeapi.RSSVideoEntry{
				ID:         "abcdefghijk",
				Title:      "Video 2",
				Link:       youtubeapi.RSSLink{Href: "http://link2"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
			youtubeapi.RSSVideoEntry{
				ID:         "lmnopqrxtuv",
				Title:      "Video 3",
				Link:       youtubeapi.RSSLink{Href: "http://link3"},
				Published:  "",
				Updated:    "",
				MediaGroup: mediaGroup,
			},
		}

		videos := []collection.Video{}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) != 3 {
			t.Errorf("Matches were not correctly found")
		}
	})

	t.Run("With no entries provied", func(t *testing.T) {
		entries := []youtubeapi.RSSVideoEntry{}

		videos := []collection.Video{
			collection.Video{
				Path:     "Video 1-12345678911.mp4",
				ID:       "12345678911",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 2-abcdefghijk.mp4",
				ID:       "abcdefghijk",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 3-lmnopqrxtuv.mp4",
				ID:       "lmnopqrxtuv",
				FileType: "mp4",
				BasePath: "/",
			},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) > 0 {
			t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
		}
	})

}

func TestIsEntryInVideoList(t *testing.T) {
	t.Run("With match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{
			URL:    "",
			Width:  0,
			Height: 0,
		}

		mediaGroup := youtubeapi.RSSMediaGroup{
			Title:       "",
			Thumbnail:   thumbnail,
			Description: "",
		}
		entry := youtubeapi.RSSVideoEntry{
			ID:         "12345678911",
			Title:      "Video 1",
			Link:       youtubeapi.RSSLink{Href: "http://link"},
			Published:  "",
			Updated:    "",
			MediaGroup: mediaGroup,
		}

		videos := []collection.Video{
			collection.Video{
				Path:     "Video 1-12345678911.mp4",
				ID:       "12345678911",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 2-abcdefghijk.mp4",
				ID:       "abcdefghijk",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 3-lmnopqrxtuv.mp4",
				ID:       "lmnopqrxtuv",
				FileType: "mp4",
				BasePath: "/",
			},
		}

		if !isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should have been found in video list")
		}
	})

	t.Run("With no match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{
			URL:    "",
			Width:  0,
			Height: 0,
		}

		mediaGroup := youtubeapi.RSSMediaGroup{
			Title:       "",
			Thumbnail:   thumbnail,
			Description: "",
		}
		entry := youtubeapi.RSSVideoEntry{
			ID:         "BADID123456",
			Title:      "Video 1",
			Link:       youtubeapi.RSSLink{Href: "http://link"},
			Published:  "",
			Updated:    "",
			MediaGroup: mediaGroup,
		}

		videos := []collection.Video{
			collection.Video{
				Path:     "Video 1-12345678911.mp4",
				ID:       "12345678911",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 2-abcdefghijk.mp4",
				ID:       "abcdefghijk",
				FileType: "mp4",
				BasePath: "/",
			},
			collection.Video{
				Path:     "Video 3-lmnopqrxtuv.mp4",
				ID:       "lmnopqrxtuv",
				FileType: "mp4",
				BasePath: "/",
			},
		}

		if isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should not have been found in video list")
		}
	})

	t.Run("With empty video list", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{
			URL:    "",
			Width:  0,
			Height: 0,
		}

		mediaGroup := youtubeapi.RSSMediaGroup{
			Title:       "",
			Thumbnail:   thumbnail,
			Description: "",
		}
		entry := youtubeapi.RSSVideoEntry{
			ID:         "BADID123456",
			Title:      "Video 1",
			Link:       youtubeapi.RSSLink{Href: "http://link"},
			Published:  "",
			Updated:    "",
			MediaGroup: mediaGroup,
		}

		videos := []collection.Video{}

		if isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should not have been found in video list")
		}
	})
}
