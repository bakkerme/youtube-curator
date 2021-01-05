package collection

import (
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"os"
	"reflect"
	"testing"
)

func TestGetAvailableYTChannels(t *testing.T) {
	videoPath := "/video/path/"
	cfg := config.Config{
		YoutubeAPIKey: "FAKE_API_KEY",
		VideoDirPath:  videoPath,
	}

	scribePath := videoPath + "65scribe" + "/config.json"
	audiopath := videoPath + "AudioPilz" + "/config.json"

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

	t.Run("GetAvailableYTChannels returns correct channels from a mock directory", func(t *testing.T) {

		returnData := map[string][]byte{
			scribePath: []byte(`{
			  "Name": "65scribe",
			  "ID": "UC8dJOqcjyiA9Zo9aOxxiCMw",
			  "RSSURL": "https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
			  "ChannelURL": "https://www.youtube.com/user/65scribe",
			  "ArchivalMode": "archive"
			}`),
			audiopath: []byte(`{
			  "Name": "AudioPilz",
			  "ID": "UCOJVsjPZcE9HxsgPKCxZfAg",
			  "RSSURL": "https://www.youtube.com/feeds/videos.xml?channel_id=UCOJVsjPZcE9HxsgPKCxZfAg",
			  "ChannelURL": "https://www.youtube.com/channel/UCOJVsjPZcE9HxsgPKCxZfAg",
			  "ArchivalMode": "archive"
			}`),
		}

		expectedReturnData := map[string]YTChannelData{
			"65scribe": YTChannelData{
				IName:         "65scribe",
				IID:           "UC8dJOqcjyiA9Zo9aOxxiCMw",
				IRSSURL:       "https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
				IChannelURL:   "https://www.youtube.com/user/65scribe",
				IArchivalMode: ArchivalModeArchive,
			},
			"AudioPilz": YTChannelData{
				IName:         "AudioPilz",
				IID:           "UCOJVsjPZcE9HxsgPKCxZfAg",
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

		scribeEqual := reflect.DeepEqual(
			expectedReturnData["65scribe"],
			(*ytChannels)["65scribe"],
		)

		audioPilzEqual := reflect.DeepEqual(
			expectedReturnData["AudioPilz"],
			(*ytChannels)["AudioPilz"],
		)

		if !scribeEqual || !audioPilzEqual {
			t.Errorf("getAvailableYTChannels did not return correct results. Expected\n%+v\ngot\n%+v", expectedReturnData, *ytChannels)
		}
	})

	t.Run("getAvailableYTChannels returns an error on directory read error", func(t *testing.T) {
		_, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                  t,
			ShouldErrorReadDir: true,
		})

		if err == nil {
			t.Error("getAvailableYTChannels should have returned an error")
		}
	})

	t.Run("getAvailableYTChannels returns an error on file read error", func(t *testing.T) {
		_, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                   t,
			ShouldErrorReadFile: true,
			ReturnReadDirValue:  dirOutput,
		})

		if err == nil {
			t.Error("getAvailableYTChannels should have returned an error")
		}
	})

	t.Run("getAvailableYTChannels returns an error on file returning invalid json", func(t *testing.T) {
		returnData := map[string][]byte{
			scribePath: []byte(`]asdf`),
			audiopath:  []byte(`[werssd`),
		}

		_, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                          t,
			ReturnReadDirValue:         dirOutput,
			ReturnReadFileValueForPath: returnData,
		})

		if err == nil {
			t.Error("getAvailableYTChannels should have returned an error")
		}
	})
}

func TestGetLocalVideosFromDisk(t *testing.T) {
	channelName := "TestChannel"
	channel := &MockYTChannel{
		IName:                     channelName,
		IID:                       "id123",
		IRSSURL:                   "https://example.com/rss.xml",
		IChannelURL:               "https://example.com/channel/",
		IArchivalMode:             ArchivalModeArchive,
		ILocalVideos:              nil,
		ShouldErrorGetLocalVideos: false,
	}
	videoDirPath := "/videos/"

	cfg := config.Config{
		VideoDirPath: videoDirPath,
	}

	t.Run("getLocalVideosByYTChannel should load the correct results", func(t *testing.T) {
		expectedVideo := &[]Video{
			Video{
				videoDirPath + channelName + "/The Macintosh LC-dCqJ6iPHus0.mp4",
				"dCqJ6iPHus0",
				"mp4",
				videoDirPath + channelName,
			},
		}

		expectedDirname := videoDirPath + channelName
		videos, err := getLocalVideos(
			channel,
			&cfg,
			&testutils.MockDirReader{
				ExpectedDirname: &expectedDirname,
				T:               t,
				ReturnReadDirValue: &[]os.FileInfo{
					mockFileInfo{
						"The Macintosh LC-dCqJ6iPHus0.mp4",
						84000000,
						false,
					},
					mockFileInfo{
						"Some Random Invalid Video.mp4",
						84000000,
						false,
					},
				},
			},
		)

		if err != nil {
			t.Errorf("getLocalVideosByYTChannel returned unexpected error: %s", err)
		}

		if !reflect.DeepEqual(*expectedVideo, *videos) {
			t.Errorf("Video list does not match response from function: Expected\n%+v\n, got\n%+v", *expectedVideo, *videos)
		}
	})

	t.Run("getLocalVideosByYTChannel should return an error if the dir lookup fails", func(t *testing.T) {
		_, err := getLocalVideos(
			channel,
			&cfg,
			&testutils.MockDirReader{
				T:                  t,
				ShouldErrorReadDir: true,
			},
		)

		if err == nil {
			t.Error("getLocalVideosByYTChannel should have thrown error")
		}

	})
}

func TestGetLocalVideosFromDirList(t *testing.T) {
	dirlist := []os.FileInfo{
		mockFileInfo{
			"The Macintosh LC-dCqJ6iPHus0.mp4",
			84000000,
			false,
		},
		mockFileInfo{
			"Bad File.description",
			31000000,
			false,
		},
		mockFileInfo{
			"The Macintosh Quadra 800--AC4HwzAK7A.mp4",
			31000000,
			false,
		},
		mockFileInfo{
			"Bad File.x",
			31000000,
			false,
		},
		mockFileInfo{
			"The Macintosh SE-_gPsIiKtybA.mp4",
			32000000,
			false,
		},
		mockFileInfo{
			"My Video-basAIdKsyIA.mkv",
			33000000,
			false,
		},
	}

	path := "/test/path"

	expectedVideos := []Video{
		Video{
			path + "/" + dirlist[0].Name(),
			"dCqJ6iPHus0",
			"mp4",
			path,
		},
		Video{
			path + "/" + dirlist[2].Name(),
			"-AC4HwzAK7A",
			"mp4",
			path,
		},
		Video{
			path + "/" + dirlist[4].Name(),
			"_gPsIiKtybA",
			"mp4",
			path,
		},
		Video{
			path + "/" + dirlist[5].Name(),
			"basAIdKsyIA",
			"mkv",
			path,
		},
	}

	t.Run("With valid Dirlist", func(t *testing.T) {
		videos, err := getLocalVideosFromDirList(&dirlist, path)
		if err != nil {
			t.Errorf("getLocalVideosFromDirList returned an error %s", err)
		}

		if !reflect.DeepEqual(expectedVideos, *videos) {
			t.Errorf("Video list does not match response from function: Expected\n%+v\n, got\n%+v", expectedVideos, *videos)
		}
	})

	t.Run("With invalid Dirlist", func(t *testing.T) {
		dirlist := []os.FileInfo{}
		path := "/test/path"
		videos, err := getLocalVideosFromDirList(&dirlist, path)
		if err != nil {
			t.Errorf("getLocalVideosFromDirList returned an error %s", err)
		}

		if len(*videos) > 0 {
			t.Errorf("getLocalVideosFromDirList with no dir list returned something %s", videos)
		}
	})
}

func TestGetVideoIDFromFileName(t *testing.T) {
	t.Run("Parses ID from standard video title", func(t *testing.T) {
		id := "OGK8gnP4TfA"
		test := "Test Video 1-OGK8gnP4TfA.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name ID parser did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})

	t.Run("Parses ID from video title with a dash", func(t *testing.T) {
		id := "zVn7GctHoVQ"
		test := "Test Video - With a Dash-zVn7GctHoVQ.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name with a dash did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})

	t.Run("Throws error for video title with no ID", func(t *testing.T) {
		test := "Test Video - With no ID.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err == nil {
			t.Errorf("Video with no ID should return error, returned: %s", testResult)
		}
	})

	t.Run("Throws error for titless video", func(t *testing.T) {
		test := ""
		testResult, err := getVideoIDFromFileName(test)
		if err == nil {
			t.Errorf("Video with no ID should return error, returned: %s", testResult)
		}
	})

	t.Run("Parses ID from video title with dash in ID", func(t *testing.T) {
		id := "1kt7-O837H8"
		test := "Test Video - ID with dash 1kt7-O837H8.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name with a dash in ID did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})

	t.Run("Parses ID from video title a dot in the title", func(t *testing.T) {
		id := "dmZR-LFp4ns"
		test := "Test Video - with an extra dot (0.8)-dmZR-LFp4ns.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name with a dot in it did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})
}

func TestIsFileMP4(t *testing.T) {
	t.Run("A file string with an mp4 extension returns true", func(t *testing.T) {
		file := "test.mp4"
		isValid, err := isMP4(file)
		if err != nil {
			t.Errorf("isMP4 returned an error %s", err)
		}
		if !isValid {
			t.Errorf("isMP4 returned false for %s", file)
		}
	})

	t.Run("A file string with a doc extension returns false", func(t *testing.T) {
		file := "test.doc"
		isValid, err := isMP4(file)
		if err != nil {
			t.Errorf("isMP4 returned an error %s", err)
		}
		if isValid {
			t.Errorf("isMP4 returned true for %s", file)
		}
	})

	t.Run("A filename with a dot at the end but no filetype extension returns an error", func(t *testing.T) {
		file := "test."
		_, err := isMP4(file)
		if err == nil {
			t.Error("isMP4 did not return an error")
		}
	})

	t.Run("A filename with no extension returns an error", func(t *testing.T) {
		file := "test"
		_, err := isMP4(file)
		if err == nil {
			t.Error("isMP4 did not return an error")
		}
	})

	t.Run("No filename returns an error", func(t *testing.T) {
		file := ""
		_, err := isMP4(file)
		if err == nil {
			t.Error("isMP4 did not return an error")
		}
	})
}

func TestIsFileMKV(t *testing.T) {
	t.Run("A file string with an mpkv extension returns true", func(t *testing.T) {
		file := "test.mkv"
		isValid, err := isMKV(file)
		if err != nil {
			t.Errorf("isMKV returned an error %s", err)
		}
		if !isValid {
			t.Errorf("isMKV returned false for %s", file)
		}
	})

	t.Run("A file string with a doc extension returns false", func(t *testing.T) {
		file := "test.doc"
		isValid, err := isMKV(file)
		if err != nil {
			t.Errorf("isMKV returned an error %s", err)
		}
		if isValid {
			t.Errorf("isMKV returned true for %s", file)
		}
	})

	t.Run("A filename with a dot at the end but no filetype extension returns an error", func(t *testing.T) {
		file := "test."
		_, err := isMKV(file)
		if err == nil {
			t.Error("isMKV did not return an error")
		}
	})

	t.Run("A filename with no extension returns an error", func(t *testing.T) {
		file := "test"
		_, err := isMKV(file)
		if err == nil {
			t.Error("isMKV did not return an error")
		}
	})

	t.Run("No filename returns an error", func(t *testing.T) {
		file := ""
		_, err := isMKV(file)
		if err == nil {
			t.Error("isMKV did not return an error")
		}
	})
}

func TestGetFileType(t *testing.T) {
	t.Run("Standard mp4 extension should come back correct", func(t *testing.T) {
		file := "test.mp4"
		fileType, err := getFileType(file)
		if err != nil {
			t.Errorf("getFileType returned an error %s", err)
		}
		if fileType != "mp4" {
			t.Errorf("getFileType did not return correct result. Expecte mp4, got %s", fileType)
		}
	})

	t.Run("Capitalised mp4 extension should come back in lower case", func(t *testing.T) {
		file := "test.MP4"
		fileType, err := getFileType(file)
		if err != nil {
			t.Errorf("getFileType returned an error %s", err)
		}
		if fileType != "mp4" {
			t.Errorf("getFileType did not return correct result. Expecte mp4, got %s", fileType)
		}
	})

	t.Run("Throws error when dot on end but no file type exists", func(t *testing.T) {
		file := "test."
		_, err := getFileType(file)
		if err == nil {
			t.Errorf("getFileType did not return an error when provided filename %s", file)
		}
	})

	t.Run("Throws error with no extension", func(t *testing.T) {
		file := "test"
		_, err := getFileType(file)
		if err == nil {
			t.Errorf("getFileType did not return an error when provided filename %s", file)
		}
	})

	t.Run("Throws error with no filename", func(t *testing.T) {
		file := ""
		_, err := getFileType(file)
		if err == nil {
			t.Errorf("getFileType did not return an error when provided filename %s", file)
		}
	})
}
