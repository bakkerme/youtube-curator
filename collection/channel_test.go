package collection

import (
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"os"
	"reflect"
	"testing"
)

func TestGetAvailableYTChannels(t *testing.T) {
	t.Run("GetAvailableYTChannels returns correct channels from a mock directory", func(t *testing.T) {
		videoPath := "/video/path/"
		cfg := config.Config{
			YoutubeAPIKey: "FAKE_API_KEY",
			VideoDirPath:  videoPath,
		}

		dirOutput := &[]os.FileInfo{
			mockFileInfo{
				name:  "65scribe",
				size:  0,
				isDir: true,
			},
			mockFileInfo{
				name:  "AudioPilz",
				size:  0,
				isDir: true,
			},
			mockFileInfo{
				name:  "somerandomvideo.mp4",
				size:  12312312312321,
				isDir: false,
			},
		}

		scribePath := videoPath + "65scribe" + "/config.json"
		audiopath := videoPath + "AudioPilz" + "/config.json"

		returnData := map[string][]byte{
			scribePath: []byte(`{
			  "Name": "65scribe",
			  "RSSURL": "https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
			  "ChannelURL": "https://www.youtube.com/user/65scribe",
			  "ArchivalMode": "archive"
			}`),
			audiopath: []byte(`{
			  "Name": "AudioPilz",
			  "RSSURL": "https://www.youtube.com/feeds/videos.xml?channel_id=UCOJVsjPZcE9HxsgPKCxZfAg",
			  "ChannelURL": "https://www.youtube.com/channel/UCOJVsjPZcE9HxsgPKCxZfAg",
			  "ArchivalMode": "archive"
			}`),
		}

		expectedReturnData := map[string]YTChannelData{
			"65scribe": YTChannelData{
				IName:         "65scribe",
				IRSSURL:       "https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
				IChannelURL:   "https://www.youtube.com/user/65scribe",
				IArchivalMode: ArchivalModeArchive,
			},
			"AudioPilz": YTChannelData{
				IName:         "AudioPilz",
				IRSSURL:       "https://www.youtube.com/feeds/videos.xml?channel_id=UCOJVsjPZcE9HxsgPKCxZfAg",
				IChannelURL:   "https://www.youtube.com/channel/UCOJVsjPZcE9HxsgPKCxZfAg",
				IArchivalMode: ArchivalModeArchive,
			},
		}

		ytChannels, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                          t,
			ReturnReadDirValue:         dirOutput,
			ExpectedDirname:            &videoPath,
			ReturnReadFileValueForPath: returnData,
		})

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expectedReturnData["65scribe"], (*ytChannels)["65scribe"]) || !reflect.DeepEqual(expectedReturnData["AudioPilz"], (*ytChannels)["AudioPilz"]) {
			t.Errorf("getAvailableYTChannels did not return correct results. Expected\n%+v\ngot\n%+v", expectedReturnData, *ytChannels)
		}
	})
}

func TestGetEntriesNotInVideoList(t *testing.T) {
	t.Run("With a Match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}

		outstandingEntry := youtubeapi.RSSVideoEntry{"wxyzabcdefg", "Video 4", youtubeapi.RSSLink{"http://link4"}, "", "", mediaGroup}
		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"abcdefghijk", "Video 2", youtubeapi.RSSLink{"http://link2"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"lmnopqrxtuv", "Video 3", youtubeapi.RSSLink{"http://link3"}, "", "", mediaGroup},
			outstandingEntry,
		}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if (*entriesToGet)[0] != outstandingEntry {
			t.Errorf("The outstanding entry not in the video list is incorrect, got %s, expected %s", (*entriesToGet)[0].ID, outstandingEntry.ID)
		}
	})

	t.Run("With no match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}

		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"abcdefghijk", "Video 2", youtubeapi.RSSLink{"http://link2"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"lmnopqrxtuv", "Video 3", youtubeapi.RSSLink{"http://link3"}, "", "", mediaGroup},
		}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) > 0 {
			t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
		}
	})

	t.Run("With no video list provided", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}

		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"abcdefghijk", "Video 2", youtubeapi.RSSLink{"http://link2"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"lmnopqrxtuv", "Video 3", youtubeapi.RSSLink{"http://link3"}, "", "", mediaGroup},
		}

		videos := []Video{}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) != 3 {
			t.Errorf("Matches were not correctly found")
		}
	})

	t.Run("With no entries provied", func(t *testing.T) {
		entries := []youtubeapi.RSSVideoEntry{}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) > 0 {
			t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
		}
	})

}

func TestIsEntryInVideoList(t *testing.T) {
	t.Run("With match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}
		entry := youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		if !isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should have been found in video list")
		}
	})

	t.Run("With no match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}
		entry := youtubeapi.RSSVideoEntry{"BADID123456", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		if isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should not have been found in video list")
		}
	})

	t.Run("With empty video list", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}
		entry := youtubeapi.RSSVideoEntry{"BADID123456", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup}

		videos := []Video{}

		if isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should not have been found in video list")
		}
	})
}
