package youtubeapi

import (
	"io/ioutil"
	"reflect"
	"testing"
)

var vlExpected = VideoListResponse{
	Kind: "youtube#videoListResponse",
	Etag: "1-lTCZCHtgPgr709KQ0ef2Mu4oM",
	Items: []Video{
		Video{
			Kind: "youtube#video",
			Etag: "jB-DuI2TOpg-o1d5hnzty8kExw8",
			ID:   "18-elPdai_1",
			Snippet: VideoSnippet{
				PublishedAt: "2012-10-01T15:27:35Z",
				ChannelID:   "UCS-WzPVpAAli-1IfEG2lN8A",
				Title:       "Test Video New",
				Description: "Test Description New",
				Thumbnails: ThumbnailDetails{
					Default: Thumbnail{
						URL:    "https://i2.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
						Width:  120,
						Height: 90,
					},
					Medium: Thumbnail{
						URL:    "https://i.ytimg.com/vi/KQA9Na4aOa1/mqdefault.jpg",
						Width:  320,
						Height: 180,
					},
					High: Thumbnail{
						URL:    "https://i.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
						Width:  480,
						Height: 360,
					},
					Standard: Thumbnail{
						URL:    "https://i.ytimg.com/vi/KQA9Na4aOa1/sddefault.jpg",
						Width:  640,
						Height: 480,
					},
					Maxres: Thumbnail{
						URL:    "https://i.ytimg.com/vi/KQA9Na4aOa1/maxresdefault.jpg",
						Width:  1280,
						Height: 720,
					},
				},
				ChannelTitle: "Test Guy",
				Tags: []string{
					"IMA",
					"TEST",
					"TAG",
				},
				CategoryID:           "22",
				LiveBroadcastContent: "none",
				DefaultLanguage:      "en",
				Localized: LocalizationDetail{
					Title:       "Test Video New",
					Description: "Test Description New",
				},
				DefaultAudioLanguage: "en",
			},
		},
		Video{
			Kind: "youtube#video",
			Etag: "GT_wjBomCnhJpvvlExadAU7E_t8",
			ID:   "OGK8gnP4TfA",
			Snippet: VideoSnippet{
				PublishedAt: "2012-03-02T19:03:16Z",
				ChannelID:   "UCS-WzPVpAAli-1IfEG2lN8A",
				Title:       "Test Video 1",
				Description: "Test Description",
				Thumbnails: ThumbnailDetails{
					Default: Thumbnail{
						URL:    "https://i.ytimg.com/vi/OGK8gnP4TfA/default.jpg",
						Width:  120,
						Height: 90,
					},
					Medium: Thumbnail{
						URL:    "https://i.ytimg.com/vi/OGK8gnP4TfA/mqdefault.jpg",
						Width:  320,
						Height: 180,
					},
					High: Thumbnail{
						URL:    "https://i.ytimg.com/vi/OGK8gnP4TfA/hqdefault.jpg",
						Width:  480,
						Height: 360,
					},
					Standard: Thumbnail{
						URL:    "https://i.ytimg.com/vi/OGK8gnP4TfA/sddefault.jpg",
						Width:  640,
						Height: 480,
					},
					Maxres: Thumbnail{
						URL:    "https://i.ytimg.com/vi/OGK8gnP4TfA/maxresdefault.jpg",
						Width:  1280,
						Height: 720,
					},
				},
				ChannelTitle: "Test Guy",
				Tags: []string{
					"The",
					"Best",
					"Test",
					"Video",
				},
				CategoryID:           "22",
				LiveBroadcastContent: "none",
				DefaultLanguage:      "en",
				Localized: LocalizationDetail{
					Title:       "test Video 1",
					Description: "Test Description",
				},
			},
		},
		Video{
			Kind: "youtube#video",
			Etag: "q3lLb6e10Mo-amZyWSBo5HNJpAU",
			ID:   "FazJqPQ6xSs",
			Snippet: VideoSnippet{
				PublishedAt: "2014-06-27T15:10:18Z",
				ChannelID:   "UCS-WzPVpAAli-1IfEG2lN8A",
				Title:       "Test Video 2",
				Description: "Test Description 2",
				Thumbnails: ThumbnailDetails{
					Default: Thumbnail{
						URL:    "https://i.ytimg.com/vi/FazJqPQ6xSs/default.jpg",
						Width:  120,
						Height: 90,
					},
					Medium: Thumbnail{
						URL:    "https://i.ytimg.com/vi/FazJqPQ6xSs/mqdefault.jpg",
						Width:  320,
						Height: 180,
					},
					High: Thumbnail{
						URL:    "https://i.ytimg.com/vi/FazJqPQ6xSs/hqdefault.jpg",
						Width:  480,
						Height: 360,
					},
					Standard: Thumbnail{
						URL:    "https://i.ytimg.com/vi/FazJqPQ6xSs/sddefault.jpg",
						Width:  640,
						Height: 480,
					},
					Maxres: Thumbnail{
						URL:    "https://i.ytimg.com/vi/FazJqPQ6xSs/maxresdefault.jpg",
						Width:  1280,
						Height: 720,
					},
				},
				ChannelTitle: "Test Guy",
				Tags: []string{
					"Yet",
					"ANOTHER",
					"test",
					"video",
				},
				CategoryID:           "26",
				LiveBroadcastContent: "none",
				DefaultLanguage:      "en",
				Localized: LocalizationDetail{
					Title:       "Test Video 2",
					Description: "Test Description 2",
				},
			},
		},
	},
	PageInfo: PageInfo{
		TotalResults:   3,
		ResultsPerPage: 3,
	},
}

func TestConvertVideoAPIResponse(t *testing.T) {
	t.Run("Correctly takes a JSON API response and parses it into expected VideoListResponse object", func(t *testing.T) {
		file, err := ioutil.ReadFile("./testfiles/videorequest.json")
		if err != nil {
			t.Errorf("Loading VideoRequest json failed: %s", err)
		}

		vl, err := convertVideoAPIResponse(string(file))
		if err != nil {
			t.Errorf("convertVideoAPIResponse returned an error %s", err)
		}

		if !reflect.DeepEqual(*vl, vlExpected) {
			t.Errorf("VideoListResponse results are different.\nExpected %+v\ngot      %+v", vlExpected, *vl)
		}
	})

	t.Run("Throws an error when invalid JSON API response is entered", func(t *testing.T) {
		_, err := convertVideoAPIResponse("")
		if err == nil {
			t.Error("convertVideoAPIResponse did not return an error with invalid JSON")
		}

		_, err = convertVideoAPIResponse("{{")
		if err == nil {
			t.Error("convertVideoAPIResponse did not return an error with invalid JSON")
		}
	})
}
