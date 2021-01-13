package youtubeapi

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
)

// MockAPI mocks access to the Youtube API
type MockAPI struct {
	GetVideoMetadataResponse    *VideoMetadataResponse
	GetVideoMetadataReturnError bool

	GetVideosForChannelReponse     *VideoMetadataResponse
	GetVideosForChannelReturnError bool
}

// GetVideoMetadata gets information on the video whos IDs are provided from the Youtube API
func (ytAPI *MockAPI) GetVideoMetadata(ids *[]string, cf *config.Config) (*VideoMetadataResponse, error) {
	if ytAPI.GetVideosForChannelReturnError {
		return nil, errors.New("Something bad happened")
	}

	if ytAPI.GetVideoMetadataResponse != nil {
		return ytAPI.GetVideoMetadataResponse, nil
	}

	vl, err := convertAPIResponse(string(SearchResponseJSON), apiSearch)
	if err != nil {
		return nil, fmt.Errorf("convertVideoAPIResponse returned an error %s", err)
	}

	for i, item := range vl.Items {
		if len(*ids) < i {
			item.ID = (*ids)[i]
		}
	}

	return vl, nil
}

// GetVideosForChannel returns videos for a provided YT Channel from the YouTube API
func (ytAPI *MockAPI) GetVideosForChannel(ytc collection.YTChannel, cf *config.Config) (*VideoMetadataResponse, error) {
	if ytAPI.GetVideosForChannelReturnError {
		return nil, errors.New("Something bad happened")
	}

	if ytAPI.GetVideosForChannelReponse != nil {
		return ytAPI.GetVideosForChannelReponse, nil
	}

	vl, err := convertAPIResponse(string(SearchResponseJSON), apiSearch)
	if err != nil {
		return nil, fmt.Errorf("convertVideoAPIResponse returned an error %s", err)
	}

	return vl, nil
}

// GetVideoMockData returns the mock data that is used for API responses
func GetVideoMockData() (*VideoMetadataResponse, error) {
	vl, err := convertAPIResponse(string(SearchResponseJSON), apiSearch)
	if err != nil {
		return nil, fmt.Errorf("convertVideoAPIResponse returned an error %s", err)
	}

	return vl, nil
}

// GetVideoEmptyMockData returns empty mock data that is used for API responses
func GetVideoEmptyMockData() (*VideoMetadataResponse, error) {
	vl, err := convertAPIResponse(string(SearchResponseEmptyJSON), apiSearch)
	if err != nil {
		return nil, fmt.Errorf("convertVideoAPIResponse returned an error %s", err)
	}

	return vl, nil

}
