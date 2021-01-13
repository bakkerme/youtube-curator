package api

// GetVideoMockData returns mock API response data
func GetVideoMockData() *[]Video {
	return &[]Video{
		Video{
			ID:          "18-elPdai_1",
			Creator:     "Test Guy",
			Description: "Test Description New",
			PublishedAt: "2012-10-01T15:27:35Z",
			Title:       "Test Video New",
		},
		Video{
			ID:          "OGK8gnP4TfA",
			Creator:     "Test Guy",
			Description: "Test Description",
			PublishedAt: "2018-12-03T23:20:21Z",
			Title:       "Test Video 1",
		},
		Video{
			ID:          "FazJqPQ6xSs",
			Creator:     "Test Guy",
			Description: "Test Description 2",
			PublishedAt: "2019-06-03T19:00:06Z",
			Title:       "Test Video 2",
		},
	}
}
