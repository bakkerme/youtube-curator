package collection

import (
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetAvailableYTChannels(t *testing.T) {
	cfg := config.Config{
		YoutubeAPIKey: "FAKE_API_KEY",
		VideoDirPath:  mockVideoDirPath,
	}

	scribePath := mockVideoDirPath + "65scribe" + "/config.json"
	danBellPath := mockVideoDirPath + "DanBell" + "/config.json"

	dirOutput := GetFileInfoDirMockData()

	t.Run("GetAvailableYTChannels returns correct channels from a mock directory", func(t *testing.T) {
		returnData := mockYTConfigJSON
		expectedReturnData := MockYTChannelData

		ytChannels, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                          t,
			ReturnReadDirValue:         dirOutput,
			ExpectedDirname:            &mockVideoDirPath,
			ReturnReadFileValueForPath: returnData,
		})

		if err != nil {
			t.Error(err)
		}

		scribeEqual := reflect.DeepEqual(
			expectedReturnData["TestGuy"],
			(*ytChannels)["TestGuy"],
		)

		danBellEqual := reflect.DeepEqual(
			expectedReturnData["TestGuy2"],
			(*ytChannels)["TestGuy2"],
		)

		if !scribeEqual || !danBellEqual {
			t.Error(testutils.MismatchError("getAvailableYTChannels", expectedReturnData, *ytChannels))
		}
	})

	t.Run("getAvailableYTChannels returns an error on directory read error", func(t *testing.T) {
		_, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                  t,
			ShouldErrorReadDir: true,
		})

		if err == nil {
			t.Error(testutils.ExpectedError("getAvailableYTChannels"))
		}
	})

	t.Run("getAvailableYTChannels returns an error on file read error", func(t *testing.T) {
		_, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                   t,
			ShouldErrorReadFile: true,
			ReturnReadDirValue:  dirOutput,
		})

		if err == nil {
			t.Error(testutils.ExpectedError("getAvailableYTChannels"))
		}
	})

	t.Run("getAvailableYTChannels returns an error on file returning invalid json", func(t *testing.T) {
		returnData := map[string][]byte{
			scribePath:  []byte(`]asdf`),
			danBellPath: []byte(`[werssd`),
		}

		_, err := getAvailableYTChannels(&cfg, &testutils.MockDirReader{
			T:                          t,
			ReturnReadDirValue:         dirOutput,
			ReturnReadFileValueForPath: returnData,
		})

		if err == nil {
			t.Error(testutils.ExpectedError("getAvailableYTChannels"))
		}
	})
}

func TestGetYTChannelConfigForDirPath(t *testing.T) {
	t.Run("getYTChannelConfigForDirPath should load the correct results", func(t *testing.T) {
		ytc, err := getYTChannelConfigForDirPath(mockPathTestGuy, &testutils.MockDirReader{
			T:                          t,
			ReturnReadFileValueForPath: mockYTConfigJSON,
		})

		if err != nil {
			t.Errorf(testutils.UnexpectedError("getYTChannelConfigForDirPath", err))
		}

		expectedYTC := MockYTChannelData["TestGuy"]
		if !reflect.DeepEqual(expectedYTC, *ytc) {
			t.Error(testutils.MismatchError("getYTChannelConfigForDirPath", expectedYTC, *ytc))
		}
	})
}

func TestGetLocalVideosFromDisk(t *testing.T) {
	channel := &MockYTChannel{
		IName:                     mockChannelName,
		IID:                       "id123",
		IRSSURL:                   "https://example.com/rss.xml",
		IChannelURL:               "https://example.com/channel/",
		IArchivalMode:             ArchivalModeArchive,
		ILocalVideos:              nil,
		ShouldErrorGetLocalVideos: false,
	}

	cfg := config.Config{
		VideoDirPath: mockVideoDirPath,
	}

	t.Run("getLocalVideosByYTChannel should load the correct results", func(t *testing.T) {
		expectedVideo := &[]LocalVideo{
			(*GetVideoMockData())[0],
		}

		expectedDirname := mockVideoDirPath + mockChannelName
		videos, err := getLocalVideos(
			channel,
			&cfg,
			&testutils.MockDirReader{
				ExpectedDirname: &expectedDirname,
				T:               t,
				ReturnReadDirValue: &[]os.FileInfo{
					(*GetFileInfoMockData())[0],
					testutils.MockFileInfo{
						"Some Random Invalid Video.mp4",
						84000000,
						false,
					},
				},
			},
		)

		if err != nil {
			t.Errorf(testutils.UnexpectedError("getLocalVideosByYTChannel", err))
		}

		if !reflect.DeepEqual(*expectedVideo, *videos) {
			t.Error(testutils.MismatchError("getLocalVideosByYTChannel", *expectedVideo, *videos))
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
			t.Error(testutils.ExpectedError("getLocalVideosByYTChannel"))
		}

	})
}

func TestGetLocalVideosFromDirList(t *testing.T) {
	dirlist := *GetFileInfoMockData()
	dirlist = append(
		dirlist,
		testutils.MockFileInfo{
			"Bad File.description",
			31000000,
			false,
		},
		testutils.MockFileInfo{
			"Bad File.x",
			31000000,
			false,
		},
	)

	expectedVideos := GetVideoMockData()

	t.Run("With valid Dirlist", func(t *testing.T) {
		videos, err := getLocalVideosFromDirList(&dirlist, mockVideoDirPath+mockChannelName)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getLocalVideosFromDirList", err))
		}

		if !reflect.DeepEqual(*expectedVideos, *videos) {
			t.Error(testutils.MismatchError("getLocalVideosFromDirList", *expectedVideos, *videos))
		}
	})

	t.Run("With invalid Dirlist", func(t *testing.T) {
		dirlist := []os.FileInfo{}
		videos, err := getLocalVideosFromDirList(&dirlist, mockVideoDirPath)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getLocalVideosFromDirList", err))
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
			t.Errorf(testutils.UnexpectedError("getVideoIDFromFileName", err))
		}
		if testResult != id {
			t.Errorf(testutils.MismatchError("getVideoIDFromFileName", id, testResult))
		}
	})

	t.Run("Parses ID from video title with a dash", func(t *testing.T) {
		id := "zVn7GctHoVQ"
		test := "Test Video - With a Dash-zVn7GctHoVQ.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getVideoIDFromFileName", err))
		}
		if testResult != id {
			t.Errorf(testutils.MismatchError("getVideoIDFromFileName", id, testResult))
		}
	})

	t.Run("Throws error for video title with no ID", func(t *testing.T) {
		test := "Test Video - With no ID.mp4"
		_, err := getVideoIDFromFileName(test)
		if err == nil {
			t.Errorf(testutils.ExpectedError("getVideoIDFromFileName"))
		}
	})

	t.Run("Throws error for titless video", func(t *testing.T) {
		test := ""
		_, err := getVideoIDFromFileName(test)
		if err == nil {
			t.Errorf(testutils.ExpectedError("getVideoIDFromFileName"))
		}
	})

	t.Run("Parses ID from video title with dash in ID", func(t *testing.T) {
		id := "1kt7-O837H8"
		test := "Test Video - ID with dash 1kt7-O837H8.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getVideoIDFromFileName", err))
		}
		if testResult != id {
			t.Errorf(testutils.MismatchError("getVideoIDFromFileName", id, testResult))
		}
	})

	t.Run("Throws error for video title with no ID", func(t *testing.T) {
		test := "Test Video - With no ID.mp4"
		_, err := getVideoIDFromFileName(test)
		if err == nil {
			t.Errorf(testutils.ExpectedError("getVideoIDFromFileName"))
		}
	})

	t.Run("Throws error for titless video", func(t *testing.T) {
		test := ""
		_, err := getVideoIDFromFileName(test)
		if err == nil {
			t.Errorf(testutils.ExpectedError("getVideoIDFromFileName"))
		}
	})

	t.Run("Parses ID from video title with dash in ID", func(t *testing.T) {
		id := "1kt7-O837H8"
		test := "Test Video - ID with dash 1kt7-O837H8.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getVideoIDFromFileName", err))
		}
		if testResult != id {
			t.Errorf(testutils.MismatchError("getVideoIDFromFileName", id, testResult))
		}
	})

	t.Run("Parses ID from video title a dot in the title", func(t *testing.T) {
		id := "dmZR-LFp4ns"
		test := "Test Video - with an extra dot (0.8)-dmZR-LFp4ns.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getVideoIDFromFileName", err))
		}
		if testResult != id {
			t.Errorf(testutils.MismatchError("getVideoIDFromFileName", id, testResult))
		}
	})
}

func TestIsFileMP4(t *testing.T) {
	t.Run("A file string with an mp4 extension returns true", func(t *testing.T) {
		file := "test.mp4"
		isValid, err := isMP4(file)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("isMP4", err))
		}
		if !isValid {
			t.Errorf("isMP4 returned false for %s", file)
		}
	})

	t.Run("A file string with a doc extension returns false", func(t *testing.T) {
		file := "test.doc"
		isValid, err := isMP4(file)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("isMP4", err))
		}
		if isValid {
			t.Errorf("isMP4 returned true for %s", file)
		}
	})

	t.Run("A filename with a dot at the end but no filetype extension returns an error", func(t *testing.T) {
		file := "test."
		_, err := isMP4(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("isMP4"))
		}
	})

	t.Run("A filename with no extension returns an error", func(t *testing.T) {
		file := "test"
		_, err := isMP4(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("isMP4"))
		}
	})

	t.Run("No filename returns an error", func(t *testing.T) {
		file := ""
		_, err := isMP4(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("isMP4"))
		}
	})
}

func TestIsFileMKV(t *testing.T) {
	t.Run("A file string with an mpkv extension returns true", func(t *testing.T) {
		file := "test.mkv"
		isValid, err := isMKV(file)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("isMKV", err))
		}
		if !isValid {
			t.Errorf("isMKV returned false for %s", file)
		}
	})

	t.Run("A file string with a doc extension returns false", func(t *testing.T) {
		file := "test.doc"
		isValid, err := isMKV(file)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("isMKV", err))
		}
		if isValid {
			t.Errorf("isMKV returned true for %s", file)
		}
	})

	t.Run("A filename with a dot at the end but no filetype extension returns an error", func(t *testing.T) {
		file := "test."
		_, err := isMKV(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("isMKV"))
		}
	})

	t.Run("A filename with no extension returns an error", func(t *testing.T) {
		file := "test"
		_, err := isMKV(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("isMKV"))
		}
	})

	t.Run("No filename returns an error", func(t *testing.T) {
		file := ""
		_, err := isMKV(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("isMKV"))
		}
	})
}

func TestGetFileType(t *testing.T) {
	t.Run("Standard mp4 extension should come back correct", func(t *testing.T) {
		file := "test.mp4"
		fileType, err := getFileType(file)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getFileType", err))
		}
		if fileType != "mp4" {
			t.Error(testutils.MismatchError("getAvailableYTChannels", "mp4", fileType))
		}
	})

	t.Run("Capitalised mp4 extension should come back in lower case", func(t *testing.T) {
		file := "test.MP4"
		fileType, err := getFileType(file)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("getFileType", err))
		}
		if fileType != "mp4" {
			t.Error(testutils.MismatchError("getAvailableYTChannels", "mp4", fileType))
		}
	})

	t.Run("Throws error when dot on end but no file type exists", func(t *testing.T) {
		file := "test."
		_, err := getFileType(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("getFileType"))
		}
	})

	t.Run("Throws error with no extension", func(t *testing.T) {
		file := "test"
		_, err := getFileType(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("getFileType"))
		}
	})

	t.Run("Throws error with no filename", func(t *testing.T) {
		file := ""
		_, err := getFileType(file)
		if err == nil {
			t.Errorf(testutils.ExpectedError("getFileType"))
		}
	})
}

func TestCheckYTChannelConfig(t *testing.T) {
	channel := MockYTChannelData["TestGuy"]

	t.Run("Should return nil for a valid config", func(t *testing.T) {
		err := checkYTChannelConfig(&channel)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("checkYTChannelConfig", err))
		}
	})

	t.Run("Should return nil for a valid config with ArchivalModeCurated", func(t *testing.T) {
		ytc := channel
		ytc.IArchivalMode = ArchivalModeCurated

		err := checkYTChannelConfig(&channel)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("checkYTChannelConfig", err))
		}
	})

	t.Run("Should return nil for a valid config with ChannelTypePlaylist", func(t *testing.T) {
		ytc := channel
		ytc.IChannelType = ChannelTypePlaylist

		err := checkYTChannelConfig(&channel)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("checkYTChannelConfig", err))
		}
	})

	checkFieldError := func(ytc *YTChannelData, field string) {
		err := checkYTChannelConfig(ytc)
		if err == nil {
			t.Error(testutils.ExpectedError("checkYTChannelConfig"))
		}

		if !strings.Contains(err.Error(), field) {
			t.Errorf("expected %s to contain %s", err.Error(), field)
		}
	}

	t.Run("Should return error for an invalid name", func(t *testing.T) {
		ytc := channel
		ytc.IName = ""
		checkFieldError(&ytc, "name")
	})

	t.Run("Should return error for an invalid id", func(t *testing.T) {
		ytc := channel
		ytc.IID = ""
		checkFieldError(&ytc, "id")
	})

	t.Run("Should return error for an invalid rssURL", func(t *testing.T) {
		ytc := channel
		ytc.IRSSURL = ""
		checkFieldError(&ytc, "rssURL")
	})

	t.Run("Should return error for an invalid channelURL", func(t *testing.T) {
		ytc := channel
		ytc.IChannelURL = ""
		checkFieldError(&ytc, "channelURL")
	})

	t.Run("Should return error for an empty archivalMode", func(t *testing.T) {
		ytc := channel
		ytc.IArchivalMode = ""
		checkFieldError(&ytc, "archivalMode")
	})

	t.Run("Should return error for an invalid archivalMode", func(t *testing.T) {
		ytc := channel
		ytc.IArchivalMode = "asdfsadfasdf"
		checkFieldError(&ytc, "archivalMode")
	})

	t.Run("Should return error for an invalid channelType", func(t *testing.T) {
		ytc := channel
		ytc.IChannelType = ""
		checkFieldError(&ytc, "channelType")
	})

	t.Run("Should return for all errors at once", func(t *testing.T) {
		ytc := channel
		ytc.IName = ""
		ytc.IID = ""
		ytc.IRSSURL = ""
		ytc.IChannelURL = ""
		ytc.IArchivalMode = ""
		ytc.IChannelType = ""

		checkFieldError(&ytc, "name")
		checkFieldError(&ytc, "id")
		checkFieldError(&ytc, "rssURL")
		checkFieldError(&ytc, "channelURL")
		checkFieldError(&ytc, "archivalMode")
		checkFieldError(&ytc, "channelType")
	})
}
