package api

import (
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"reflect"
	"testing"
)

var cf = config.Config{
	VideoDirPath: "/a/path",
}

func TestGetChannels(t *testing.T) {
	t.Run("getChannels should return correct results", func(t *testing.T) {
		expectedYTChannels := []collection.YTChannelData{
			collection.YTChannelData{
				IName:         "Channel1",
				IID:           "asdfasdf",
				IRSSURL:       "http://testurl",
				IChannelURL:   "http://testurl",
				IArchivalMode: collection.ArchivalModeArchive,
			},
			collection.YTChannelData{
				IName:         "Channel2",
				IID:           "asdfasdf",
				IRSSURL:       "http://testurl2",
				IChannelURL:   "http://testurl2",
				IArchivalMode: collection.ArchivalModeArchive,
			},
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{
				"Channel1": collection.MockYTChannel{
					IName:                     expectedYTChannels[0].IName,
					IID:                       expectedYTChannels[0].IID,
					IRSSURL:                   expectedYTChannels[0].IRSSURL,
					IChannelURL:               expectedYTChannels[0].IChannelURL,
					IArchivalMode:             expectedYTChannels[0].IArchivalMode,
					ILocalVideos:              nil,
					ShouldErrorGetLocalVideos: false,
				},
				"Channel2": collection.MockYTChannel{
					IName:                     expectedYTChannels[1].IName,
					IID:                       expectedYTChannels[1].IID,
					IRSSURL:                   expectedYTChannels[1].IRSSURL,
					IChannelURL:               expectedYTChannels[1].IChannelURL,
					IArchivalMode:             expectedYTChannels[1].IArchivalMode,
					ILocalVideos:              nil,
					ShouldErrorGetLocalVideos: false,
				},
			},
		}

		ytChannels, err := getChannels(&cf, &ytcl)

		if err != nil {
			t.Error(err)
		}

		match1 := false
		match2 := false
		for _, ytch := range *ytChannels {
			if reflect.DeepEqual(expectedYTChannels[0], ytch) {
				match1 = true
			}

			if reflect.DeepEqual(expectedYTChannels[1], ytch) {
				match2 = true
			}
		}

		if !match1 || !match2 {
			t.Errorf("getChannels did not return correct result. Expected\n%+v\ngot\n%+v", expectedYTChannels, *ytChannels)
		}
	})

	t.Run("getChannels should return an empty YTChannelData slice when no channels are available", func(t *testing.T) {
		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{},
		}

		ytChannels, err := getChannels(&cf, &ytcl)

		if err != nil {
			t.Error(err)
		}

		if len(*ytChannels) != 0 {
			t.Errorf("getChannels did not return correct result. Expected 0 results, ngot\n%+v", *ytChannels)
		}
	})

	t.Run("getChannels should return an error when the channel loader returns an error", func(t *testing.T) {
		ytcl := collection.MockYTChannelLoad{
			ReturnValue: nil,
			ShouldError: true,
		}

		_, err := getChannels(&cf, &ytcl)

		if err == nil {
			t.Error("Should have returned an error")
		}
	})
}

func TestGetChannelByID(t *testing.T) {
	t.Run("getChannelByID should return correct result", func(t *testing.T) {
		expectedYTChannel := collection.MockYTChannel{
			IName:                     "Channel1",
			IID:                       "asdfasdf",
			IRSSURL:                   "http://testurl",
			IChannelURL:               "http://testurl",
			IArchivalMode:             collection.ArchivalModeArchive,
			ILocalVideos:              nil,
			ShouldErrorGetLocalVideos: false,
		}

		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{
				"Channel1": collection.MockYTChannel{
					IName:                     "Channel2",
					IID:                       "asdfasdf",
					IRSSURL:                   "http://testurl2",
					IChannelURL:               "http://testurl2",
					IArchivalMode:             collection.ArchivalModeArchive,
					ILocalVideos:              nil,
					ShouldErrorGetLocalVideos: false,
				},
				"Channel2": expectedYTChannel,
			},
		}

		ytChannels, err := getChannelByID("Channel2", &cf, &ytcl)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expectedYTChannel, *ytChannels) {
			t.Errorf("getChannels did not return correct result. Expected\n%+v\ngot\n%+v", expectedYTChannel, *ytChannels)
		}
	})

	t.Run("getChannelByID should return an error when the channel loader returns an error", func(t *testing.T) {
		ytcl := collection.MockYTChannelLoad{
			ReturnValue: nil,
			ShouldError: true,
		}

		_, err := getChannelByID("someID", &cf, &ytcl)

		if err == nil {
			t.Error("Should have returned an error")
		}
	})

	t.Run("getChannelByID should return a nil pointer when no YTChannel is found", func(t *testing.T) {
		ytcl := collection.MockYTChannelLoad{
			ReturnValue: &map[string]collection.YTChannel{},
		}

		ytChannel, err := getChannelByID("something bad", &cf, &ytcl)

		if err != nil {
			t.Error(err)
		}

		if ytChannel != nil {
			t.Errorf("getChannelByID was expected to return a nil pointer. Got %+v", *ytChannel)
		}
	})
}

func TestCheckChannelUpdates(t *testing.T) {
	t.Run("checkChannelUpdates returns correct respsonse", func(t *testing.T) {
		localVideoMockData := *collection.GetVideoMockData()
		response, err := checkChannelUpdates(
			"Channel1",
			&cf,
			&collection.MockYTChannelLoad{
				ReturnValue: &map[string]collection.YTChannel{
					"Channel1": collection.MockYTChannel{
						IName:         "Test Guy",
						IID:           "UCS-WzPVpAAli-1IfEG2lN8A",
						IRSSURL:       "http://testurl1",
						IChannelURL:   "http://testurl1",
						IArchivalMode: collection.ArchivalModeArchive,
						IChannelType:  collection.ChannelTypeChannel,
						ILocalVideos: &[]collection.LocalVideo{
							localVideoMockData[1],
							localVideoMockData[2],
						},
						ShouldErrorGetLocalVideos: false,
					},
				},
			},
			&youtubeapi.MockAPI{},
		)
		if err != nil {
			t.Errorf("checkChannelUpdates returned an error %s", err)
		}

		mockData := *GetVideoMockData()
		expectedResponse := []Video{
			mockData[0],
		}

		if !reflect.DeepEqual(expectedResponse, *response) {
			t.Errorf(testutils.MismatchError("checkChannelUpdates", expectedResponse, *response))
		}
	})

	t.Run("checkChannelUpdates returns nothing when all videos are local", func(t *testing.T) {
		response, err := checkChannelUpdates(
			"Channel1",
			&cf,
			&collection.MockYTChannelLoad{
				ReturnValue: &map[string]collection.YTChannel{
					"Channel1": collection.MockYTChannel{
						IName:                     "Test Guy",
						IID:                       "UCS-WzPVpAAli-1IfEG2lN8A",
						IRSSURL:                   "http://testurl1",
						IChannelURL:               "http://testurl1",
						IArchivalMode:             collection.ArchivalModeArchive,
						ILocalVideos:              collection.GetVideoMockData(),
						ShouldErrorGetLocalVideos: false,
					},
				},
			},
			&youtubeapi.MockAPI{},
		)
		if err != nil {
			t.Errorf("checkChannelUpdates returned an error %s", err)
		}

		if len(*response) > 0 {
			t.Errorf("checkChannelUpdates did not return correct response. Expected an empty slice, got %+v", response)
		}
	})

	t.Run("checkChannelUpdates returns nothing remote has no videos", func(t *testing.T) {
		channelMockResponse, err := youtubeapi.GetVideoEmptyMockData()
		if err != nil {
			t.Error(err)
		}

		response, err := checkChannelUpdates(
			"Channel1",
			&cf,
			&collection.MockYTChannelLoad{
				ReturnValue: &map[string]collection.YTChannel{
					"Channel1": collection.MockYTChannel{
						IName:                     "Test Guy",
						IID:                       "UCS-WzPVpAAli-1IfEG2lN8A",
						IRSSURL:                   "http://testurl1",
						IChannelURL:               "http://testurl1",
						IArchivalMode:             collection.ArchivalModeArchive,
						ILocalVideos:              collection.GetVideoMockData(),
						ShouldErrorGetLocalVideos: false,
					},
				},
			},
			&youtubeapi.MockAPI{
				GetVideosForChannelReponse: channelMockResponse,
			},
		)
		if err != nil {
			t.Errorf("checkChannelUpdates returned an error %s", err)
		}

		if len(*response) > 0 {
			t.Errorf("checkChannelUpdates did not return correct response. Expected an empty slice, got %+v", response)
		}
	})

	t.Run("checkChannelUpdates returns everything when no videos are local", func(t *testing.T) {
		response, err := checkChannelUpdates(
			"Channel1",
			&cf,
			&collection.MockYTChannelLoad{
				ReturnValue: &map[string]collection.YTChannel{
					"Channel1": collection.MockYTChannel{
						IName:                     "Test Guy",
						IID:                       "UCS-WzPVpAAli-1IfEG2lN8A",
						IRSSURL:                   "http://testurl1",
						IChannelURL:               "http://testurl1",
						IArchivalMode:             collection.ArchivalModeArchive,
						ILocalVideos:              &[]collection.LocalVideo{},
						ShouldErrorGetLocalVideos: false,
					},
				},
			},
			&youtubeapi.MockAPI{},
		)
		if err != nil {
			t.Errorf("checkChannelUpdates returned an error %s", err)
		}

		expectedResponse := GetVideoMockData()

		if !reflect.DeepEqual(*expectedResponse, *response) {
			t.Errorf(testutils.MismatchError("checkChannelUpdates", *expectedResponse, *response))
		}
	})

	t.Run("returns error when channel cannot be found", func(t *testing.T) {
		_, err := checkChannelUpdates(
			"Channel1",
			&cf,
			&collection.MockYTChannelLoad{
				ReturnValue: &map[string]collection.YTChannel{
					"Channel2": collection.MockYTChannel{
						IName:                     "Test Guy",
						IID:                       "UCS-WzPVpAAli-1IfEG2lN8A",
						IRSSURL:                   "http://testurl1",
						IChannelURL:               "http://testurl1",
						IArchivalMode:             collection.ArchivalModeArchive,
						ILocalVideos:              &[]collection.LocalVideo{},
						ShouldErrorGetLocalVideos: false,
					},
				},
			},
			&youtubeapi.MockAPI{},
		)

		if err == nil {
			t.Errorf("expected error")
		}
	})

	t.Run("returns error if Youtube API errors", func(t *testing.T) {
		_, err := checkChannelUpdates(
			"Channel1",
			&cf,
			&collection.MockYTChannelLoad{
				ReturnValue: &map[string]collection.YTChannel{
					"Channel1": collection.MockYTChannel{
						IName:                     "Test Guy",
						IID:                       "UCS-WzPVpAAli-1IfEG2lN8A",
						IRSSURL:                   "http://testurl1",
						IChannelURL:               "http://testurl1",
						IArchivalMode:             collection.ArchivalModeArchive,
						ILocalVideos:              &[]collection.LocalVideo{},
						ShouldErrorGetLocalVideos: false,
					},
				},
			},
			&youtubeapi.MockAPI{GetVideosForChannelReturnError: true},
		)

		if err == nil {
			t.Errorf("expected error")
		}
	})

	t.Run("returns error if collection local video loader errors", func(t *testing.T) {
		_, err := checkChannelUpdates(
			"Channel1",
			&cf,
			&collection.MockYTChannelLoad{
				ReturnValue: &map[string]collection.YTChannel{
					"Channel1": collection.MockYTChannel{
						IName:                     "Test Guy",
						IID:                       "UCS-WzPVpAAli-1IfEG2lN8A",
						IRSSURL:                   "http://testurl1",
						IChannelURL:               "http://testurl1",
						IArchivalMode:             collection.ArchivalModeArchive,
						ILocalVideos:              &[]collection.LocalVideo{},
						ShouldErrorGetLocalVideos: true,
					},
				},
			},
			&youtubeapi.MockAPI{},
		)

		if err == nil {
			t.Errorf("expected error")
		}
	})
}
