package collection

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/mkvmetadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/mp4metadata"
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
			Path:     "/a/path/Channel1/20201118 - Installing Red Hat Linux 8.0 on the $5 Windows 98 PC!-44E3kV_6p24.asd",
			ID:       "44E3kV_6p24",
			FileType: "asd",
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

func TestGetMetadataCommandProviderForFileType(t *testing.T) {
	t.Run("should return correct provider for mp4", func(t *testing.T) {
		cmdProv, err := getMetadataCommandProviderForFileType("test.mp4")
		if err != nil {
			t.Errorf("Received unexpected error %s", err)
		}

		expect := mp4metadata.CommandProvider{}
		if cmdProv != expect {
			t.Errorf("Did not receive correct command provider. Got %+v", cmdProv)
		}
	})

	t.Run("should return correct provider for mkv", func(t *testing.T) {
		cmdProv, err := getMetadataCommandProviderForFileType("test.mkv")
		if err != nil {
			t.Errorf("Received unexpected error %s", err)
		}

		expect := mkvmetadata.CommandProvider{}
		if cmdProv != expect {
			t.Errorf("Did not receive correct command provider. Got %+v", cmdProv)
		}
	})

	t.Run("should return error for unsupported file type", func(t *testing.T) {
		_, err := getMetadataCommandProviderForFileType("test.asdf")
		if err == nil {
			t.Error("Expected to receive error")
		}
	})

	t.Run("should return error for no file type", func(t *testing.T) {
		_, err := getMetadataCommandProviderForFileType("test")
		if err == nil {
			t.Error("Expected to receive error")
		}
	})

	t.Run("should return error for dot with no file type", func(t *testing.T) {
		_, err := getMetadataCommandProviderForFileType("test.")
		if err == nil {
			t.Error("Expected to receive error")
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

		ytcl := &MockYTChannelLoad{
			ReturnValue: &map[string]YTChannel{
				"Channel1": MockYTChannel{
					IName:                     channel1,
					IID:                       "id123",
					IRSSURL:                   "123abc.com",
					IChannelURL:               "123abc.com",
					IArchivalMode:             ArchivalModeArchive,
					ILocalVideos:              &channel1Videos,
					ShouldErrorGetLocalVideos: false,
				},
				"Channel2": MockYTChannel{
					IName:                     channel1,
					IID:                       "id123",
					IRSSURL:                   "123abc.com",
					IChannelURL:               "123abc.com",
					IArchivalMode:             ArchivalModeArchive,
					ILocalVideos:              &channel2Videos,
					ShouldErrorGetLocalVideos: false,
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

		ytcl := &MockYTChannelLoad{
			ReturnValue: &map[string]YTChannel{
				"Channel1": MockYTChannel{
					IName:                     channel1,
					IID:                       "id123",
					IRSSURL:                   "123abc.com",
					IChannelURL:               "123abc.com",
					IArchivalMode:             ArchivalModeArchive,
					ILocalVideos:              &channel1Videos,
					ShouldErrorGetLocalVideos: false,
				},
				"Channel2": MockYTChannel{
					IName:                     channel1,
					IID:                       "id123",
					IRSSURL:                   "123abc.com",
					IChannelURL:               "123abc.com",
					IArchivalMode:             ArchivalModeArchive,
					ILocalVideos:              &channel2Videos,
					ShouldErrorGetLocalVideos: false,
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

		ytcl := &MockYTChannelLoad{
			ShouldError: true,
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

		ytcl := &MockYTChannelLoad{
			ReturnValue: &map[string]YTChannel{
				"Channel1": MockYTChannel{
					IName:                     channel,
					IID:                       "id123",
					IRSSURL:                   "123abc.com",
					IChannelURL:               "123abc.com",
					IArchivalMode:             ArchivalModeArchive,
					ILocalVideos:              &expectedVideos,
					ShouldErrorGetLocalVideos: false,
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

		ytcl := &MockYTChannelLoad{
			ShouldError: true,
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

		ytcl := &MockYTChannelLoad{
			ReturnValue: nil,
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

		ytcl := &MockYTChannelLoad{
			ReturnValue: &map[string]YTChannel{
				"Channel1": MockYTChannel{
					IName:         "Channel1",
					IID:           "id123",
					IRSSURL:       "123abc.com",
					IChannelURL:   "123abc.com",
					IArchivalMode: ArchivalModeArchive,
					ILocalVideos: &[]Video{
						Video{
							"Video1",
							"dCqJ6iPHus0",
							"mp4",
							"/a/path",
						},
					},
					ShouldErrorGetLocalVideos: false,
				},
				"Channel2": MockYTChannel{
					IName:         "Channel2",
					IID:           "id123",
					IRSSURL:       "123abc.com",
					IChannelURL:   "123abc.com",
					IArchivalMode: ArchivalModeArchive,
					ILocalVideos: &[]Video{
						Video{
							"Video1",
							"dCqJ6iPHus0",
							"mp4",
							"/a/path",
						},
					},
					ShouldErrorGetLocalVideos: true,
				},
			},
		}

		_, err := getAllLocalVideos(cf, ytcl)

		if err == nil {
			t.Error("getAllLocalVideos should have thrown an error")
		}
	})
}
