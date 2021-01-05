package youtubeapi

import (
	"fmt"
	"io/ioutil"
)

var vlVideo1 = Video{
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

var vlVideo2 = Video{
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
}

var vlVideo3 = Video{
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
}

// VLExpectedFull contains an expected VideoListResponse for the videoresponse.json file
// in the testfiles directory
var vlExpectedFull = VideoMetadataResponse{
	Kind:  "youtube#videoListResponse",
	Etag:  "1-lTCZCHtgPgr709KQ0ef2Mu4oM",
	Items: []Video{vlVideo1, vlVideo2, vlVideo3},
	PageInfo: PageInfo{
		TotalResults:   3,
		ResultsPerPage: 3,
	},
}

// VLExpectedSingleVideo contains an expected VideoListResponse for the videoresponse_single.json file
// in the testfiles directory
var vlExpectedSingleVideo = VideoMetadataResponse{
	Kind:  "youtube#videoListResponse",
	Etag:  "1-lTCZCHtgPgr709KQ0ef2Mu4oM",
	Items: []Video{vlVideo1},
	PageInfo: PageInfo{
		TotalResults:   1,
		ResultsPerPage: 1,
	},
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
