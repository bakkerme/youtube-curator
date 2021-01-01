package api

import (
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"reflect"
	"testing"
)

func TestGetChannels(t *testing.T) {
	t.Run("getChannels should return correct results", func(t *testing.T) {
		cfg := config.Config{
			VideoDirPath: "/a/path",
		}

		expectedYTChannels := []collection.YTChannelData{
			collection.YTChannelData{
				IName:         "Channel1",
				IRSSURL:       "http://testurl",
				IChannelURL:   "http://testurl",
				IArchivalMode: "archive",
			},
			collection.YTChannelData{
				IName:         "Channel2",
				IRSSURL:       "http://testurl2",
				IChannelURL:   "http://testurl2",
				IArchivalMode: "curated",
			},
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{
				"Channel1": collection.MockYTChannel{
					IName:                     expectedYTChannels[0].IName,
					IRSSURL:                   expectedYTChannels[0].IRSSURL,
					IChannelURL:               expectedYTChannels[0].IChannelURL,
					IArchivalMode:             expectedYTChannels[0].IArchivalMode,
					ILocalVideos:              nil,
					ShouldErrorGetLocalVideos: false,
				},
				"Channel2": collection.MockYTChannel{
					IName:                     expectedYTChannels[1].IName,
					IRSSURL:                   expectedYTChannels[1].IRSSURL,
					IChannelURL:               expectedYTChannels[1].IChannelURL,
					IArchivalMode:             expectedYTChannels[1].IArchivalMode,
					ILocalVideos:              nil,
					ShouldErrorGetLocalVideos: false,
				},
			},
		}

		ytChannels, err := getChannels(&cfg, &ytcl)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expectedYTChannels, *ytChannels) {
			t.Errorf("getChannels did not return correct result. Expected\n%+v\ngot\n%+v", expectedYTChannels, *ytChannels)
		}
	})

	t.Run("getChannels should return an empty YTChannelData slice when no channels are available", func(t *testing.T) {
		cfg := config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{},
		}

		ytChannels, err := getChannels(&cfg, &ytcl)

		if err != nil {
			t.Error(err)
		}

		if len(*ytChannels) != 0 {
			t.Errorf("getChannels did not return correct result. Expected 0 results, ngot\n%+v", *ytChannels)
		}
	})

	t.Run("getChannels should return an error when the channel loader returns an error", func(t *testing.T) {
		cfg := config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: nil,
			ShouldError: true,
		}

		_, err := getChannels(&cfg, &ytcl)

		if err == nil {
			t.Error("Should have returned an error")
		}
	})
}

func TestGetChannelByID(t *testing.T) {
	t.Run("getChannelByID should return correct result", func(t *testing.T) {
		cfg := config.Config{
			VideoDirPath: "/a/path",
		}

		expectedYTChannel := collection.MockYTChannel{
			IName:                     "Channel1",
			IRSSURL:                   "http://testurl",
			IChannelURL:               "http://testurl",
			IArchivalMode:             "archive",
			ILocalVideos:              nil,
			ShouldErrorGetLocalVideos: false,
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{
				"Channel1": collection.MockYTChannel{
					IName:                     "Channel2",
					IRSSURL:                   "http://testurl2",
					IChannelURL:               "http://testurl2",
					IArchivalMode:             "curated",
					ILocalVideos:              nil,
					ShouldErrorGetLocalVideos: false,
				},
				"Channel2": expectedYTChannel,
			},
		}

		ytChannels, err := getChannelByID("Channel2", &cfg, &ytcl)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expectedYTChannel, *ytChannels) {
			t.Errorf("getChannels did not return correct result. Expected\n%+v\ngot\n%+v", expectedYTChannel, *ytChannels)
		}
	})

	t.Run("getChannelByID should return an error when the channel loader returns an error", func(t *testing.T) {
		cfg := config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: nil,
			ShouldError: true,
		}

		_, err := getChannelByID("someID", &cfg, &ytcl)

		if err == nil {
			t.Error("Should have returned an error")
		}
	})

	t.Run("getChannelByID should return a nil pointer when no YTChannel is found", func(t *testing.T) {
		cfg := config.Config{
			VideoDirPath: "/a/path",
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{},
		}

		ytChannel, err := getChannelByID("something bad", &cfg, &ytcl)

		if err != nil {
			t.Error(err)
		}

		if ytChannel != nil {
			t.Errorf("getChannelByID was expected to return a nil pointer. Got %+v", *ytChannel)
		}
	})
}
