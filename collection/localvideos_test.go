package collection

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetVideoMetadata(t *testing.T) {
	t.Run("getVideoMetadata should return correct result", func(t *testing.T) {
		publishedAt, err := time.Parse("2006-01-02", "1992-05-01")
		if err != nil {
			t.Error(err)
		}

		duration, err := time.ParseDuration("3m")
		if err != nil {
			t.Error(err)
		}

		video := Video{
			Path:     "/a/path/Channel1/20201118 - Installing Red Hat Linux 8.0 on the $5 Windows 98 PC!-44E3kV_6p24.mkv",
			ID:       "44E3kV_6p24",
			FileType: "mkv",
			BasePath: "/a/path/Channel1/",
		}

		mt := videometadata.Metadata{
			Title:       "Installing Red Hat Linux 8.0 on the $5 Windows 98 PC!",
			Description: "A description is here",
			Creator:     "Channel1",
			PublishedAt: &publishedAt,
			Duration:    &duration,
		}

		expectedVWM := VideoWithMetadata{
			mt,
			video,
		}

		mockVideoMetadata := videometadata.MockVideoMetadata{
			GetReturn: &videometadata.Response{
				Metadata:   &mt,
				ParseError: nil,
			},
			GetReturnError: false,
			SetReturnError: false,
		}

		vwm, err := getVideoMetadata(&video, &mockVideoMetadata)

		if err != nil {
			t.Error(err)
		}

		if fmt.Sprintf("%+v", expectedVWM) != fmt.Sprintf("%+v", *vwm) {
			t.Errorf("getVideoMetadata returned incorrect result. Expected\n%+v\ngot\n%+v", expectedVWM, *vwm)
		}
	})

	t.Run("getVideoMetadata should return error if file type is not supported", func(t *testing.T) {
		video := Video{
			Path:     "/a/path/Channel1/20201118 - Installing Red Hat Linux 8.0 on the $5 Windows 98 PC!-44E3kV_6p24.mkv",
			ID:       "44E3kV_6p24",
			FileType: "mkv",
			BasePath: "/a/path/Channel1/",
		}

		_, err := getVideoMetadata(&video, &videometadata.VideoMetadata{})
		if err == nil {
			t.Errorf("getVideoMetadata should return an error")
		}
	})

	t.Run("getVideoMetadata should return error if metadata provider fails", func(t *testing.T) {
		video := Video{
			Path:     "/a/path/Channel1/20201118 - Installing Red Hat Linux 8.0 on the $5 Windows 98 PC!-44E3kV_6p24.mkv",
			ID:       "44E3kV_6p24",
			FileType: "mkv",
			BasePath: "/a/path/Channel1/",
		}

		mockVideoMetadata := videometadata.MockVideoMetadata{
			GetReturn:      nil,
			GetReturnError: true,
			SetReturnError: false,
		}

		_, err := getVideoMetadata(&video, &mockVideoMetadata)

		if err == nil {
			t.Errorf("getVideoMetadata should return an error")
		}
	})
}

func TestGetVideoByID(t *testing.T) {
	t.Run("getVideoByID should return correct result", func(t *testing.T) {
		path := "/test/pass/"
		channel1 := "Channel1"
		channel2 := "Channel2"

		cf := &config.Config{
			VideoDirPath: path,
		}

		channel1Videos := []Video{
			Video{
				fmt.Sprintf("%s%s/%s", path, channel1, "20201118 - Installing Red Hat Linux 8.0 on the $5 Windows 98 PC!-44E3kV_6p24.mkv"),
				"44E3kV_6p24",
				"mp4",
				path + channel1,
			},
		}

		channel2Videos := []Video{
			Video{
				fmt.Sprintf("%s%s/%s", path, channel2, "20200622 - The US electrical system is not 120V-jMmUoZh3Hq4.mkv"),
				"jMmUoZh3Hq4",
				"mp4",
				path + channel2,
			},
		}

		ytcl := &mockYTChannelLoad{
			returnValue: &map[string]YTChannel{
				"Channel1": mockYTChannel{
					channel1,
					"123abc.com",
					"123abc.com",
					ArchivalModeArchive,
					&channel1Videos,
					false,
				},
				"Channel2": mockYTChannel{
					channel1,
					"123abc.com",
					"123abc.com",
					ArchivalModeArchive,
					&channel2Videos,
					false,
				},
			},
		}

		video, err := getVideoByID("jMmUoZh3Hq4", cf, ytcl)

		if err != nil {
			t.Errorf("getVideoByID returned an unexpected error %s", err)
		}

		expectedVideo := channel2Videos[0]
		if !reflect.DeepEqual(*video, expectedVideo) {
			t.Errorf("getVideoByID returned incorrect result. Expected\n%+vgot\n%+v", expectedVideo, *video)
		}
	})

	t.Run("getVideoByID should return nil if ID is not found", func(t *testing.T) {
		path := "/test/pass/"
		channel1 := "Channel1"
		channel2 := "Channel2"

		cf := &config.Config{
			VideoDirPath: path,
		}

		channel1Videos := []Video{
			Video{
				fmt.Sprintf("%s%s/%s", path, channel1, "20201118 - Installing Red Hat Linux 8.0 on the $5 Windows 98 PC!-44E3kV_6p24.mkv"),
				"44E3kV_6p24",
				"mp4",
				path + channel1,
			},
		}

		channel2Videos := []Video{
			Video{
				fmt.Sprintf("%s%s/%s", path, channel2, "20200622 - The US electrical system is not 120V-jMmUoZh3Hq4.mkv"),
				"jMmUoZh3Hq4",
				"mp4",
				path + channel2,
			},
		}

		ytcl := &mockYTChannelLoad{
			returnValue: &map[string]YTChannel{
				"Channel1": mockYTChannel{
					channel1,
					"123abc.com",
					"123abc.com",
					ArchivalModeArchive,
					&channel1Videos,
					false,
				},
				"Channel2": mockYTChannel{
					channel1,
					"123abc.com",
					"123abc.com",
					ArchivalModeArchive,
					&channel2Videos,
					false,
				},
			},
		}

		video, err := getVideoByID("SOMEBADID", cf, ytcl)

		if err != nil {
			t.Errorf("getVideoByID returned an unexpected error %s", err)
		}

		if video != nil {
			t.Errorf("video should be nil if ID can't be found")
		}
	})

	t.Run("getVideoByID should return an error if it can't get all the videos", func(t *testing.T) {
		cf := &config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := &mockYTChannelLoad{
			shouldError: true,
		}

		_, err := getVideoByID("someID", cf, ytcl)

		if err == nil {
			t.Error("getVideoByID should have thrown an error")
		}
	})
}

func TestGetAllLocalVideos(t *testing.T) {
	t.Run("getAllLocalVideos should load the correct results", func(t *testing.T) {
		path := "/test/pass/"
		channel := "Channel1"
		fileName := "The Macintosh LC-dCqJ6iPHus0.mp4"

		cf := &config.Config{
			VideoDirPath: path,
		}

		expectedVideos := []Video{
			Video{
				fmt.Sprintf("%s%s/%s", path, channel, fileName),
				"dCqJ6iPHus0",
				"mp4",
				path + channel,
			},
		}

		ytcl := &mockYTChannelLoad{
			returnValue: &map[string]YTChannel{
				"Channel1": mockYTChannel{
					channel,
					"123abc.com",
					"123abc.com",
					ArchivalModeArchive,
					&expectedVideos,
					false,
				},
			},
		}

		videos, err := getAllLocalVideos(cf, ytcl)

		if err != nil {
			t.Errorf("getAllLocalVideos returned an unexpected error %s", err)
		}

		if !reflect.DeepEqual(expectedVideos, *videos) {
			t.Errorf("getAllLocalVideos did not return correct results. Expected\n%+v\ngot\n%+v", expectedVideos, *videos)
		}
	})

	t.Run("getAllLocalVideos should error when the ytChannelLoader fails", func(t *testing.T) {
		cf := &config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := &mockYTChannelLoad{
			shouldError: true,
		}

		_, err := getAllLocalVideos(cf, ytcl)

		if err == nil {
			t.Error("getAllLocalVideos should have thrown an error")
		}
	})

	t.Run("getAllLocalVideos should return an empty []Video when the channel loader has no channels", func(t *testing.T) {
		cf := &config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := &mockYTChannelLoad{
			returnValue: nil,
		}

		videos, err := getAllLocalVideos(cf, ytcl)

		if err != nil {
			t.Errorf("getAllLocalVideos returned an unexpected error %s", err)
		}

		if len(*videos) > 0 {
			t.Errorf("videos should be empty. Got %+v", *videos)
		}
	})

	t.Run("getAllLocalVideos should return an error if any Channel can't load videos", func(t *testing.T) {
		cf := &config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := &mockYTChannelLoad{
			returnValue: &map[string]YTChannel{
				"Channel1": mockYTChannel{
					"Channel1",
					"123abc.com",
					"123abc.com",
					ArchivalModeArchive,
					&[]Video{
						Video{
							"Video1",
							"dCqJ6iPHus0",
							"mp4",
							"/a/path",
						},
					},
					false,
				},
				"Channel2": mockYTChannel{
					"Channel2",
					"123abc.com",
					"123abc.com",
					ArchivalModeArchive,
					&[]Video{
						Video{
							"Video1",
							"dCqJ6iPHus0",
							"mp4",
							"/a/path",
						},
					},
					true,
				},
			},
		}

		_, err := getAllLocalVideos(cf, ytcl)

		if err == nil {
			t.Error("getAllLocalVideos should have thrown an error")
		}
	})
}

func TestGetLocalVideosFromDisk(t *testing.T) {
	t.Run("getLocalVideosByYTChannel should load the correct results", func(t *testing.T) {
		channelName := "TestChannel"
		channel := &mockYTChannel{
			channelName,
			"https://example.com/rss.xml",
			"https://example.com/channel/",
			ArchivalModeArchive,
			nil,
			false,
		}

		videoDirPath := "/videos/"

		expectedVideo := &[]Video{
			Video{
				videoDirPath + channelName + "/The Macintosh LC-dCqJ6iPHus0.mp4",
				"dCqJ6iPHus0",
				"mp4",
				videoDirPath + channelName,
			},
		}

		cfg := config.Config{
			VideoDirPath: videoDirPath,
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
		channelName := "TestChannel"
		channel := &mockYTChannel{
			channelName,
			"https://example.com/rss.xml",
			"https://example.com/channel/",
			ArchivalModeArchive,
			nil,
			false,
		}

		cfg := config.Config{
			VideoDirPath: "/videos/",
		}

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
