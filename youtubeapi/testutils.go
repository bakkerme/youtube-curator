package youtubeapi

import (
	"fmt"
	"io/ioutil"
)

var vlVideoSingle = Video{
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
}

var searchResult = []Video{
	Video{
		Kind: "youtube#searchResult",
		Etag: "mXlJ6j7wGlbm9SjkEyvihSwo0tU",
		ID:   "18-elPdai_1",
		Snippet: VideoSnippet{
			PublishedAt: "2012-10-01T15:27:35Z",
			ChannelID:   "UCS-WzPVpAAli-1IfEG2lN8A",
			Title:       "Test Video New",
			Description: "Test Description New",
			Thumbnails: ThumbnailDetails{
				Default: Thumbnail{
					URL:    "https://i.ytimg.com/vi/18-elPdai_1/default.jpg",
					Width:  120,
					Height: 90,
				},
				Medium: Thumbnail{
					URL:    "https://i.ytimg.com/vi/18-elPdai_1/mqdefault.jpg",
					Width:  320,
					Height: 180,
				},
				High: Thumbnail{
					URL:    "https://i.ytimg.com/vi/18-elPdai_1/hqdefault.jpg",
					Width:  480,
					Height: 360,
				},
			},
			ChannelTitle:         "Test Guy",
			LiveBroadcastContent: "none",
			PublishTime:          "2012-10-01T15:27:35Z",
		},
	},
	Video{
		Kind: "youtube#searchResult",
		Etag: "-fusrPK0jUxsR3-7UT7as7j4sGM",
		ID:   "OGK8gnP4TfA",
		Snippet: VideoSnippet{
			PublishedAt: "2018-12-03T23:20:21Z",
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
			},
			ChannelTitle:         "Test Guy",
			LiveBroadcastContent: "none",
			PublishTime:          "2018-12-03T23:20:21Z",
		},
	},
	Video{
		Kind: "youtube#searchResult",
		Etag: "nWfU-BRD9p-BGwQ_oFpSv7YmaeI",
		ID:   "FazJqPQ6xSs",
		Snippet: VideoSnippet{
			PublishedAt: "2019-06-03T19:00:06Z",
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
			},
			ChannelTitle:         "Test Guy",
			LiveBroadcastContent: "none",
			PublishTime:          "2019-06-03T19:00:06Z",
		},
	},
}

var playlistItemsResponse = []Video{
	Video{
		Kind: "youtube#playlistItem",
		Etag: "b_mpir6DuPAQXBP2QtAT9y2Dp8c",
		ID:   "GhMOw141DKg",
		Snippet: VideoSnippet{
			PublishedAt: "2020-09-30T19:22:21Z",
			ChannelID:   "UCjU-Cwjfqbo2hMRItlXwnnQ",
			Title:       "Test Video New",
			Description: "Test Description New",
			Thumbnails: ThumbnailDetails{
				Default: Thumbnail{
					URL:    "https://i.ytimg.com/vi/GhMOw141DKg/default.jpg",
					Width:  120,
					Height: 90,
				},
				Medium: Thumbnail{
					URL:    "https://i.ytimg.com/vi/GhMOw141DKg/mqdefault.jpg",
					Width:  320,
					Height: 180,
				},
				High: Thumbnail{
					URL:    "https://i.ytimg.com/vi/GhMOw141DKg/hqdefault.jpg",
					Width:  480,
					Height: 360,
				},
				Standard: Thumbnail{
					URL:    "https://i.ytimg.com/vi/GhMOw141DKg/sddefault.jpg",
					Width:  640,
					Height: 480,
				},
			},
			ResourceID: ResourceID{
				Kind: "youtube#video",
				ID:   "GhMOw141DKg",
			},
			ChannelTitle: "Test Guy",
			PlaylistID:   "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
			Position:     0,
		},
	},
	Video{
		Kind: "youtube#playlistItem",
		Etag: "Nzlqatt6khGGAufkcqgLhEZ6fbg",
		ID:   "YEa2aj9KYQA",
		Snippet: VideoSnippet{
			PublishedAt: "2020-03-30T00:22:24Z",
			ChannelID:   "UCjU-Cwjfqbo2hMRItlXwnnQ",
			Title:       "Test Video 1",
			Description: "Test Description 1",
			Thumbnails: ThumbnailDetails{
				Default: Thumbnail{
					URL:    "https://i.ytimg.com/vi/YEa2aj9KYQA/default.jpg",
					Width:  120,
					Height: 90,
				},
				Medium: Thumbnail{
					URL:    "https://i.ytimg.com/vi/YEa2aj9KYQA/mqdefault.jpg",
					Width:  320,
					Height: 180,
				},
				High: Thumbnail{
					URL:    "https://i.ytimg.com/vi/YEa2aj9KYQA/hqdefault.jpg",
					Width:  480,
					Height: 360,
				},
				Standard: Thumbnail{
					URL:    "https://i.ytimg.com/vi/YEa2aj9KYQA/sddefault.jpg",
					Width:  640,
					Height: 480,
				},
			},
			ResourceID: ResourceID{
				Kind: "youtube#video",
				ID:   "YEa2aj9KYQA",
			},
			ChannelTitle: "Test Guy",
			PlaylistID:   "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
			Position:     1,
		},
	},
	Video{
		Kind: "youtube#playlistItem",
		Etag: "5LQBWt1kwMx7Uqtb2gTHP0bnHOw",
		ID:   "fZYBhmteJDE",
		Snippet: VideoSnippet{
			PublishedAt: "2020-01-23T16:36:27Z",
			ChannelID:   "UCjU-Cwjfqbo2hMRItlXwnnQ",
			Title:       "Test Video 2",
			Description: "Test Description 2",
			Thumbnails: ThumbnailDetails{
				Default: Thumbnail{
					URL:    "https://i.ytimg.com/vi/fZYBhmteJDE/default.jpg",
					Width:  120,
					Height: 90,
				},
				Medium: Thumbnail{
					URL:    "https://i.ytimg.com/vi/fZYBhmteJDE/mqdefault.jpg",
					Width:  320,
					Height: 180,
				},
				High: Thumbnail{
					URL:    "https://i.ytimg.com/vi/fZYBhmteJDE/hqdefault.jpg",
					Width:  480,
					Height: 360,
				},
				Standard: Thumbnail{
					URL:    "https://i.ytimg.com/vi/fZYBhmteJDE/sddefault.jpg",
					Width:  640,
					Height: 480,
				},
			},
			ResourceID: ResourceID{
				Kind: "youtube#video",
				ID:   "fZYBhmteJDE",
			},
			ChannelTitle: "Test Guy",
			PlaylistID:   "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
			Position:     2,
		},
	},
}

// expectedSearchResponse contains an expected VideoMetadataResponse for the videoresponse.json file
// in the testfiles directory
var expectedSearchResponse = VideoMetadataResponse{
	Kind:  "youtube#searchListResponse",
	Etag:  "JUwDkhuk2MfT9JohAdWLTXV45aM",
	Items: searchResult,
	PageInfo: PageInfo{
		TotalResults:   522,
		ResultsPerPage: 3,
	},
	NextPageToken: "CAUQAA",
	RegionCode:    "AU",
}

// expectedVideoResponse contains an expected VideoMetadataResponse for the videoresponse_single.json file
// in the testfiles directory
var expectedVideoResponse = VideoMetadataResponse{
	Kind:  "youtube#videoListResponse",
	Etag:  "1-lTCZCHtgPgr709KQ0ef2Mu4oM",
	Items: []Video{vlVideoSingle},
	PageInfo: PageInfo{
		TotalResults:   1,
		ResultsPerPage: 1,
	},
}

// expectedPlaylistItemsResponse contains an expected VideoMetadataResponse for the videoresponse.json file
// in the testfiles directory
var expectedPlaylistItemsResponse = VideoMetadataResponse{
	Kind:  "youtube#playlistItemListResponse",
	Etag:  "c6a05HjuPsmPxbDxMmt-196SvPI",
	Items: playlistItemsResponse,
	PageInfo: PageInfo{
		TotalResults:   54,
		ResultsPerPage: 3,
	},
	NextPageToken: "CAUQAA",
}

// LoadVideoListSingleTestFile loads the test file that is the counterpart to the VLExpected
// struct, that way it can be loaded up for unmarshalling comparison testing
func loadFile(filePath string) ([]byte, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return file, fmt.Errorf("Loading VideoRequest json failed: %s", err)
	}

	return file, nil
}
